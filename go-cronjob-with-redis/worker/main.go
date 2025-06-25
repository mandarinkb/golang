package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mandarinkb/go-cronjob-with-redis/util"
	"github.com/redis/go-redis/v9"
)

func processCronJobs(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second) // ตรวจสอบทุก 2 วินาที
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping cron job processor gracefully...")
			return

		case <-ticker.C:
			now := time.Now().Unix()

			jobs, err := util.RDB.ZRangeByScore(util.CTX, "cron_jobs", &redis.ZRangeBy{
				Min: "-inf", Max: fmt.Sprintf("%d", now),
			}).Result()
			if err != nil {
				log.Println("[ERROR] Fetching jobs:", err)
				continue
			}

			if len(jobs) == 0 {
				log.Println("[INFO] No jobs to run at", time.Now().Format(time.RFC3339))
				continue
			}

			for _, jobID := range jobs {
				job, err := util.RDB.HGetAll(util.CTX, "job_name-"+jobID).Result()
				if err != nil || len(job) == 0 {
					log.Printf("[WARN] Job data missing for ID: %s\n", jobID)
					continue
				}

				if job["paused"] == "true" {
					log.Printf("[SKIP] Job %s is paused.\n", jobID)
					continue
				}

				log.Printf("[RUN] ID: %s | Command: %s\n", jobID, job["command"])

				nextRun := util.GetNextRun(job["schedule"])
				util.RDB.ZRem(util.CTX, "cron_jobs", jobID)
				util.RDB.ZAdd(util.CTX, "cron_jobs", redis.Z{
					Score:  float64(nextRun),
					Member: jobID,
				})

				log.Printf("[SCHEDULED] Next run at: %s\n", time.Unix(nextRun, 0).Format(time.RFC3339))
			}
		}
	}
}

func main() {
	log.Println("Starting cron job worker...")

	ctx, cancel := context.WithCancel(context.Background())

	// ตั้งค่าให้กด Ctrl+C เพื่อ stop worker
	// สร้าง channel เพื่อรอรับสัญญาณจาก OS เช่น Ctrl+C (SIGINT) หรือ docker stop (SIGTERM)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// เมื่อได้รับสัญญาณ → เรียก cancel() เพื่อแจ้งให้ processCronJobs() หยุดทำงาน
	go func() {
		<-sigs
		log.Println("Received stop signal")
		cancel()
	}()

	// เรียกให้ worker ทำงาน โดยส่ง context ที่เราสามารถ cancel ได้
	processCronJobs(ctx)
}
