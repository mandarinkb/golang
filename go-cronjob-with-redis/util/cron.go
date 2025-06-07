package util

import (
	"log"
	"time"

	"github.com/robfig/cron"
)

// ใช้ Cron Parser คำนวณเวลาถัดไป
func GetNextRun(schedule string) int64 {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	expr, err := parser.Parse(schedule)
	if err != nil {
		log.Fatal("Invalid cron expression:", err)
	}
	return expr.Next(time.Now()).Unix()
}
