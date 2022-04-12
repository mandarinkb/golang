package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/database"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/repository"
	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/utils"
	"github.com/tidwall/gjson"
)

type bigcModel struct {
	Name              string  `json:"name"`
	FinalPriceInclTax float64 `json:"final_price_incl_tax"` //price
	PriceInclTax      float64 `json:"price_incl_tax"`       //originalPrice
	Image             string  `json:"image"`
	UrlKey            string  `json:"url_key"` //for productUrl
}

func Bigc(web Web) error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	data, err := bigcApi(web.EntityId, strconv.Itoa(web.BigcPage))
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	//ดึงหน้าเว็บไซต์
	arrProduct := gjson.Get(string(data), "result.items")

	//แกะหน้าเว็บ
	var bigcProducts []bigcModel
	err = json.Unmarshal([]byte(arrProduct.String()), &bigcProducts)
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}

	for i := range bigcProducts {
		if bigcProducts[i].PriceInclTax != bigcProducts[i].FinalPriceInclTax {
			discount := ((bigcProducts[i].PriceInclTax - bigcProducts[i].FinalPriceInclTax) / bigcProducts[i].PriceInclTax * 100)
			// ปัดเศษ
			discount = math.Round(discount)
			product := Product{
				Timestamp:     time.Now().Format(time.RFC3339),
				WebName:       web.WebName,
				ProductName:   bigcProducts[i].Name,
				Category:      web.Category,
				Price:         bigcProducts[i].FinalPriceInclTax,
				OriginalPrice: bigcProducts[i].PriceInclTax,
				Discount:      discount,
				ProductUrl:    web.WebUrl + bigcProducts[i].UrlKey,
				Image:         bigcProducts[i].Image,
				Icon:          web.IconUrl,
			}
			jsonProducts, err := json.Marshal(product)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}

			db, err := database.Conn()
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			defer db.Close()

			dbRepo, err := repository.NewSwitchDatabaseDB(db).GetInActivateDatabaseName()
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}

			dbname := dbRepo.DatabaseName
			err = insertToElasticsearch(dbname, string(jsonProducts))
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "bigc : ", bigcProducts[i].Name)
		}
	}
	return nil
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
