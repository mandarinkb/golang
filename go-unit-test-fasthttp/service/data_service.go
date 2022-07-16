package service

import (
	"fmt"
	"strconv"

	"github.com/mandarinkb/test-git/repository"
)

type dataService struct {
	dataRepo repository.DataReposity
}

func NewDataService(dataRepo repository.DataReposity) DataService {
	return dataService{dataRepo}
}

func (s dataService) GetAll() ([]repository.Data, error) {
	dataMock, err := s.dataRepo.GetAll()
	if err != nil {
		return nil, err
	}
	// (business logic) all CreditAmount - 5
	newData := []repository.Data{}
	for _, data := range dataMock {
		stringData := data.CreditAmount
		floatData, err := strconv.ParseFloat(stringData, 32)
		if err != nil {
			return nil, err
		}
		floatData = floatData - 5
		data.CreditAmount = fmt.Sprintf("%.4f", floatData)
		newData = append(newData, data)
	}

	return newData, nil
}

func (s dataService) GetById(id int) (*repository.Data, error) {
	dataMock, err := s.dataRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	// (business logic) all CreditAmount - 5
	stringData := dataMock.CreditAmount
	floatData, err := strconv.ParseFloat(stringData, 32)
	if err != nil {
		return nil, err
	}
	floatData = floatData - 5
	dataMock.CreditAmount = fmt.Sprintf("%.4f", floatData)

	return dataMock, nil
}
