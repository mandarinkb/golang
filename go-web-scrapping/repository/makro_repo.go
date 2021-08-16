package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

var makroDomail = "https://www.makroclick.com/th/products/"

type makroModel struct {
	ProductName       string  `json:"productName"`
	InVatSpecialPrice float32 `json:"inVatSpecialPrice"` //price
	InVatPrice        float32 `json:"inVatPrice"`        //originalPrice
	Image             string  `json:"image"`
	ProductCode       string  `json:"productCode"` //for productUrl
}

// ต้องมี NewMakro
// เพราะเป็นตัวเชื่อมใช้ implement interface
// โดยจะเชื่อมเฉพาะ ที่เป็น receiver function เท่านั้น
func NewMakro() ProductRepositoy {
	return Product{}
}

func (Product) Makro() ([]Product, error) {
	url := "https://www.makroclick.com/th"
	categories, err := getCategoryMakro(url)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	var productSlice []Product
	for _, category := range categories {
		//หา menuId จาก category
		menuId := makroId(category)
		// ดึงหน้าแรก เพื่อหาจำนวนหน้าเว็บไซต์ทั้งหมดก่อน
		data, err := makroApi(menuId, "1")
		if err != nil {
			fmt.Println(err)
			//return nil, err
		}
		totalPage := totalPage(data)
		// วนดึงข้อมูลทุกหน้าเว็บไซต์
		for i := 1; i < int(totalPage); i++ {
			// เรียก makro api
			data, err := makroApi(menuId, strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
				//return nil, err
			}
			//ดึงหน้าเว็บไซต์
			arrProduct := gjson.Get(string(data), "content")
			//แกะหน้าเว็บ
			var makroProducts []makroModel
			e := json.Unmarshal([]byte(arrProduct.String()), &makroProducts)
			if e != nil {
				fmt.Println(e)
				//return nil, e
			}

			for i := range makroProducts {
				fmt.Println(time.Now().Format(time.RFC3339), " : ", "makro : ", makroProducts[i].ProductName)
				product := Product{
					Timestamp:     time.Now().Format(time.RFC3339),
					WebName:       "makro",
					ProductName:   makroProducts[i].ProductName,
					Category:      category,
					Price:         makroProducts[i].InVatSpecialPrice,
					OriginalPrice: makroProducts[i].InVatPrice,
					ProductUrl:    makroDomail + makroProducts[i].ProductCode,
					Image:         makroProducts[i].Image,
					Icon:          "https://www.makroclick.com/static/images/logo.png",
				}
				productSlice = append(productSlice, product) //เก็บข้อมูลลง slice Product
			}
		}
	}
	return productSlice, nil
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
		fmt.Println(err)
		//return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		//return nil, err
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
		fmt.Println(err)
		//return nil, err
	}
	return category, nil
}

func totalPage(data []byte) int64 {
	return gjson.Get(string(data), "totalPages").Int()
}
