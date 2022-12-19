package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-cronjob-with-crontab-ui/config"
	"github.com/go-cronjob-with-crontab-ui/logger"
	"github.com/go-cronjob-with-crontab-ui/model"
	"github.com/go-cronjob-with-crontab-ui/repository"
	"github.com/go-cronjob-with-crontab-ui/utils"
)

type cronJobService struct {
	cronJobRepo repository.CronJobRepository
	log         *logger.Logger
}

func NewCronJobService(cronJobRepo repository.CronJobRepository) CronJobService {
	return &cronJobService{
		cronJobRepo: cronJobRepo,
		log:         logger.L().Named("CronJobService")}
}

func (c *cronJobService) RunJobService(ctx context.Context) error {
	isFirstExcute := true
	file := config.C().Crontab.Path

	//for check error read file
	_, err := utils.ReadFile(file)
	if err != nil {
		c.log.WithContext(ctx).Debugf("Read file error : ", err)
		return err
	}
	callback := func() {
		ctxWithCorrID := utils.MakeContextCorrelationID(ctx)
		fileLines, _ := utils.ReadFile(file)

		c.syncData(ctxWithCorrID, fileLines, model.CronJobDetail)
		c.log.WithContext(ctxWithCorrID).Infof("cron detail : ", model.CronJobDetail)

		if model.CronJobDetail != nil {
			isFirstExcute = false
			c.runJobWhenEventChange(ctxWithCorrID)
			fmt.Println("CronJobDetail : ", model.CronJobDetail)
		}
		if isFirstExcute {
			fmt.Println("is first start")
			c.runJobWhenServiceStart(ctxWithCorrID, fileLines)
		}

		c.removeStatusServicePending(ctxWithCorrID)
	}
	// call callback function
	callback()

	model.Cron.Start()

	utils.ReadFileSystemEventChange(file, callback)

	return nil
}

func (c *cronJobService) syncData(ctx context.Context, fileLines []model.Crontab, cronDetail []model.CronJobData) {
	for _, line := range fileLines {
		for i, inMemoryCron := range cronDetail {
			if line.ID == inMemoryCron.CronIDRef {
				model.CronJobDetail = utils.RemoveElement(model.CronJobDetail, i)
				cronData := model.CronJobData{
					CronIDRef:        line.ID,
					CronName:         line.Name,
					CronExpression:   line.Schedule,
					CronStopped:      line.Stopped,
					IsCronDataChange: false,
					CronEntryID:      inMemoryCron.CronEntryID,
				}
				// check cronexpression change
				if line.Schedule != inMemoryCron.CronExpression || line.Stopped != inMemoryCron.CronStopped {
					cronData.IsCronDataChange = true
				}
				model.CronJobDetail = append(model.CronJobDetail, cronData)
				break
			}
		}
	}

}

func (c *cronJobService) runJobWhenServiceStart(ctx context.Context, fileLines []model.Crontab) {
	for _, line := range fileLines {
		cronData := model.CronJobData{
			CronIDRef:      line.ID,
			CronName:       line.Name,
			CronExpression: line.Schedule,
			CronStopped:    line.Stopped,
		}
		if !line.Stopped {
			switch line.Name {
			case config.C().CronJobService.ServiceOne:
				cronData.CronEntryID, _ = c.cronJobRepo.RunJob(line.Schedule, c.serviceOne)
			case config.C().CronJobService.ServiceTwo:
				cronData.CronEntryID, _ = c.cronJobRepo.RunJob(line.Schedule, c.serviceTwo)
			}
		}
		model.CronJobDetail = append(model.CronJobDetail, cronData)
	}
}

func (c *cronJobService) runJobWhenEventChange(ctx context.Context) {
	for i, inMemoryCron := range model.CronJobDetail {
		if inMemoryCron.IsCronDataChange {
			switch inMemoryCron.CronName {
			case config.C().CronJobService.ServiceOne:
				if model.StatusExcuteServiceOne != "pending" {
					if inMemoryCron.CronStopped {
						c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
					} else {
						c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
						model.CronJobDetail = utils.RemoveElement(model.CronJobDetail, i)
						inMemoryCron.CronEntryID, _ = c.cronJobRepo.RunJob(inMemoryCron.CronExpression, c.serviceOne)
						model.CronJobDetail = append(model.CronJobDetail, inMemoryCron)
					}
				} else {
					model.CronJobDetailForRemove = append(model.CronJobDetailForRemove, inMemoryCron)
					fmt.Println("service one is running")
					fmt.Println("cronJob detail for remove", model.CronJobDetailForRemove)
				}
			case config.C().CronJobService.ServiceTwo:
				if model.StatusExcuteServiceTwo != "pending" {
					if inMemoryCron.CronStopped {
						c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
					} else {
						c.cronJobRepo.RemoveJob(inMemoryCron.CronEntryID)
						model.CronJobDetail = utils.RemoveElement(model.CronJobDetail, i)
						inMemoryCron.CronEntryID, _ = c.cronJobRepo.RunJob(inMemoryCron.CronExpression, c.serviceTwo)
						model.CronJobDetail = append(model.CronJobDetail, inMemoryCron)
					}
				} else {
					model.CronJobDetailForRemove = append(model.CronJobDetailForRemove, inMemoryCron)
					fmt.Println("service two is running")
					fmt.Println("cronJob detail for remove", model.CronJobDetailForRemove)
				}
			}
		}
	}
}

func (c *cronJobService) removeStatusServicePending(ctx context.Context) {
	for i, inMemoryCronRemove := range model.CronJobDetailForRemove {
		c.cronJobRepo.RemoveJob(inMemoryCronRemove.CronEntryID)
		utils.RemoveElement(model.CronJobDetailForRemove, i)
		fmt.Println("remove cronjob status pending : ", inMemoryCronRemove.Marshal())
	}
}

func (c *cronJobService) serviceOne() {
	model.StatusExcuteServiceOne = "pending"
	flag := true
	count := 0
	for flag {

		if count >= 50 {
			flag = false
		}

		fmt.Println(utils.CurrentLocalDate(), "serviceOne running : ", count)
		count++
		time.Sleep(1 * time.Second)
	}

	model.StatusExcuteServiceOne = ""
}

func (c *cronJobService) serviceTwo() {
	model.StatusExcuteServiceTwo = "pending"
	fmt.Println("======================")
	fmt.Println(utils.CurrentLocalDate(), "serviceTwo running")
	fmt.Println("======================")
	model.StatusExcuteServiceTwo = ""
}
