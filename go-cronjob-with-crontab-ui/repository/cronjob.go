package repository

import "github.com/robfig/cron/v3"

type CronJobRepository interface {
	RunJob(cronExpression string, cmd func()) (cron.EntryID, error)
	RemoveJob(id cron.EntryID)
}
