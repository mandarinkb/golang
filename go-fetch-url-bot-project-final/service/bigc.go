package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mandarinkb/go-fetch-url-bot-project-final/database"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/repository"
	"github.com/mandarinkb/go-fetch-url-bot-project-final/utils"
	"github.com/tidwall/gjson"
)

type categoryBigc struct {
	EntityId float32 `json:"entity_id"`
	Name     string  `json:"name"`
}

func BigcMainPage(web Web) error {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	redis := database.RedisConn()
	defer redis.Close()
	categories, err := getCategoryBigc()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}

	for _, row := range categories {
		detectNotTakeCategory, err := repository.IsNotTakeCategory(web.WebName, row.Name)
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		if !detectNotTakeCategory {
			newCategory, err := repository.GetNewCategory(row.Name)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			web.Category = newCategory
			web.EntityId = floatToString(row.EntityId)

			webStr, err := json.Marshal(web)
			if err != nil {
				logger.Error(err.Error(), utils.Url("-"),
					utils.User("-"), utils.Type(utils.TypeBot))
			}
			redis.RPush(context.Background(), "fetchUrl", string(webStr))
			fmt.Println(time.Now().Format(time.RFC3339), " : ", "bigc : ", row.Name, " entityId : ", row.EntityId)
		}
	}
	return nil
}
func BigcFindUrlInPage(web Web) {
	logger, err := utils.LogConf()
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// close log
	defer logger.Sync()

	redis := database.RedisConn()
	defer redis.Close()
	// ดึงหน้าแรกก่อนเพื่อหาหน้าสุดท้าย
	bigcData, err := bigcApi(web.EntityId, "1")
	if err != nil {
		logger.Error(err.Error(), utils.Url("-"),
			utils.User("-"), utils.Type(utils.TypeBot))
	}
	// หา last page ต่อ
	totalPage := lastPage(bigcData)
	for i := 1; i < int(totalPage); i++ {
		web.BigcPage = i
		webStr, err := json.Marshal(web)
		if err != nil {
			logger.Error(err.Error(), utils.Url("-"),
				utils.User("-"), utils.Type(utils.TypeBot))
		}
		redis.RPush(context.Background(), "detailUrl", string(webStr))
		fmt.Println(time.Now().Format(time.RFC3339), " : ", "bigc : entityId : ", web.EntityId, " page=", i)
	}
}

func getCategoryBigc() ([]categoryBigc, error) {
	// call category api from bigc
	url := "https://www.bigc.co.th/api/categories/mainCategory?_store=2"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//ดึงหน้าเว็บไซต์
	arrCategory := gjson.Get(string(body), "result")

	// สร้าง slice categoryBigc
	var categories []categoryBigc
	err = json.Unmarshal([]byte(arrCategory.String()), &categories)
	if err != nil {
		return nil, err
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

func lastPage(data []byte) int64 {
	return gjson.Get(string(data), "result.lastPage").Int()
}

func floatToString(f float32) string {
	return fmt.Sprintf("%v", f)
}
