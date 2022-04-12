package service

import (
	"log"
	"net/http"
	"strings"

	"github.com/mandarinkb/go-web-scraping-bot-with-kafka/utils"
)

type Web struct {
	WebId     int    `json:"webId"`
	WebName   string `json:"webName"`
	WebUrl    string `json:"webUrl"`
	WebStatus string `json:"webStatus"`
	IconUrl   string `json:"iconUrl"`
	Category  string `json:"category"`
	MenuId    string `json:"menuId"`    // for makro
	MakroPage int    `json:"makroPage"` // for makro
	EntityId  string `json:"entityId"`  // for bigc
	BigcPage  int    `json:"bigcPage"`  // for bigc
}
type Product struct {
	Timestamp     string  `json:"timestamp"`
	WebName       string  `json:"webName"`
	ProductName   string  `json:"productName"`
	Category      string  `json:"category"`
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalPrice"`
	Discount      float64 `json:"discount"`
	ProductUrl    string  `json:"productUrl"`
	Image         string  `json:"image"`
	Icon          string  `json:"icon"`
}
type ProductService interface {
	Tescolotus(web Web) error
}

var elasticsearchUrl string

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	elasticsearchUrl = config.Elasticsearch
}
func insertToElasticsearch(dbName string, product string) error {

	url := elasticsearchUrl + "/" + dbName + "/txt"
	method := "POST"

	payload := strings.NewReader(product)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
