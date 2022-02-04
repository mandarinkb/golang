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

type LogSystem struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Caller    string `json:"caller"`
	User      string `json:"user"`
	Massage   string `json:"msg"`
	Url       string `json:"url"`
	TypeLog   string `json:"typeLog"`
}

var elasticsearchUrl string

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	elasticsearchUrl = config.Elasticsearch
}

type LogRepository interface {
	GetLogs(date string) ([]LogSystem, error)
}

func NewLogRepo() LogRepository {
	return LogSystem{}
}

func (LogSystem) GetLogs(date string) ([]LogSystem, error) {
	url := elasticsearchUrl + "/web_scrapping_log/_search"
	method := "POST"

	payload := strings.NewReader(`{
	  "from": 0,
	  "size": 1000,
	  "sort": {
		  "timestamp": "desc"
	  },
	  "query": {
		  "bool": {
			  "must": {
				  "match_phrase": {
					  "timestamp": "` + date + `"
				  }
			  }
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
	return getSource(string(body))
}

func getSource(body string) ([]LogSystem, error) {
	// อ้างอิงจาก key โดยการใส่ dot ไปเรื่อยๆ
	// กรณีข้อมูลเป็น json array ให้ใส่ # คั่นแล้วตามด้วย key ที่ต้องการ
	source := gjson.Get(body, "hits.hits.#._source")
	// fmt.Println(source.String())
	arrSource := []LogSystem{}
	err := json.Unmarshal([]byte(source.String()), &arrSource)
	if err != nil {
		return nil, err
	}
	return arrSource, nil
}
