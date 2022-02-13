package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mandarinkb/go-fetch-url-bot-project-final/utils"
	"github.com/tidwall/gjson"
)

type Category struct {
	WebName  string `json:"webName"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

var elasticsearchUrl string

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	elasticsearchUrl = config.Elasticsearch
}

func GetNewCategory(category string) (string, error) {
	// url := "http://127.0.0.1/elasticsearch/web_scrapping_categories/_search"
	url := elasticsearchUrl + "/web_scrapping_categories/_search"
	method := "POST"

	payload := strings.NewReader(`{"query": {"bool": {"must": {"match_phrase": {"tag": "` + category + `"}}}}}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return getCategory(string(body))
}

// category  tescolotus ที่ไม่เอา
func IsNotTakeCategory(webName string, inputCategory string) (bool, error) {
	// url := "http://127.0.0.1/elasticsearch/web_scrapping_ignore_categories/_search"
	url := elasticsearchUrl + "/web_scrapping_ignore_categories/_search"
	method := "POST"

	payload := strings.NewReader(`{"query": {"bool": {"must": {"match_phrase": {"webName": "` + webName + `"}}}}}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	return getDataNotTakeCategory(string(body), inputCategory)
}

func getCategory(body string) (string, error) {
	// อ้างอิงจาก key โดยการใส่ dot ไปเรื่อยๆ
	// กรณีข้อมูลเป็น json array ให้ใส่ # คั่นแล้วตามด้วย key ที่ต้องการ
	source := gjson.Get(body, "hits.hits.#._source")
	arrCate := []Category{}
	err := json.Unmarshal([]byte(source.String()), &arrCate)
	if err != nil {
		return "", err
	}
	var newCate string
	for _, row := range arrCate {
		newCate = row.Category
	}
	return newCate, nil
}

// category  tescolotus ที่ไม่เอา
func getDataNotTakeCategory(body string, inputCategory string) (bool, error) {
	// อ้างอิงจาก key โดยการใส่ dot ไปเรื่อยๆ
	// กรณีข้อมูลเป็น json array ให้ใส่ # คั่นแล้วตามด้วย key ที่ต้องการ
	source := gjson.Get(body, "hits.hits.#._source")
	arrCate := []Category{}
	err := json.Unmarshal([]byte(source.String()), &arrCate)
	if err != nil {
		return false, err
	}
	var categories string
	for _, row := range arrCate {
		categories = row.Category
	}

	cateSlice := strings.Split(categories, ",")
	checkIgnore := false
	for _, cate := range cateSlice {
		if inputCategory == cate {
			checkIgnore = true
		}
	}
	return checkIgnore, nil
}
