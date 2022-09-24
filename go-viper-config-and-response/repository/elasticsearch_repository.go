package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mandarinkb/go-viper-config-and-response/config"
	"github.com/mandarinkb/go-viper-config-and-response/external/elastic"
)

type elasticSearchRepository struct {
	HTTPESClient *http.Client
}

func NewElasticSearchRepository() ElasticSearchRepository {
	return &elasticSearchRepository{
		HTTPESClient: elastic.NewHttpESClient(),
	}
}

func (es *elasticSearchRepository) ESBuildJsonToByte(data interface{}) ([]byte, error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return dataByte, err
}

func (es *elasticSearchRepository) ESPost(index, indexType string, msg []byte) ([]byte, error) {
	// prepare data
	url := fmt.Sprintf("%s/%s/%s", config.C().Elastic.ServerWithPort, index, indexType)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return nil, err
	}

	// call api
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := es.HTTPESClient.Do(req)
	if err != nil {
		return nil, err
	}

	// get response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}
