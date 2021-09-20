package service

import (
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/utils"
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
			UserId:   row.UserId,
			Username: row.Username,
			UserRole: row.UserRole,
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
	// user := UserResponse{
	// 	UserId:   userRepo.UserId,
	// 	Username: userRepo.Username,
	// 	UserRole: userRepo.UserRole}
	userRes := mapDataUser(*userRepo)
	return &userRes, nil
}

func (s userService) Create(user UserRequest) (*UserResponse, error) {
	passwordHash, err := utils.NewBcrypt().HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	// map data from request
	addUser := repository.User{
		UserId:   user.UserId,
		Username: user.Username,
		Password: passwordHash,
		UserRole: user.UserRole,
	}
	userRes := mapDataUser(addUser)
	return &userRes, s.userRepo.Create(addUser)
}

func (s userService) Update(user UserRequest) (*UserResponse, error) {
	passwordHash, err := utils.NewBcrypt().HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	// map data from request
	updateUser := repository.User{
		UserId:   user.UserId,
		Username: user.Username,
		Password: passwordHash,
		UserRole: user.UserRole,
	}
	userRes := mapDataUser(updateUser)
	return &userRes, s.userRepo.Update(updateUser)
}

func (s userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

// ฟังก็ชัน map data จาก repository ไปยัง service
func mapDataUser(userRepo repository.User) UserResponse {
	return UserResponse{
		UserId:   userRepo.UserId,
		Username: userRepo.Username,
		UserRole: userRepo.UserRole}
}
