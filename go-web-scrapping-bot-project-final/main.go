package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mandarinkb/go-web-scrapping-bot-project-final/database"
	"github.com/mandarinkb/go-web-scrapping-bot-project-final/repository"
	"github.com/mandarinkb/go-web-scrapping-bot-project-final/service"
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
	fmt.Println(time.Now(), "webscrapping bot start")
	redis := database.RedisConn()
	defer redis.Close()
	rWeb := service.Web{}

	checkStartUrl := true
	for checkStartUrl {
		detail, err := redis.RPop(context.Background(), "detailUrl").Result() //detail
		if err != nil {
			fmt.Println(err)
			checkStartUrl = false
		}
		if detail != "" {
			json.Unmarshal([]byte(detail), &rWeb)
			switch rWeb.WebName {
			case "tescolotus":
				webscrapping(rWeb)
			}
			// case "makroclick":
			// 	fmt.Println("makroclick")
			// case "bigc":
			// 	fmt.Println("bigc")

			// }
		}
	}
	fmt.Println(time.Now(), "webscrapping bot stop")
}
func webscrapping(web service.Web) {
	db, err := database.Conn()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()
	swRepo := repository.NewSwitchDatabaseDB(db)
	service.NewSwitchDatabaseService(swRepo).ProdudtDetail(web)
}

func main() {
	cronJob()
	// ต้องมีคำสั่งนี้ ไม่งั้นฟังก์ชันการตั้งเวลาทำงานจะไม่ทำงาน เพราะต้องใช้ goroutine
	for {
		time.Sleep(time.Second)
	}
}
