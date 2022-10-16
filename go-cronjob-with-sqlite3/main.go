package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"

	"github.com/mandarinkb/go-cronjob-with-sqlite3/assets"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/config"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/logger"
	"github.com/mandarinkb/go-cronjob-with-sqlite3/repository"
)

// var wg sync.WaitGroup

func job(exp string) {
	//wg.Add(1)
	c := cron.New()
	entryID, _ := c.AddFunc(exp, start)
	fmt.Println("entry ID: ", entryID)
	c.Start()
	// ต้องมีคำสั่งนี้ ไม่งั้นฟังก์ชันการตั้งเวลาทำงานจะไม่ทำงาน เพราะต้องใช้ goroutine
	select {}
	//wg.Wait()
}
func start() {
	fmt.Println("start running...")
	fmt.Println(time.Now())
}
func main() {
	config.LoadConfig("config", "config")
	assets.LoadAssets("assets", "error")
	mainLog := logger.InitialLogger()
	mainLog.Info("main log")

	// Open the created SQLite File
	sqliteDB, err := sql.Open("sqlite3", "./sqlite/cronjob.db")
	if err != nil {
		mainLog.Errorf("error connect sqlite3 error: %v", err)
		panic(err)
	}
	defer sqliteDB.Close()

	cronJobRepo := repository.NewCronJobRepository(sqliteDB)
	cronJobData, err := cronJobRepo.Read()
	if err != nil {
		mainLog.Error(err.Error())
	}
	fmt.Println(cronJobData)
	for _, v := range cronJobData {
		if v.CronUseInFunc == "start" {
			job(v.CronExpression) //"0/1 * 1/1 * ?"
		}
	}

}
