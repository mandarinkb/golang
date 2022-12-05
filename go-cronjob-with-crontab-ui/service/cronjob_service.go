package service

import (
	"context"
	"fmt"

	"github.com/go-cronjob-with-crontab-ui/config"
	"github.com/go-cronjob-with-crontab-ui/repository"
	"github.com/go-cronjob-with-crontab-ui/utils"
)

type cronJobService struct {
	cronJobRepo repository.CronJobRepository
}

func NewCronJobService(cronJobRepo repository.CronJobRepository) CronJobService {
	return &cronJobService{cronJobRepo}
}

func (c *cronJobService) RunJobService(ctx context.Context) error {
	file := config.C().Crontab.Path

	//for check error read file
	_, err := utils.ReadFile(file)
	if err != nil {
		return err
	}

	callback := func() {

		fileLines, _ := utils.ReadFile(file)
		fmt.Println(fileLines)

		if len(repository.CronJobDetail) > 0 {
			// remove old cronjob
			for i, v := range repository.CronJobDetail {
				if v.CronName == "demo" {
					c.cronJobRepo.RemoveJob(v.CronEntryID)
					repository.CronJobDetail = utils.RemoveElement(repository.CronJobDetail, i)
				}
			}
		}

		for _, v := range fileLines {
			if v.Name == "demo" {
				id, _ := c.cronJobRepo.RunJob(v.Schedule, start)
				cronData := repository.CronJobData{
					CronEntryID:    id,
					CronIDRef:      v.Id,
					CronName:       v.Name,
					CronExpression: v.Schedule,
				}
				repository.CronJobDetail = append(repository.CronJobDetail, cronData)
			}
		}
		fmt.Println(repository.CronJobDetail)

	}
	// call function
	callback()

	repository.Cron.Start()

	utils.ReadFileSystemEventChange(file, callback)

	return nil
}

func start() {
	fmt.Println(utils.CurrentLocalDate(), "hello world")
}
