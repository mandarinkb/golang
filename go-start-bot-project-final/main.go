package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-start-bot-project-final/database"
	"github.com/mandarinkb/go-start-bot-project-final/repository"
	"github.com/mandarinkb/go-start-bot-project-final/service"
	"github.com/robfig/cron/v3"
)

var (
	Cron    *cron.Cron
	CronId  cron.EntryID
	cronStr string
)

// ตัดตัวอักษรตัวแรกออก
func trimFirstPrefix(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

// ตัด 2 ตัวอักษรข้างหน้าออก
func trimTwoLetterPrefix(cron string) string {
	var newCron string
	newCron = trimFirstPrefix(cron)
	newCron = trimFirstPrefix(newCron)
	return newCron
}

// สั่งให้ bot ทำงาน
func startBot() {
	db, err := database.Conn()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	webRepo := repository.NewWeb(db)
	err = service.NewWebService(webRepo).Read()
	if err != nil {
		fmt.Print(err)
	}
}

// ดึงค่าตั้งเวาการทำงานจาก database
func getCronExpression() {
	db, err := database.Conn()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	scheduleRepo := repository.NewScheduleRepo(db)
	scheduleServ, err := service.NewScheduleService(scheduleRepo).Read()
	if err != nil {
		fmt.Println(err)
	}
	var cr string
	for _, row := range scheduleServ {
		if row.ProjectName == "project-final-start-bot" {
			cr = row.CronExpression
		}
	}
	cronStr = trimTwoLetterPrefix(cr)
}

// ตั้งเวาการทำงาน
func cronJob() {
	getCronExpression()
	fmt.Println(cronStr)
	Cron = cron.New()
	CronId, _ = Cron.AddFunc(cronStr, func() {
		startBot()
	})
	Cron.Start()
}

func restartHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		Cron.Remove(CronId)
		cronJob()
	}
}

func main() {
	cronJob()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/restart", restartHandler())
	router.Run(":8081")
}
