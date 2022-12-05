package repository

import "github.com/robfig/cron/v3"

var (
	Cron *cron.Cron = cron.New()
	//store cron job detail
	CronJobDetail []CronJobData
)

type CronJobData struct {
	CronEntryID    cron.EntryID
	CronIDRef      string
	CronName       string
	CronExpression string
}

type CronJobRepository interface {
	RunJob(cronExpression string, cmd func()) (cron.EntryID, error)
	RemoveJob(id cron.EntryID)
}
