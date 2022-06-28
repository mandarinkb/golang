package service

import "github.com/mandarinkb/go-gorm-session-cookie/repository"

type UserService interface {
	Authenticate(username string, password string) (*repository.User, error)
	Read() ([]repository.User, error)
	ReadById(id int) (*repository.User, error)
	Create(user repository.User) error
	Update(user repository.User) error
	Delete(id int) error
}
