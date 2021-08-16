package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"webscrapping/util"

	"github.com/tidwall/gjson"
)

var bigcDomain = "https://www.bigc.co.th/"

type bigcModel struct {
	Name              string  `json:"name"`
	FinalPriceInclTax float32 `json:"final_price_incl_tax"` //price
	PriceInclTax      float32 `json:"price_incl_tax"`       //originalPrice
	Image             string  `json:"image"`
	UrlKey            string  `json:"url_key"` //for productUrl
}
type categoryBigc struct {
	EntityId float32 `json:"entity_id"`
	Name     string  `json:"name"`
}

// ต้องมี NewBigC
// เพราะเป็นตัวเชื่อมใช้ implement interface
// โดยจะเชื่อมเฉพาะ ที่เป็น receiver function เท่านั้น
func NewBigC() ProductRepositoy {
	return Product{}
}
func (Product) Bigc() ([]Product, error) {
	categories, err := getCategoryBigc()
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	var productSlice []Product
	for _, row := range categories {
		// ดึงหน้าแรกก่อนเพื่อหาหน้าสุดท้าย
		bigcData, err := bigcApi(util.FloatToString(row.EntityId), "1")
		if err != nil {
			fmt.Println(err)
			//return nil, err
		}
		// หา last page ต่อ
		totalPage := lastPage(bigcData)
		for i := 1; i < int(totalPage); i++ {
			data, err := bigcApi(util.FloatToString(row.EntityId), strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
				//return nil, err
			}
			//ดึงหน้าเว็บไซต์
			arrProduct := gjson.Get(string(data), "result.items")

			//แกะหน้าเว็บ
			var bigcProducts []bigcModel
			err = json.Unmarshal([]byte(arrProduct.String()), &bigcProducts)
			if err != nil {
				fmt.Println(err)
				//return nil, err
			}

			for i := range bigcProducts {
				fmt.Println(time.Now().Format(time.RFC3339), " : ", "bigc : ", bigcProducts[i].Name)
				product := Product{
					Timestamp:     time.Now().Format(time.RFC3339),
					WebName:       "bigc",
					ProductName:   bigcProducts[i].Name,
					Category:      row.Name,
					Price:         bigcProducts[i].FinalPriceInclTax,
					OriginalPrice: bigcProducts[i].PriceInclTax,
					ProductUrl:    bigcDomain + bigcProducts[i].UrlKey,
					Image:         bigcProducts[i].Image,
					Icon:          "https://www.bigc.co.th/_nuxt/img/CI-Bigc-resize.108a02e.png",
				}
				productSlice = append(productSlice, product) //เก็บข้อมูลลง slice Product
				// j, _ := json.Marshal(product)
				// fmt.Println(string(j))
				// return productSlice, nil
			}
		}
	}
	return productSlice, nil
}

func getCategoryBigc() ([]categoryBigc, error) {
	// call category api from bigc
	url := "https://www.bigc.co.th/api/categories/mainCategory?_store=2"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
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
	//ดึงหน้าเว็บไซต์
	arrCategory := gjson.Get(string(body), "result")

	// สร้าง slice categoryBigc
	var categories []categoryBigc
	err = json.Unmarshal([]byte(arrCategory.String()), &categories)
	if err != nil {
		fmt.Println(err)
		//return nil, err
	}
	return categories, nil
}

func bigcApi(cateId string, pageNo string) ([]byte, error) {
	url := "https://www.bigc.co.th/api/categories/getproductListBycateId?_store=2&stock_id=1"
	method := "POST"

	payload := strings.NewReader(`{
	  "cate_id": "` + cateId + `",
	  "ignore_items": "",
	  "page_no": ` + pageNo + `,
	  "page_size": 36,
	  "selected_categories": "",
	  "selected_brands": "",
	  "sort_by": "bestsellers:desc",
	  "price_from": "",
	  "price_to": "",
	  "filter": [],
	  "stock_id": 1
  }`)

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

func lastPage(data []byte) int64 {
	return gjson.Get(string(data), "result.lastPage").Int()
}
