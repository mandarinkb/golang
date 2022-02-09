package repository

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mandarinkb/go-start-bot-project-final/utils"
)

var elasticsearchUrl string

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	elasticsearchUrl = config.Elasticsearch
}
func DeleteIndex(index string) error {
	url := elasticsearchUrl + "/" + index
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	_ = body
	if err != nil {
		return err
	}
	return nil
}
