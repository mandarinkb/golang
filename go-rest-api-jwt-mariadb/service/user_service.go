package service

import (
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserServ(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (s userService) Read() ([]UserResponse, error) {
	users := []UserResponse{}
	userRepo, err := s.userRepo.Read()
	if err != nil {
		return nil, err
	}

	for _, row := range userRepo {
		dataRepo := UserResponse{
			USER_ID:   row.USER_ID,
			USERNAME:  row.USERNAME,
			USER_ROLE: row.USER_ROLE,
		}
		users = append(users, dataRepo)
	}

	return users, nil
}

func (s userService) ReadById(id int) (*UserResponse, error) {
	userRepo, err := s.userRepo.ReadById(id)
	if err != nil {
		return nil, err
	}
	user := UserResponse{
		USER_ID:   userRepo.USER_ID,
		USERNAME:  userRepo.USERNAME,
		USER_ROLE: userRepo.USER_ROLE}

	return &user, nil
}

func (s userService) Create(user UserRequest) error {
	addUser := repository.User{
		USER_ID:   user.USER_ID,
		USERNAME:  user.USERNAME,
		PASSWORD:  user.PASSWORD,
		USER_ROLE: user.USER_ROLE,
	}
	return s.userRepo.Create(addUser)
}

func (s userService) Update(user UserRequest) error {
	updateUser := repository.User{
		USER_ID:   user.USER_ID,
		USERNAME:  user.USERNAME,
		PASSWORD:  user.PASSWORD,
		USER_ROLE: user.USER_ROLE,
	}
	return s.userRepo.Update(updateUser)
}

func (s userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}
