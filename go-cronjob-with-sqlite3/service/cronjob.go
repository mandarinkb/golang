package service

import "github.com/mandarinkb/go-cronjob-with-sqlite3/repository"

type CronJobService interface {
	GetCronJob() (cronjob []repository.CronJob, err error)
	GetCronJobByID(id int) (*repository.CronJob, error)
	CreateCronJob(cronjob repository.CronJob) error
	UpdateCronJob(cronjob repository.CronJob) error
	DeleteCronJob(id int) error
}
