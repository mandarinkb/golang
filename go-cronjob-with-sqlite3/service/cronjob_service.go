package service

import "github.com/mandarinkb/go-cronjob-with-sqlite3/repository"

type cronJobService struct {
	cronJobRepo repository.CronJobRepository
}

func NewCronJobService(cronJobRepo repository.CronJobRepository) CronJobService {
	return &cronJobService{cronJobRepo}
}

func (s *cronJobService) GetCronJob() (cronjob []repository.CronJob, err error) {
	return s.cronJobRepo.GetCronJob()
}

func (s *cronJobService) GetCronJobByID(id int) (*repository.CronJob, error) {
	return s.cronJobRepo.GetCronJobByID(id)
}

func (s *cronJobService) CreateCronJob(cronjob repository.CronJob) error {
	return s.cronJobRepo.CreateCronJob(cronjob)
}

func (s *cronJobService) UpdateCronJob(cronjob repository.CronJob) error {
	return s.cronJobRepo.UpdateCronJob(cronjob)
}

func (s *cronJobService) DeleteCronJob(id int) error {
	return s.cronJobRepo.DeleteCronJob(id)
}
