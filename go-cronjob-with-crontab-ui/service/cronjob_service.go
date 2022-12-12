package service

import (
	"context"
	"fmt"

	"github.com/go-cronjob-with-crontab-ui/config"
	"github.com/go-cronjob-with-crontab-ui/model"
	"github.com/go-cronjob-with-crontab-ui/repository"
	"github.com/go-cronjob-with-crontab-ui/utils"
)

var statusRunning string

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
	isRunning := true
	callback := func() {
		if statusRunning != "pending" {
			fileLines, _ := utils.ReadFile(file)

			syncChange(fileLines, repository.CronJobDetail)

			if repository.CronJobDetail != nil {
				isRunning = false
				for i, inMemoryCron := range repository.CronJobDetail {
					if inMemoryCron.IsCronDataChange {
						switch inMemoryCron.CronIDRef {
						case "6ASM7cWbuEL2WaPA": // demo
							if inMemoryCron.CronStopped {
								c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
							} else {
								c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
								repository.CronJobDetail = utils.RemoveElement(repository.CronJobDetail, i)

								id, _ := c.cronJobRepo.RunJob(inMemoryCron.CronExpression, start)
								inMemoryCron.CronEntryID = id

								repository.CronJobDetail = append(repository.CronJobDetail, inMemoryCron)
							}
						case "v0IqBRBhIGkzDyUM":
							if inMemoryCron.CronStopped {
								c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
							} else {
								c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
								repository.CronJobDetail = utils.RemoveElement(repository.CronJobDetail, i)

								id, _ := c.cronJobRepo.RunJob(inMemoryCron.CronExpression, start2)
								inMemoryCron.CronEntryID = id

								repository.CronJobDetail = append(repository.CronJobDetail, inMemoryCron)
							}
						}
					}
				}

				fmt.Println("CronJobDetail : ", repository.CronJobDetail)
			}
			if isRunning {
				fmt.Println("is first start")

				for _, line := range fileLines {
					cronData := repository.CronJobData{
						CronIDRef:      line.Id,
						CronName:       line.Name,
						CronExpression: line.Schedule,
						CronStopped:    line.Stopped,
					}
					if !line.Stopped {
						if line.Name == "demo" {
							id, _ := c.cronJobRepo.RunJob(line.Schedule, start)
							cronData.CronEntryID = id
						}
						if line.Name == "2m" {
							id, _ := c.cronJobRepo.RunJob(line.Schedule, start2)
							cronData.CronEntryID = id
						}
					}
					repository.CronJobDetail = append(repository.CronJobDetail, cronData)
				}
			}

		}

	}
	// call function
	callback()

	repository.Cron.Start()

	utils.ReadFileSystemEventChange(file, callback)

	return nil
}

func start() {
	statusRunning = "pending"
	fmt.Println(utils.CurrentLocalDate(), "hello world")
	statusRunning = ""
}

func start2() {
	statusRunning = "pending"
	fmt.Println(utils.CurrentLocalDate(), "================hello world2")
	statusRunning = ""
}

func syncChange(fileLines []model.Crontab, cronDetail []repository.CronJobData) {
	for _, line := range fileLines {
		for i, inMemoryCron := range cronDetail {
			if line.Id == inMemoryCron.CronIDRef {
				repository.CronJobDetail = utils.RemoveElement(repository.CronJobDetail, i)
				cronData := repository.CronJobData{
					CronIDRef:        line.Id,
					CronName:         line.Name,
					CronExpression:   line.Schedule,
					CronStopped:      line.Stopped,
					IsCronDataChange: false,
					CronEntryID:      inMemoryCron.CronEntryID, // original
				}
				fmt.Println("inMemoryCron.CronEntryID: ", inMemoryCron, " : ", inMemoryCron.CronEntryID)
				// check cronexpression change
				if line.Schedule != inMemoryCron.CronExpression || line.Stopped != inMemoryCron.CronStopped {
					cronData.IsCronDataChange = true
				}
				repository.CronJobDetail = append(repository.CronJobDetail, cronData)
				break
			}
		}
	}

}
