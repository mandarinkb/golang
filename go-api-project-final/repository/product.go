package repository

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/mandarinkb/go-api-project-final/utils"
	"github.com/tidwall/gjson"
)

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
type ProductRepository interface {
	GetSearch(index string, from string, jsonName string) ([]Product, error)
	GetFilterSearch(index string, from string, jsonName string,
		jsonWebName string, min string, max string) ([]Product, error)
	GetCategory(index string, from string, category string, jsonWebName string) ([]Product, error)
	GetHistory(index string, from string) ([]Product, error)
	GetHistoryWithHistoryData(index string, from string, jsonName string) ([]Product, error)
}

var elasUrl string

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	elasUrl = config.Elasticsearch
}

func NewProductRepo() ProductRepository {
	return Product{}
}

func (Product) GetSearch(index string, from string, jsonName string) ([]Product, error) {
	url := elasUrl + "/" + index + "/_search"
	method := "POST"

	payload := strings.NewReader(`{"from": ` + from + `,"size": 50,
  	"sort": {"discount": "desc"},
  	"query": {
		"bool": {
		   "must": [` + jsonName + `]
		}
	  }
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
	return getSourceProduct(string(body))
}
func (p Product) GetFilterSearch(index string, from string, jsonName string,
	jsonWebName string, min string, max string) ([]Product, error) {
	url := elasUrl + "/" + index + "/_search"
	method := "POST"

	payload := strings.NewReader(`{"from": ` + from + `,"size": 50,
	  "sort": {"discount": "desc"},
	  "query": {
	  		"bool": {
	  			"must": [` + jsonName + `
	  				,{"dis_max": {"queries": [` + jsonWebName + `]}},
	  					{"range": {"price": {"gte": ` + min + `,"lte": ` + max + `}}}
					]
	  			}
	 		}
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
	// fmt.Println(string(body))
	return getSourceProduct(string(body))
}
func (p Product) GetCategory(index string, from string, category string, jsonWebName string) ([]Product, error) {
	url := elasUrl + "/" + index + "/_search"
	method := "POST"

	payload := strings.NewReader(`{"from": ` + from + `,"size": 50,
  	"sort": {"discount": "desc"},
  	"query": {
	  "bool": {
		  "must": [{ "match": {"category": "` + category + `"}},
			  {"dis_max": 
				  { "queries": [` + jsonWebName + `]
				  }
			  }]
		  }
	  }
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

	return getSourceProduct(string(body))
}
func (p Product) GetHistory(index string, from string) ([]Product, error) {
	url := elasUrl + "/" + index + "/_search"
	method := "POST"

	payload := strings.NewReader(`{"from": ` + from + `,"size": 50,"sort": {"discount": "desc"},"query": {"match_all": {}}}`)

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

	return getSourceProduct(string(body))
}
func (p Product) GetHistoryWithHistoryData(index string, from string, jsonName string) (product []Product, err error) {
	url := elasUrl + "/" + index + "/_search"
	method := "POST"

	payload := strings.NewReader(`{"from": ` + from + `,"size": 50,
  	"sort": {"discount": "desc"},
  	"query": {
		"bool": {
		   "should": [` + jsonName + `]
		}
	  }
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
	return getSourceProduct(string(body))
}

func getSourceProduct(body string) ([]Product, error) {
	// อ้างอิงจาก key โดยการใส่ dot ไปเรื่อยๆ
	// กรณีข้อมูลเป็น json array ให้ใส่ # คั่นแล้วตามด้วย key ที่ต้องการ
	source := gjson.Get(body, "hits.hits.#._source")
	arrSource := []Product{}
	err := json.Unmarshal([]byte(source.String()), &arrSource)
	if err != nil {
		return nil, err
	}
	return arrSource, nil
}
