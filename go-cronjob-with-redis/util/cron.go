package util

import (
	"log"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// GetNextRun คำนวณเวลาถัดไปจาก schedule
func GetNextRun(schedule string) int64 {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	expr, err := parser.Parse(schedule)
	if err != nil {
		Logger.Error("Invalid cron expression", slog.String("schedule", schedule), slog.Any("error", err))
		log.Fatal(err)
	}
	return expr.Next(time.Now()).Unix()
}
