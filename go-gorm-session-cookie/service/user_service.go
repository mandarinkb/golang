package service

import (
	"errors"

	"github.com/mandarinkb/go-gorm-session-cookie/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userService{userRepo}
}
func (s userService) Authenticate(username string, password string) (*repository.User, error) {
	user, err := s.userRepo.Authenticate(username)
	if err != nil {
		return nil, err
	}
	if password != user.Password {
		return nil, errors.New("invalid username or password")
	}
	return user, nil
}
func (s userService) Read() ([]repository.User, error) {
	users, err := s.userRepo.Read()
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (s userService) ReadById(id int) (*repository.User, error) {
	user, err := s.userRepo.ReadById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s userService) Create(user repository.User) error {
	err := s.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}
func (s userService) Update(user repository.User) error {
	err := s.userRepo.Update(user)
	if err != nil {
		return err
	}
	return nil
}
func (s userService) Delete(id int) error {
	err := s.userRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
