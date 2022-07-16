package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type mockData struct{}

func NewMock() DataReposity {
	return mockData{}
}
func (mockData) GetAll() ([]Data, error) {
	var sliceMock []Data
	ndjsonFile, err := os.Open("./repository/data.ndjson")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer ndjsonFile.Close()
	data := json.NewDecoder(ndjsonFile)
	for {
		var mock Data
		// Decode one JSON document.
		err := data.Decode(&mock)
		if err != nil {
			//end of stream.
			break
		}
		// Do something with the value.
		sliceMock = append(sliceMock, mock)
	}
	return sliceMock, nil
}

func (mockData) GetById(id int) (*Data, error) {
	ndjsonFile, err := os.Open("./repository/data.ndjson")
	if err != nil {
		return nil, err
	}
	defer ndjsonFile.Close()
	data := json.NewDecoder(ndjsonFile)
	for {
		var mock Data
		// Decode one JSON document.
		err := data.Decode(&mock)
		if err != nil {
			//end of stream.
			break
		}
		// query RequestCardTopupCompensateLogID == id
		if mock.RequestCardTopupCompensateLogID == id {
			return &mock, nil
		}
	}
	return nil, errors.New("not found")
}
