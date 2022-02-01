package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mandarinkb/go-fetch-url-bot-project-final/database"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/service"
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
	fmt.Println(time.Now(), "fetch url bot start")
	redis := database.RedisConn()
	defer redis.Close()

	data := service.Web{}
	// จัดการหมวดหมู่สินค้าใหม่
	checkStartUrl := true
	for checkStartUrl {
		startUrl, err := redis.RPop(context.Background(), "startUrl").Result()
		if err != nil {
			fmt.Println("checkStartUrl : ", err)
			checkStartUrl = false
		}
		if startUrl != "" {
			json.Unmarshal([]byte(startUrl), &data)
			switch data.WebName {
			case "tescolotus":
				service.MainPage(data)
			case "makroclick":
				fmt.Println("makroclick")
			case "bigc":
				fmt.Println("bigc")

			}
		}

	}
	// หา url ของสินค้าในแต่ละหมวดหมู่ โดยจะหาทุกๆหน้า
	checkFetchUrl := true
	for checkFetchUrl {
		fetchUrl, err := redis.RPop(context.Background(), "fetchUrl").Result()
		if err != nil {
			fmt.Println(err)
			checkFetchUrl = false
		}
		if fetchUrl != "" {
			json.Unmarshal([]byte(fetchUrl), &data)
			switch data.WebName {
			case "tescolotus":
				service.FindUrlInPage(data)
			case "makroclick":
				fmt.Println("makroclick")
			case "bigc":
				fmt.Println("bigc")

			}
		}

	}
	fmt.Println(time.Now(), "fetch url bot stop")
}

func main() {
	cronJob()
	// ต้องมีคำสั่งนี้ ไม่งั้นฟังก์ชันการตั้งเวลาทำงานจะไม่ทำงาน เพราะต้องใช้ goroutine
	for {
		time.Sleep(time.Second)
	}
}

// ผลิตภัณฑ์เพื่อสุขภาพ & ความงาม
// เครื่องดื่ม, ขนมขบเคี้ยว & ของหวาน
// ผลิตภัณฑ์ทำความสะอาดในครัวเรือน
// อาหารสด, แช่แข็ง & เบเกอรี่
// อาหารแห้ง & อาหารกระป๋อง
// เครื่องใช้ไฟฟ้า & อุปกรณ์ภายในบ้าน
// สินค้าอื่นๆ
// ผลิตภัณฑ์สำหรับเด็ก
// ผลิตภัณฑ์สำหรับสัตว์เลี้ยง
// เทศกาลต้อนรับเปิดเทอม
// เสื้อผ้าเครื่องแต่งกาย
// ดูทั้งหมด

//
//
//
//
//
//
// สินค้าอื่นๆ
//
//
// เทศกาลต้อนรับเปิดเทอม
//
// ดูทั้งหมด
