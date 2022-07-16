package service

import "github.com/mandarinkb/test-git/repository"

type DataService interface {
	GetAll() ([]repository.Data, error)
	GetById(id int) (*repository.Data, error)
}
