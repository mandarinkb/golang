package main

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"

	"github.com/mandarinkb/go-cronjob-with-sqlite3/assets"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/config"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/handler"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/logger"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/repository"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/service"
)

func main() {
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	mainLog := logger.InitialLogger()
	mainLog.Info("main log")

	// Open the created SQLite File
	sqliteDB, err := sql.Open("sqlite3", "./sqlite/cronjob.db")
	if err != nil {
		panic(err)
	}
	defer sqliteDB.Close()

	cronJobRepo := repository.NewCronJobRepository(sqliteDB)
	cronJobServ := service.NewCronJobService(cronJobRepo)
	cronJobHandler := handler.NewCronJobHandler(cronJobServ)

	app := fiber.New()
	app.Get("/api/cronjob", cronJobHandler.GetCronJob)
	app.Listen(":3000")

	// cronJobRepo := repository.NewCronJobRepository(sqliteDB)

	// err = cronJobRepo.DeleteCronJob(2)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// cronJobData, err := cronJobRepo.GetCronJob()
	// if err != nil {
	// 	mainLog.Error(err.Error())
	// }
	// fmt.Println(cronJobData)

	// for _, v := range cronJobData {
	// 	if v.CronFunctionRef == "ex_cronjob_func" {
	// 		id, _ := cronJobRepo.RunJob(v.CronExpression, start)
	// 		cronData := repository.CronJobData{
	// 			CronEntryID:     id,
	// 			CronFunctionRef: v.CronName,
	// 		}
	// 		repository.CronJobDetail = append(repository.CronJobDetail, cronData)

	// 		id2, _ := cronJobRepo.RunJob("0/2 * 1/1 * ?", start2)
	// 		cronData.CronEntryID = id2
	// 		cronData.CronFunctionRef = "every_2m"
	// 		repository.CronJobDetail = append(repository.CronJobDetail, cronData)

	// 		id3, _ := cronJobRepo.RunJob("0/3 * 1/1 * ?", start3)
	// 		cronData.CronEntryID = id3
	// 		cronData.CronFunctionRef = "every_3m"
	// 		repository.CronJobDetail = append(repository.CronJobDetail, cronData)
	// 	}
	// }
	// fmt.Println(repository.CronJobDetail)
	// for i, v := range repository.CronJobDetail {
	// 	if v.CronFunctionRef == "every_2m" {
	// 		cronJobRepo.RemoveJob(v.CronEntryID)
	// 		repository.CronJobDetail = removeElement(repository.CronJobDetail, i)
	// 	}
	// }
	// fmt.Println("remove : ", repository.CronJobDetail)
	// repository.Cron.Start()
	// // ต้องมีคำสั่งนี้ ไม่งั้นฟังก์ชันการตั้งเวลาทำงานจะไม่ทำงาน เพราะต้องใช้ goroutine
	// select {}
}
