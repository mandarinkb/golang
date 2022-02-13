package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mandarinkb/go-fetch-url-bot-project-final/database"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/service"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/utils"
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

	logger.Info("[fetch url bot] start", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))

	fmt.Println(time.Now(), " : fetch url bot start")
	redis := database.RedisConn()
	defer redis.Close()

	webData := service.Web{}
	// จัดการหมวดหมู่สินค้าใหม่
	checkStartUrl := true
	for checkStartUrl {
		startUrl, err := redis.RPop(context.Background(), "startUrl").Result()
		if err != nil {
			checkStartUrl = false
		}
		if startUrl != "" {
			err = json.Unmarshal([]byte(startUrl), &webData)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
				checkStartUrl = false
			}
			if webData.WebStatus == "1" {
				switch webData.WebName {
				case "tescolotus":
					service.TescolotusMainPage(webData)
				case "makroclick":
					service.MakroMainPage(webData)
				case "bigc":
					service.BigcMainPage(webData)
				}
			}
		}
	}
	// หา url ของสินค้าในแต่ละหมวดหมู่ โดยจะหาทุกๆหน้า
	checkFetchUrl := true
	for checkFetchUrl {
		fetchUrl, err := redis.RPop(context.Background(), "fetchUrl").Result()
		if err != nil {
			checkFetchUrl = false
		}
		if fetchUrl != "" {
			json.Unmarshal([]byte(fetchUrl), &webData)
			if webData.WebStatus == "1" {
				switch webData.WebName {
				case "tescolotus":
					service.TescolotusFindUrlInPage(webData)
				case "makroclick":
					service.MakroFindUrlInPage(webData)
				case "bigc":
					service.BigcFindUrlInPage(webData)
				}
			}
		}
	}

	logger.Info("[fetch url bot] stop", utils.Url("-"),
		utils.User("-"), utils.Type(utils.TypeBot))
	fmt.Println(time.Now(), " : fetch url bot stop")
}

func main() {
	cronJob()
	// ต้องมีคำสั่งนี้ ไม่งั้นฟังก์ชันการตั้งเวลาทำงานจะไม่ทำงาน เพราะต้องใช้ goroutine
	for {
		time.Sleep(time.Second)
	}

}

// #######tesco lotus########
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

// ########bigc########
//   โปรโมชั่นโบรชัวร์
//   พร้อมรับมือ โควิด-19
//   สินค้าส่งด่วน 1 ชม.
//   ร้านค้าส่ง
//   สินค้าเฉพาะบิ๊กซี
//   อาหารสด, แช่แข็ง/ ผักผลไม้
//   อาหารแห้ง/ เครื่องปรุง
//   เครื่องดื่ม/ ขนมขบเคี้ยว
//   สุขภาพและความงาม
//   แม่และเด็ก
//   ของใช้ในครัวเรือน/ สัตว์เลี้ยง
//   เครื่องใช้ไฟฟ้า/ อุปกรณ์อิเล็กทรอนิกส์
//   บ้านและไลฟ์สไตล์
//   เครื่องเขียน/ อุปกรณ์สำนักงาน
//   เสื้อผ้า/ เครื่องประดับ
//   ร้านเพรียวฟาร์มาซี
//   สินค้าแบรนด์เบสิโค

// #######makro#########
// ผักและผลไม้
// เนื้อสัตว์
// ปลาและอาหารทะเล
// ไข่ นม เนย ชีส
// ผลิตภัณฑ์แปรรูปแช่เย็น
// ผลิตภัณฑ์เนื้อสัตว์แปรรูป
// อาหารแช่แข็ง
// เบเกอรีและวัตถุดิบสำหรับทำเบเกอรี
// อาหารแห้ง
// เครื่องดื่มและขนมขบเคี้ยว
// อุปกรณ์และของใช้ในครัวเรือน
// ผลิตภัณฑ์ทำความสะอาด
// เครื่องเขียนและอุปกรณ์สำนักงาน
// เครื่องใช้ไฟฟ้า
// สุขภาพและความงาม
// สมาร์ทและไลฟ์สไตล์
// แม่และเด็ก
// ผลิตภัณฑ์สำหรับสัตว์เลี้ยง
// Own Brand
// สินค้าสั่งพิเศษ และสินค้าเทศกาล
