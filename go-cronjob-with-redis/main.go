package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mandarinkb/go-cronjob-with-redis/util"
	"github.com/redis/go-redis/v9"
)

func generateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// add job: store hash + zadd (score=nextRun)
func addCronJob(ctx context.Context, inputID, name, schedule, command string) (string, error) {
	id := inputID
	if id == "" {
		id = generateID()
	}
	nextRun := util.GetNextRun(schedule)

	// pipeline: HSET + ZADD
	pipe := util.RDB.Pipeline()
	pipe.HSet(ctx, util.JobKey(id), map[string]interface{}{
		"id":       id,
		"name":     name,
		"command":  command,
		"schedule": schedule,
		"paused":   "false",
	})
	pipe.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: id})
	_, err := pipe.Exec(ctx)
	if err != nil {
		return "", err
	}

	util.Logger.Info("Added job",
		"id", id,
		"name", name,
		"schedule", schedule,
		"next_run", time.Unix(nextRun, 0).Format(time.RFC3339),
	)

	return id, nil
}

func deleteCronJob(ctx context.Context, id string) error {
	pipe := util.RDB.Pipeline()
	pipe.ZRem(ctx, util.ZsetKey(), id)
	pipe.Del(ctx, util.JobKey(id))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	util.Logger.Info("Deleted job", slog.String("id", id))
	return nil
}

func updateCronJob(ctx context.Context, id, newName, newSchedule, newCommand string) error {
	// We'll replace hash fields + update ZSet with new nextRun
	nextRun := util.GetNextRun(newSchedule)
	pipe := util.RDB.Pipeline()
	pipe.HSet(ctx, util.JobKey(id), map[string]interface{}{
		"name":     newName,
		"command":  newCommand,
		"schedule": newSchedule,
	})
	pipe.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: id})
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	util.Logger.Info("Updated job", slog.String("id", id), slog.String("schedule", newSchedule))
	return nil
}

func getAllCronJobs(ctx context.Context) ([]map[string]interface{}, error) {
	ids, err := util.RDB.ZRange(ctx, util.ZsetKey(), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	jobs := make([]map[string]interface{}, 0, len(ids))
	if len(ids) == 0 {
		return jobs, nil
	}

	pipe := util.RDB.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, 0, len(ids))
	for _, id := range ids {
		cmds = append(cmds, pipe.HGetAll(ctx, util.JobKey(id)))
	}

	if _, err := pipe.Exec(ctx); err != nil {
		util.Logger.Error("Pipeline exec failed", slog.Any("error", err))
	}

	for i, id := range ids {
		data, err := cmds[i].Result()
		if err != nil {
			util.Logger.Warn("Failed to fetch job data", slog.String("id", id), slog.Any("error", err))
			continue
		}
		if len(data) == 0 {
			util.Logger.Warn("No data for job", slog.String("id", id))
			continue
		}
		jobs = append(jobs, map[string]interface{}{
			"id":       id,
			"name":     data["name"],
			"schedule": data["schedule"],
			"command":  data["command"],
			"paused":   data["paused"],
		})
	}
	return jobs, nil
}

func getCronJob(ctx context.Context, id string) (map[string]interface{}, error) {
	data, err := util.RDB.HGetAll(ctx, util.JobKey(id)).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("job not found")
	}
	return map[string]interface{}{
		"id":       id,
		"name":     data["name"],
		"schedule": data["schedule"],
		"command":  data["command"],
		"paused":   data["paused"],
	}, nil
}

func jobExists(ctx context.Context, id string) (bool, error) {
	n, err := util.RDB.Exists(ctx, util.JobKey(id)).Result()
	return n > 0, err
}

// utility for structured slog fields
func slogFields(id, name, schedule string, nextRun int64) []slog.Attr {
	return []slog.Attr{
		slog.String("id", id),
		slog.String("name", name),
		slog.String("schedule", schedule),
		slog.String("next_run", time.Unix(nextRun, 0).Format(time.RFC3339)),
	}
}

func main() {
	util.Init()
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// create
	app.Post("/cron", func(c *fiber.Ctx) error {
		type CronJobInput struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Command  string `json:"command"`
		}
		var input CronJobInput
		if err := c.BodyParser(&input); err != nil {
			util.Logger.Warn("invalid body", slog.Any("error", err))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
		}
		if input.Name == "" || input.Schedule == "" || input.Command == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields"})
		}
		id, err := addCronJob(util.CTX, input.ID, input.Name, input.Schedule, input.Command)
		if err != nil {
			util.Logger.Error("add job failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add cron job"})
		}
		return c.JSON(fiber.Map{"message": "Cron job added", "id": id})
	})

	// delete
	app.Delete("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
		}
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := deleteCronJob(util.CTX, id); err != nil {
			util.Logger.Error("delete failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete"})
		}
		return c.JSON(fiber.Map{"message": "deleted", "id": id})
	})

	// update
	app.Put("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
		}
		type UpdateInput struct {
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Command  string `json:"command"`
		}
		var in UpdateInput
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
		}
		if in.Schedule == "" || in.Command == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields"})
		}
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := updateCronJob(util.CTX, id, in.Name, in.Schedule, in.Command); err != nil {
			util.Logger.Error("update failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update"})
		}
		return c.JSON(fiber.Map{"message": "updated", "id": id})
	})

	// list
	app.Get("/cron", func(c *fiber.Ctx) error {
		jobs, err := getAllCronJobs(util.CTX)
		if err != nil {
			util.Logger.Error("getAll failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch"})
		}
		return c.JSON(jobs)
	})

	// get single
	app.Get("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
		}
		job, err := getCronJob(util.CTX, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
		}
		return c.JSON(job)
	})

	// run now
	app.Patch("/cron/:id/now", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
		}
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		now := float64(time.Now().Unix())
		if err := util.RDB.ZAdd(util.CTX, util.ZsetKey(), redis.Z{Score: now, Member: id}).Err(); err != nil {
			util.Logger.Error("set now failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed"})
		}
		return c.JSON(fiber.Map{"message": "scheduled now", "id": id, "time": now})
	})

	// pause
	app.Post("/cron/:id/pause", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
		}
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := util.RDB.HSet(util.CTX, util.JobKey(id), "paused", "true").Err(); err != nil {
			util.Logger.Error("pause failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed"})
		}
		return c.JSON(fiber.Map{"message": "paused", "id": id})
	})

	// resume
	app.Post("/cron/:id/resume", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
		}
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := util.RDB.HSet(util.CTX, util.JobKey(id), "paused", "false").Err(); err != nil {
			util.Logger.Error("resume failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed"})
		}
		return c.JSON(fiber.Map{"message": "resumed", "id": id})
	})

	// graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		util.Logger.Info("shutting down api...")
		app.Shutdown()
	}()

	port := util.GetEnv("API_PORT", "3000")
	util.Logger.Info("starting api", slog.String("port", port))
	app.Listen(":" + port)
}
