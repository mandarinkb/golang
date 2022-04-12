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

type makroModel struct {
	ProductName       string  `json:"productName"`
	InVatSpecialPrice float64 `json:"inVatSpecialPrice"` //price
	InVatPrice        float64 `json:"inVatPrice"`        //originalPrice
	Image             string  `json:"image"`
	ProductCode       string  `json:"productCode"` //for productUrl
}

func Makroclick(web Web) error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	// เรียก makro api
	data, err := makroApi(web.MenuId, strconv.Itoa(web.MakroPage))
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	//ดึงหน้าเว็บไซต์
	arrProduct := gjson.Get(string(data), "content")
	//แกะหน้าเว็บ
	var makroProducts []makroModel
	err = json.Unmarshal([]byte(arrProduct.String()), &makroProducts)
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}

	for i := range makroProducts {
		// เก็บเฉพาะสินค้าลดราคา
		if makroProducts[i].InVatSpecialPrice != makroProducts[i].InVatPrice {
			discount := (((makroProducts[i].InVatPrice - makroProducts[i].InVatSpecialPrice) / makroProducts[i].InVatPrice) * 100)
			// ปัดเศษ
			discount = math.Round(discount)
			product := Product{
				Timestamp:     time.Now().Format(time.RFC3339),
				WebName:       web.WebName,
				ProductName:   makroProducts[i].ProductName,
				Category:      web.Category,
				Price:         makroProducts[i].InVatSpecialPrice,
				OriginalPrice: makroProducts[i].InVatPrice,
				Discount:      discount,
				ProductUrl:    web.WebUrl + "/products/" + makroProducts[i].ProductCode,
				Image:         makroProducts[i].Image,
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
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "makro : ", makroProducts[i].ProductName)
		}

	}
	return nil
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
