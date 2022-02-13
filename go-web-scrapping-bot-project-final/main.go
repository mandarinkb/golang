package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mandarinkb/go-web-scrapping-bot-project-final/database"
	"github.com/mandarinkb/go-web-scrapping-bot-project-final/repository"
	"github.com/mandarinkb/go-web-scrapping-bot-project-final/service"
	"github.com/mandarinkb/go-web-scrapping-bot-project-final/utils"
	"github.com/robfig/cron/v3"
)

// ตั้งเวาการทำงาน โดยทำงาน ทุกๆ 1 นาที
func cronJob() {
	c := cron.New()
	c.AddFunc("0/1 * 1/1 * ?", func() {
		run()
	})
	c.Start()
}
func run() {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	db, err := database.Conn()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	defer db.Close()

	logger.Info("[web scrapping bot] start", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))
	fmt.Println(time.Now(), "web scrapping bot start")

	redis := database.RedisConn()
	defer redis.Close()
	webData := service.Web{}

	checkStartUrl := true
	for checkStartUrl {
		detail, err := redis.RPop(context.Background(), "detailUrl").Result() //detail
		if err != nil {
			checkStartUrl = false
		}
		if detail != "" {
			json.Unmarshal([]byte(detail), &webData)
			if webData.WebStatus == "1" {
				switch webData.WebName {
				case "tescolotus":
					swRepo := repository.NewSwitchDatabaseDB(db)
					service.NewProductService(swRepo).Tescolotus(webData)
				case "makroclick":
					service.Makroclick(webData)
				case "bigc":
					service.Bigc(webData)
				}
			}
		}
	}
	logger.Info("[web scrapping bot] stop", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))
	fmt.Println(time.Now(), "web scrapping bot stop")
}

func main() {
	cronJob()
	// ต้องมีคำสั่งนี้ ไม่งั้นฟังก์ชันการตั้งเวลาทำงานจะไม่ทำงาน เพราะต้องใช้ goroutine
	for {
		time.Sleep(time.Second)
	}
}
