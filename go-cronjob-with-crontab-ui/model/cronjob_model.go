package model

import (
	"encoding/json"

	"github.com/robfig/cron/v3"
)

var (
	Cron *cron.Cron = cron.New()

	// store cron job detail
	CronJobDetail []CronJobData

	// store cron job case status pending
	// and remove after excute success
	CronJobDetailForRemove []CronJobData

	StatusExcuteServiceOne string
	StatusExcuteServiceTwo string
)

type CronJobData struct {
	CronEntryID      cron.EntryID
	CronIDRef        string
	CronName         string
	CronExpression   string
	CronStopped      bool
	IsCronDataChange bool
}

func (c *CronJobData) Marshal() string {
	js, _ := json.Marshal(c)
	return string(js)
}
