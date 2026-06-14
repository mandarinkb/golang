package util

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// GetNextRun คำนวณเวลาถัดไปจาก schedule
// return error แทน log.Fatal เพื่อให้ caller จัดการได้
func GetNextRun(schedule string) (int64, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	expr, err := parser.Parse(schedule)
	if err != nil {
		Logger.Error("Invalid cron expression", slog.String("schedule", schedule), slog.Any("error", err))
		return 0, fmt.Errorf("invalid cron expression %q: %w", schedule, err)
	}
	return expr.Next(time.Now()).Unix(), nil
}
