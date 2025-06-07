package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mandarinkb/go-cronjob-with-redis/util"
	"github.com/redis/go-redis/v9"
	// ใช้ Cron Parser
)

// สร้าง UUID แบบไม่มี dash
func generateID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// เพิ่ม Cron Job เข้า Redis
func addCronJob(inputID, namme, schedule, command string) (string, error) {
	id := inputID
	if inputID == "" {
		id = generateID()
	}

	nextRun := util.GetNextRun(schedule)

	// เก็บ Job ลง Redis
	if err := util.RDB.ZAdd(util.CTX, "cron_jobs", redis.Z{Score: float64(nextRun), Member: id}).Err(); err != nil {
		return "", err
	}
	if err := util.RDB.HSet(util.CTX, "job_name-"+id, map[string]interface{}{
		"command":  command,
		"schedule": schedule,
	}).Err(); err != nil {
		return "", err
	}

	fmt.Println("Added job:", id, "Next run at:", time.Unix(nextRun, 0))
	return id, nil
}

// ลบ Cron Job ออกจาก Redis
func deleteCronJob(id string) error {
	// ลบจาก ZSet
	if err := util.RDB.ZRem(util.CTX, "cron_jobs", id).Err(); err != nil {
		return err
	}
	// ลบข้อมูลรายละเอียด
	if err := util.RDB.Del(util.CTX, "job_name-"+id).Err(); err != nil {
		return err
	}
	fmt.Println("Deleted job:", id)
	return nil
}

// อัปเดต Cron Job ใน Redis
func updateCronJob(id, newName, newSchedule, newCommand string) error {
	// ลบ Job เก่า
	if err := deleteCronJob(id); err != nil {
		return err
	}
	// เพิ่ม Job ใหม่
	_, err := addCronJob(id, newName, newSchedule, newCommand)
	if err != nil {
		return err
	}

	return nil
}

// ดึง Job ทั้งหมดจาก Redis
func getAllCronJobs() ([]map[string]interface{}, error) {
	// ดึง IDs ทั้งหมดจาก ZSet
	ids, err := util.RDB.ZRange(util.CTX, "cron_jobs", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var jobs []map[string]interface{}
	for _, id := range ids {
		data, err := util.RDB.HGetAll(util.CTX, "job_name-"+id).Result()
		if err != nil {
			continue
		}
		if len(data) == 0 {
			continue
		}
		job := map[string]interface{}{
			"id":       id,
			"name":     data["name"],
			"schedule": data["schedule"],
			"command":  data["command"],
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// ดึง Job ตาม ID จาก Redis
func getCronJob(id string) (map[string]interface{}, error) {
	data, err := util.RDB.HGetAll(util.CTX, "job_name-"+id).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("job not found")
	}
	job := map[string]interface{}{
		"id":       id,
		"name":     data["name"],
		"schedule": data["schedule"],
		"command":  data["command"],
	}
	return job, nil
}

func main() {
	app := fiber.New()

	app.Post("/cron", func(c *fiber.Ctx) error {
		type CronJobInput struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Command  string `json:"command"`
		}

		var input CronJobInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		if input.Name == "" || input.Schedule == "" || input.Command == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing required fields",
			})
		}

		id, err := addCronJob(input.ID, input.Name, input.Schedule, input.Command)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add cron job",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Cron job added successfully",
			"id":      id,
		})
	})

	app.Delete("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID is required",
			})
		}

		err := deleteCronJob(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete cron job",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Cron job deleted successfully",
			"id":      id,
		})
	})

	app.Put("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID is required",
			})
		}

		type UpdateCronJobInput struct {
			Name     string `json:"name"`
			Schedule string `json:"schedule"`
			Command  string `json:"command"`
		}

		var input UpdateCronJobInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		if input.Schedule == "" || input.Command == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing required fields",
			})
		}

		err := updateCronJob(id, input.Name, input.Schedule, input.Command)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update cron job",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Cron job updated successfully",
			"id":      id,
		})
	})

	app.Get("/cron", func(c *fiber.Ctx) error {
		jobs, err := getAllCronJobs()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch cron jobs",
			})
		}
		return c.JSON(jobs)
	})

	app.Get("/cron/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID is required",
			})
		}

		job, err := getCronJob(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cron job not found",
			})
		}

		return c.JSON(job)
	})

	fmt.Println("Starting server on port 3000")
	app.Listen(":3000")
}
