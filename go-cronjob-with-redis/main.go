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

// ── Standalone Job helpers ────────────────────────────────────────────────────

func addCronJob(ctx context.Context, inputID, name, schedule, command string) (string, error) {
	id := inputID
	if id == "" {
		id = generateID()
	}
	nextRun, err := util.GetNextRun(schedule)
	if err != nil {
		return "", err
	}
	pipe := util.RDB.Pipeline()
	pipe.HSet(ctx, util.JobKey(id), map[string]interface{}{
		"id": id, "name": name, "command": command,
		"schedule": schedule, "paused": "false",
	})
	pipe.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: id})
	if _, err := pipe.Exec(ctx); err != nil {
		return "", err
	}
	util.Logger.Info("Added job",
		slog.String("id", id), slog.String("name", name),
		slog.String("schedule", schedule),
		slog.String("next_run", time.Unix(nextRun, 0).Format(time.RFC3339)),
	)
	return id, nil
}

func deleteCronJob(ctx context.Context, id string) error {
	pipe := util.RDB.Pipeline()
	pipe.ZRem(ctx, util.ZsetKey(), id)
	pipe.ZRem(ctx, util.ProcessingKey(), id)
	pipe.Del(ctx, util.JobKey(id))
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	util.Logger.Info("Deleted job", slog.String("id", id))
	return nil
}

func updateCronJob(ctx context.Context, id, newName, newSchedule, newCommand string) error {
	nextRun, err := util.GetNextRun(newSchedule)
	if err != nil {
		return err
	}
	pipe := util.RDB.Pipeline()
	pipe.HSet(ctx, util.JobKey(id), map[string]interface{}{
		"name": newName, "command": newCommand, "schedule": newSchedule,
	})
	pipe.ZAdd(ctx, util.ZsetKey(), redis.Z{Score: float64(nextRun), Member: id})
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	util.Logger.Info("Updated job",
		slog.String("id", id), slog.String("schedule", newSchedule),
		slog.String("next_run", time.Unix(nextRun, 0).Format(time.RFC3339)),
	)
	return nil
}

func getAllCronJobs(ctx context.Context) ([]map[string]interface{}, error) {
	zItems, err := util.RDB.ZRangeWithScores(ctx, util.ZsetKey(), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(zItems) == 0 {
		return []map[string]interface{}{}, nil
	}
	ids := make([]string, len(zItems))
	scores := make(map[string]float64, len(zItems))
	for i, z := range zItems {
		id := z.Member.(string)
		ids[i] = id
		scores[id] = z.Score
	}
	pipe := util.RDB.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(ids))
	for i, id := range ids {
		cmds[i] = pipe.HGetAll(ctx, util.JobKey(id))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		util.Logger.Error("Pipeline exec failed", slog.Any("error", err))
	}
	jobs := make([]map[string]interface{}, 0, len(ids))
	for i, id := range ids {
		data, _ := cmds[i].Result()
		if len(data) == 0 {
			continue
		}
		jobs = append(jobs, map[string]interface{}{
			"id": id, "name": data["name"], "schedule": data["schedule"],
			"command": data["command"], "paused": data["paused"],
			"next_run": time.Unix(int64(scores[id]), 0).Format(time.RFC3339),
		})
	}
	return jobs, nil
}

func getCronJob(ctx context.Context, id string) (map[string]interface{}, error) {
	pipe := util.RDB.Pipeline()
	hashCmd := pipe.HGetAll(ctx, util.JobKey(id))
	scoreCmd := pipe.ZScore(ctx, util.ZsetKey(), id)
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}
	data, err := hashCmd.Result()
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("job not found")
	}
	result := map[string]interface{}{
		"id": id, "name": data["name"], "schedule": data["schedule"],
		"command": data["command"], "paused": data["paused"],
	}
	if score, err := scoreCmd.Result(); err == nil {
		result["next_run"] = time.Unix(int64(score), 0).Format(time.RFC3339)
	}
	return result, nil
}

func jobExists(ctx context.Context, id string) (bool, error) {
	n, err := util.RDB.Exists(ctx, util.JobKey(id)).Result()
	return n > 0, err
}

// ── Workflow helpers ──────────────────────────────────────────────────────────

func createWorkflow(ctx context.Context, wf *util.Workflow) error {
	if err := util.ValidateWorkflow(wf); err != nil {
		return err
	}

	// Pipeline EXISTS ทุก job_id ใน 1 round trip แทนการ loop N ครั้ง
	pipe := util.RDB.Pipeline()
	existsCmds := make([]*redis.IntCmd, len(wf.Steps))
	for i, step := range wf.Steps {
		existsCmds[i] = pipe.Exists(ctx, util.JobKey(step.JobID))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis error checking jobs: %w", err)
	}
	for i, step := range wf.Steps {
		n, err := existsCmds[i].Result()
		if err != nil {
			return fmt.Errorf("redis error checking job %q: %w", step.JobID, err)
		}
		if n == 0 {
			return fmt.Errorf("step %q references job_id %q which does not exist", step.ID, step.JobID)
		}
	}

	nextRun, err := util.GetNextRun(wf.Schedule)
	if err != nil {
		return err
	}
	fields, err := util.MarshalWorkflow(wf)
	if err != nil {
		return err
	}
	pipe2 := util.RDB.Pipeline()
	pipe2.HSet(ctx, util.WorkflowKey(wf.ID), fields)
	pipe2.ZAdd(ctx, util.WorkflowZsetKey(), redis.Z{Score: float64(nextRun), Member: wf.ID})
	if _, err := pipe2.Exec(ctx); err != nil {
		return err
	}
	util.Logger.Info("Created workflow",
		slog.String("id", wf.ID), slog.String("name", wf.Name),
		slog.String("schedule", wf.Schedule),
		slog.String("next_run", time.Unix(nextRun, 0).Format(time.RFC3339)),
	)
	return nil
}

func deleteWorkflow(ctx context.Context, id string) error {
	pipe := util.RDB.Pipeline()
	pipe.ZRem(ctx, util.WorkflowZsetKey(), id)
	pipe.ZRem(ctx, util.WorkflowRunZsetKey(), id)
	pipe.Del(ctx, util.WorkflowKey(id))
	// history ยังคงอยู่เป็น audit trail
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	util.Logger.Info("Deleted workflow", slog.String("id", id))
	return nil
}

func getWorkflow(ctx context.Context, id string) (*util.Workflow, error) {
	data, err := util.RDB.HGetAll(ctx, util.WorkflowKey(id)).Result()
	if err != nil {
		return nil, err
	}
	return util.UnmarshalWorkflow(data)
}

func workflowExists(ctx context.Context, id string) (bool, error) {
	n, err := util.RDB.Exists(ctx, util.WorkflowKey(id)).Result()
	return n > 0, err
}

func getAllWorkflows(ctx context.Context) ([]map[string]interface{}, error) {
	zItems, err := util.RDB.ZRangeWithScores(ctx, util.WorkflowZsetKey(), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if len(zItems) == 0 {
		return []map[string]interface{}{}, nil
	}

	// Pipeline HGetAll ทุก workflow ใน 1 round trip แทนการ loop N ครั้ง
	pipe := util.RDB.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(zItems))
	scores := make([]float64, len(zItems))
	for i, z := range zItems {
		cmds[i] = pipe.HGetAll(ctx, util.WorkflowKey(z.Member.(string)))
		scores[i] = z.Score
	}
	if _, err := pipe.Exec(ctx); err != nil {
		util.Logger.Error("getAllWorkflows pipeline failed", slog.Any("error", err))
	}

	result := make([]map[string]interface{}, 0, len(zItems))
	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil || len(data) == 0 {
			continue
		}
		wf, err := util.UnmarshalWorkflow(data)
		if err != nil {
			continue
		}
		result = append(result, map[string]interface{}{
			"id":         wf.ID,
			"name":       wf.Name,
			"schedule":   wf.Schedule,
			"paused":     wf.Paused,
			"step_count": len(wf.Steps),
			"on_failure": wf.OnFailure,
			"next_run":   time.Unix(int64(scores[i]), 0).Format(time.RFC3339),
			"created_at": wf.CreatedAt,
		})
	}
	return result, nil
}

// isCronExprError ช่วยตรวจว่า error มาจาก invalid cron expression
func isCronExprError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "invalid cron expression")
}

// ── main ──────────────────────────────────────────────────────────────────────

func main() {
	util.Init()
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// ── Standalone Job routes ─────────────────────────────────────────────────

	app.Post("/cron", func(c *fiber.Ctx) error {
		var input struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Command  string `json:"command"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
		}
		if input.Name == "" || input.Schedule == "" || input.Command == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields: name, schedule, command"})
		}
		id, err := addCronJob(util.CTX, input.ID, input.Name, input.Schedule, input.Command)
		if err != nil {
			if isCronExprError(err) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
			util.Logger.Error("add job failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add cron job"})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "cron job added", "id": id})
	})

	app.Get("/cron", func(c *fiber.Ctx) error {
		jobs, err := getAllCronJobs(util.CTX)
		if err != nil {
			util.Logger.Error("getAll failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch jobs"})
		}
		return c.JSON(jobs)
	})

	app.Get("/cron/:id", func(c *fiber.Ctx) error {
		job, err := getCronJob(util.CTX, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		return c.JSON(job)
	})

	app.Put("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var in struct {
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Command  string `json:"command"`
		}
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
		}
		if in.Schedule == "" || in.Command == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing required fields: schedule, command"})
		}
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := updateCronJob(util.CTX, id, in.Name, in.Schedule, in.Command); err != nil {
			if isCronExprError(err) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
			util.Logger.Error("update failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update job"})
		}
		return c.JSON(fiber.Map{"message": "updated", "id": id})
	})

	app.Delete("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := deleteCronJob(util.CTX, id); err != nil {
			util.Logger.Error("delete failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete job"})
		}
		return c.JSON(fiber.Map{"message": "deleted", "id": id})
	})

	app.Patch("/cron/:id/now", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := util.RDB.ZAdd(util.CTX, util.ZsetKey(), redis.Z{
			Score: float64(time.Now().Unix()), Member: id,
		}).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to schedule"})
		}
		return c.JSON(fiber.Map{"message": "scheduled now", "id": id})
	})

	app.Post("/cron/:id/pause", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := util.RDB.HSet(util.CTX, util.JobKey(id), "paused", "true").Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to pause"})
		}
		return c.JSON(fiber.Map{"message": "paused", "id": id})
	})

	app.Post("/cron/:id/resume", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		if err := util.RDB.HSet(util.CTX, util.JobKey(id), "paused", "false").Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to resume"})
		}
		return c.JSON(fiber.Map{"message": "resumed", "id": id})
	})

	app.Get("/cron/:id/history", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := jobExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		records, err := util.RDB.LRange(util.CTX, util.JobHistoryKey(id), 0, 99).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch history"})
		}
		return c.JSON(fiber.Map{"id": id, "history": records})
	})

	// ── Workflow routes ───────────────────────────────────────────────────────

	// POST /workflow — สร้าง workflow ใหม่
	app.Post("/workflow", func(c *fiber.Ctx) error {
		var input struct {
			ID        string              `json:"id"`
			Name      string              `json:"name"`
			Schedule  string              `json:"schedule"`
			Steps     []util.WorkflowStep `json:"steps"`
			OnFailure string              `json:"on_failure"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
		}

		policy := util.OnFailurePolicy(input.OnFailure)
		if policy == "" {
			policy = util.OnFailureStop
		}
		if policy != util.OnFailureStop && policy != util.OnFailureContinue {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "on_failure must be 'stop' or 'continue'"})
		}

		id := input.ID
		if id == "" {
			id = generateID()
		}
		now := time.Now().UTC().Format(time.RFC3339)
		wf := &util.Workflow{
			ID:        id,
			Name:      input.Name,
			Schedule:  input.Schedule,
			Steps:     input.Steps,
			OnFailure: policy,
			Paused:    false,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := createWorkflow(util.CTX, wf); err != nil {
			if isCronExprError(err) || strings.Contains(err.Error(), "does not exist") ||
				strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "cycle") ||
				strings.Contains(err.Error(), "required") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
			util.Logger.Error("create workflow failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create workflow"})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "workflow created", "id": id})
	})

	// GET /workflow — list workflows
	app.Get("/workflow", func(c *fiber.Ctx) error {
		workflows, err := getAllWorkflows(util.CTX)
		if err != nil {
			util.Logger.Error("list workflows failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch workflows"})
		}
		return c.JSON(workflows)
	})

	// GET /workflow/:id — ดู workflow detail
	app.Get("/workflow/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		wf, err := getWorkflow(util.CTX, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
		}
		// เพิ่ม next_run
		score, _ := util.RDB.ZScore(util.CTX, util.WorkflowZsetKey(), id).Result()
		return c.JSON(fiber.Map{
			"id": wf.ID, "name": wf.Name, "schedule": wf.Schedule,
			"steps": wf.Steps, "on_failure": wf.OnFailure,
			"paused":     wf.Paused,
			"next_run":   time.Unix(int64(score), 0).Format(time.RFC3339),
			"created_at": wf.CreatedAt,
			"updated_at": wf.UpdatedAt,
		})
	})

	// DELETE /workflow/:id — ลบ workflow
	app.Delete("/workflow/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := workflowExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
		}
		if err := deleteWorkflow(util.CTX, id); err != nil {
			util.Logger.Error("delete workflow failed", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete workflow"})
		}
		return c.JSON(fiber.Map{"message": "deleted", "id": id})
	})

	// POST /workflow/:id/pause — หยุด workflow
	app.Post("/workflow/:id/pause", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := workflowExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
		}
		if err := util.RDB.HSet(util.CTX, util.WorkflowKey(id), "paused", "true", "updated_at",
			time.Now().UTC().Format(time.RFC3339)).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to pause"})
		}
		return c.JSON(fiber.Map{"message": "paused", "id": id})
	})

	// POST /workflow/:id/resume — เปิด workflow อีกครั้ง
	app.Post("/workflow/:id/resume", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := workflowExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
		}
		if err := util.RDB.HSet(util.CTX, util.WorkflowKey(id), "paused", "false", "updated_at",
			time.Now().UTC().Format(time.RFC3339)).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to resume"})
		}
		return c.JSON(fiber.Map{"message": "resumed", "id": id})
	})

	// PATCH /workflow/:id/now — trigger workflow ทันที
	app.Patch("/workflow/:id/now", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := workflowExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
		}
		if err := util.RDB.ZAdd(util.CTX, util.WorkflowZsetKey(), redis.Z{
			Score: float64(time.Now().Unix()), Member: id,
		}).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to schedule"})
		}
		util.Logger.Info("workflow scheduled for immediate run", slog.String("id", id))
		return c.JSON(fiber.Map{"message": "scheduled now", "id": id})
	})

	// GET /workflow/:id/history — ดู execution history ของ workflow
	app.Get("/workflow/:id/history", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if ok, _ := workflowExists(util.CTX, id); !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workflow not found"})
		}
		records, err := util.RDB.LRange(util.CTX, util.WorkflowHistoryKey(id), 0, 49).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch history"})
		}
		return c.JSON(fiber.Map{"id": id, "history": records})
	})

	// GET /workflow/:id/runs/:run_id — ดู step-level detail ของ run หนึ่งครั้ง
	app.Get("/workflow/:id/runs/:run_id", func(c *fiber.Ctx) error {
		wfID := c.Params("id")
		runID := c.Params("run_id")
		data, err := util.RDB.HGetAll(util.CTX, util.WorkflowRunKey(wfID, runID)).Result()
		if err != nil || len(data) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "run not found"})
		}
		state, err := util.UnmarshalRunState(data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse run state"})
		}
		return c.JSON(state)
	})

	// graceful shutdown
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		util.Logger.Info("shutting down api server...")
		_ = app.Shutdown()
	}()

	port := util.GetEnv("API_PORT", "3000")
	util.Logger.Info("starting api server", slog.String("port", port))
	if err := app.Listen(":" + port); err != nil {
		util.Logger.Error("server error", slog.Any("error", err))
	}
}
