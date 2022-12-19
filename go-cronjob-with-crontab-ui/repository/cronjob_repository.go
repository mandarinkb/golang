package repository

import (
	"github.com/go-cronjob-with-crontab-ui/model"
	"github.com/robfig/cron/v3"
)

type cronJobRepository struct{}

func NewCronJobRepository() CronJobRepository {
	return &cronJobRepository{}
}

func (c *cronJobRepository) RunJob(cronExpression string, cmd func()) (cron.EntryID, error) {
	return model.Cron.AddFunc(cronExpression, cmd)
}

func (c *cronJobRepository) RemoveJob(id cron.EntryID) {
	model.Cron.Remove(id)
}
