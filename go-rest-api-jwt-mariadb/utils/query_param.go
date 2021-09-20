package utils

import "strconv"

type Param struct {
	ProductName string
	Page        string
	Limit       string
}

type queryParam struct{}

func NewQueryParam() queryParam {
	return queryParam{}
}

func (queryParam) DefaultLimit() string {
	return "10"
}

func (queryParam) SetPage(pageStr string, limitStr string) (int, int, error) {
	//strconv.Atoi() convert string to int
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, 0, err
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, err
	}

	rowCount := (page - 1) * limit
	return rowCount, limit, nil
}
