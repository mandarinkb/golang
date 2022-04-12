package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gocolly/colly"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/database"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/repository"
	"github.com/mandarinkb/go-fetch-url-bot-with-kafka/utils"
	"github.com/tidwall/gjson"
)

var (
	topicFetchUrl  string = "fetch-url"
	keyFetchUrl    string = "fetch-url"
	topicDetailUrl string = "detail-url"
	keyDetail      string = "detail-url"
)

func MakroMainPage(web Web) error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	categories, err := getCategoryMakro(web.WebUrl)
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	for _, category := range categories {
		detectNotTakeCategory, err := repository.IsNotTakeCategory(web.WebName, category)
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		if !detectNotTakeCategory {
			//หา menuId จาก category
			menuId := makroId(category)
			newCategory, err := repository.GetNewCategory(category)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			web.Category = newCategory
			web.MenuId = menuId

			producer := database.KafkaProducerConn()
			defer producer.Close()
			webStr, err := json.Marshal(web)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			message := &sarama.ProducerMessage{
				Topic: topicFetchUrl,
				Key:   sarama.StringEncoder(keyFetchUrl),
				Value: sarama.StringEncoder(string(webStr)),
			}

			_, offset, err := producer.SendMessage(message)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			fmt.Println("offset : ", offset)
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "makro : ", category, " menuId : ", menuId)
		}

	}
	return nil
}
func MakroFindUrlInPage(web Web) {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	producer := database.KafkaProducerConn()
	defer producer.Close()
	// ดึงหน้าแรก เพื่อหาจำนวนหน้าเว็บไซต์ทั้งหมดก่อน
	data, err := makroApi(web.MenuId, "1")
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	totalPage := totalPage(data)
	// วนดึงข้อมูลทุกหน้าเว็บไซต์
	for i := 1; i < int(totalPage); i++ {
		web.MakroPage = i
		webStr, err := json.Marshal(web)
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		message := &sarama.ProducerMessage{
			Topic: topicDetailUrl,
			Key:   sarama.StringEncoder(keyDetail),
			Value: sarama.StringEncoder(string(webStr)),
		}
		_, offset, err := producer.SendMessage(message)
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		fmt.Println("offset : ", offset)
		fmt.Println(time.Now().Format(time.RFC3339), " : ", "makro : menuId : ", web.MenuId, " page=", i)
	}
}

func makroId(category string) string {
	var id string
	switch category {
	case "ผักและผลไม้":
		id = "3874"
	case "เนื้อสัตว์":
		id = "3896"
	case "ปลาและอาหารทะเล":
		id = "4147"
	case "ไข่ นม เนย ชีส":
		id = "3353"
	case "ผลิตภัณฑ์แปรรูปแช่เย็น":
		id = "82"
	case "ผลิตภัณฑ์เนื้อสัตว์แปรรูป":
		id = "4227"
	case "อาหารแช่แข็ง":
		id = "3932"
	case "เบเกอรีและวัตถุดิบสำหรับทำเบเกอรี":
		id = "3803"
	case "อาหารแห้ง":
		id = "2465"
	case "เครื่องดื่มและขนมขบเคี้ยว":
		id = "2462"
	case "อุปกรณ์และของใช้ในครัวเรือน":
		id = "2460"
	case "ผลิตภัณฑ์ทำความสะอาด":
		id = "4112"
	case "เครื่องเขียนและอุปกรณ์สำนักงาน":
		id = "2464"
	case "เครื่องใช้ไฟฟ้า":
		id = "2461"
	case "สุขภาพและความงาม":
		id = "2466"
	case "สมาร์ทและไลฟ์สไตล์":
		id = "4056"
	case "แม่และเด็ก":
		id = "2467"
	case "ผลิตภัณฑ์สำหรับสัตว์เลี้ยง":
		id = "2468"
	case "Own Brand":
		id = "2634"
	case "สินค้าสั่งพิเศษ และสินค้าเทศกาล":
		id = "3350"
	default:
		fmt.Printf("category not match")
	}
	return id
}

func makroApi(menuId string, page string) ([]byte, error) {
	url := "https://ocs-prod-api.makroclick.com/next-product/public/api/product/search%0A"
	method := "POST"

	payload := strings.NewReader(`{
	  "locale": "th_TH",
	  "minPrice": null,
	  "maxPrice": null,
	  "menuId":` + menuId + `,
	  "hierarchies": [],
	  "customerType": "MKC",
	  "page": ` + page + `,
	  "pageSize": 32,
	  "sorting": "SORTING_LAST_UPDATE",
	  "reloadPrice": true}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getCategoryMakro(url string) ([]string, error) {
	var category []string
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"))
	// ดึงจาก body tag
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// ดึงจาก class MenuCategoryPopOver__MenuListView-sc-77t7qb-2
		e.ForEach(".MenuCategoryPopOver__MenuListView-sc-77t7qb-2", func(_ int, el *colly.HTMLElement) {
			//fmt.Println(el.ChildText("p"))
			category = append(category, el.ChildText("p"))
		})
	})
	// start scraping (ไว้ล่างสุด)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func totalPage(data []byte) int64 {
	return gjson.Get(string(data), "totalPages").Int()
}
