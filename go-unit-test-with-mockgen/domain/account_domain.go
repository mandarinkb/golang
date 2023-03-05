package domain

import "github.com/mandarinkb/go-unit-test-with-mockgen/model"

type AccountRepository interface {
	GetInfo(id int) (*model.Account, error)
}
