package repository

import (
	"github.com/robfig/cron/v3"
)

type cronJobRepository struct{}

func NewCronJobRepository() CronJobRepository {
	return &cronJobRepository{}
}

func (c *cronJobRepository) RunJob(cronExpression string, cmd func()) (cID cron.EntryID, err error) {
	return Cron.AddFunc(cronExpression, cmd)
}

func (c *cronJobRepository) RemoveJob(id cron.EntryID) {
	Cron.Remove(id)
}
