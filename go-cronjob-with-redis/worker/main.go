package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mandarinkb/go-cronjob-with-redis/util"
	"github.com/redis/go-redis/v9"
)

func processCronJobs() {
	for {
		now := time.Now().Unix()

		// ดึง Job ที่ถึงเวลาแล้ว
		jobs, err := util.RDB.ZRangeByScore(util.CTX, "cron_jobs", &redis.ZRangeBy{
			Min: "-inf", Max: fmt.Sprintf("%d", now),
		}).Result()
		if err != nil {
			fmt.Println(err)
			break
		}

		log.Print("jobs:", jobs)
		for _, jobID := range jobs {
			job, _ := util.RDB.HGetAll(util.CTX, "job_name-"+jobID).Result()
			fmt.Println("Running:", job["command"])

			// อัปเดตเวลาถัดไป
			nextRun := util.GetNextRun(job["schedule"])
			util.RDB.ZRem(util.CTX, "cron_jobs", jobID) // ลบออกก่อน
			util.RDB.ZAdd(util.CTX, "cron_jobs", redis.Z{Score: float64(nextRun), Member: jobID})

			fmt.Println("Next run at:", time.Unix(nextRun, 0))
		}

		time.Sleep(5 * time.Second) // ตรวจสอบทุก 5 วินาที
	}
}

func main() {
	processCronJobs() // เริ่ม Worker
}
