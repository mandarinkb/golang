package service

import (
	"errors"

	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/repository"
	"github.com/mandarinkb/go-api-project-final/utils"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserServ(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (s userService) Authenticate(username string, password string) (*middleware.TokenResponse, error) {
	var errIncorrect = errors.New("username or password incorrect")
	userRepo, err := s.userRepo.Authenticate(username)
	// กรณีไม่พบ username ใน database
	if err != nil {
		return nil, errIncorrect
	}

	// กรณีพบ ตรวจสอบ password ต่อ และ รหัสผ่านถูกต้อง
	if utils.CheckPasswordHash(password, userRepo.Password) {
		td, err := middleware.GenerateToken(*userRepo)
		if err != nil {
			return nil, err
		}
		resToken := middleware.TokenResponse{
			AccessToken:  td.AccessToken,
			RefreshToken: td.RefreshToken,
		}
		return &resToken, nil
	}
	// กรณีรหัสผ่านไม่ถูกต้อง
	return nil, errIncorrect
}

func (s userService) Read() (users []UserResponse, err error) {
	userRepo, err := s.userRepo.Read()
	if err != nil {
		return nil, err
	}

	for _, row := range userRepo {
		users = append(users, mapDataUserResponse(row))
	}

	return users, nil
}

func (s userService) ReadById(id int) (*UserResponse, error) {
	userRepo, err := s.userRepo.ReadById(id)
	if err != nil {
		return nil, err
	}
	userRes := mapDataUserResponse(*userRepo)
	return &userRes, nil
}

func (s userService) Create(user UserRequest) (*UserResponse, error) {
	userReq := mapDataUserRequest(user)
	userRes := mapDataUserResponse(userReq)
	return &userRes, s.userRepo.Create(userReq)
}

func (s userService) Update(user UserRequest) (*UserResponse, error) {
	// passwordHash, err := utils.NewBcrypt().HashPassword(user.Password)
	// if err != nil {
	// 	return nil, err
	// }
	// // map data from request
	// updateUser := repository.User{
	// 	UserId:   user.UserId,
	// 	Username: user.Username,
	// 	Password: passwordHash,
	// 	UserRole: user.Role,
	// }
	userReq := mapDataUserRequest(user)
	userRes := mapDataUserResponse(userReq)
	return &userRes, s.userRepo.Update(userReq)
}

func (s userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

// แปลงค่า เพื่อส่งไปยัง repository
func mapDataUserRequest(user UserRequest) repository.User {
	passwordHash, _ := utils.HashPassword(user.Password)
	return repository.User{
		UserId:   user.UserId,
		Username: user.Username,
		Password: passwordHash,
		UserRole: user.Role,
	}
}

// แปลงค่า เพื่อส่งไปยัง handler
func mapDataUserResponse(userRepo repository.User) UserResponse {
	return UserResponse{
		Username: userRepo.Username,
		UserRole: userRepo.UserRole}
}
