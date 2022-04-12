package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/controller"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/database"

	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/utils"
)

func main() {
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

	consumer := database.KafkaConsumerConn()
	defer consumer.Close()

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient("bot-group", consumer)
	if err != nil {
		panic(err)
	}
	defer group.Close()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{"start-url", "fetch-url"}
		handler := controller.ConsumerGroupHandler{}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
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
