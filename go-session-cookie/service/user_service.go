package service

import (
	"errors"

	"github.com/mandarinkb/go-session-cookie/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserServ(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (s userService) Authenticate(username string, password string) error {
	userAuth, err := s.userRepo.Authenticate(username)
	// กรณีไม่พบ username ใน database
	if err != nil || password != userAuth.Password {
		return errors.New("invalid username or password.")
	}
	return nil
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
	//ตรวจสอบ user ที่สร้างว่ามีในระบบหร่ือไม่
	_, err := s.userRepo.Authenticate(user.Username)
	// กรณีไม่พบ username ใน database
	if err != nil {
		userReq := mapDataUserRequest(user)
		userRes := mapDataUserResponse(userReq)
		return &userRes, s.userRepo.Create(userReq)
	}
	return nil, errors.New("username already exists")
}

func (s userService) Update(user UserRequest) (*UserResponse, error) {
	userReq := mapDataUserRequest(user)
	userRes := mapDataUserResponse(userReq)
	return &userRes, s.userRepo.Update(userReq)
}

func (s userService) Delete(id int) (*UserResponse, error) {
	userRepo, err := s.userRepo.ReadById(id)
	if err != nil {
		return nil, err
	}
	userRes := mapDataUserResponse(*userRepo)
	return &userRes, s.userRepo.Delete(id)
}

// แปลงค่า เพื่อส่งไปยัง repository
func mapDataUserRequest(user UserRequest) repository.User {
	return repository.User{
		UserId:   user.UserId,
		Username: user.Username,
		Password: user.Password,
		UserRole: user.Role,
	}
}

// แปลงค่า เพื่อส่งไปยัง handler
func mapDataUserResponse(userRepo repository.User) UserResponse {
	return UserResponse{
		UserId:   userRepo.UserId,
		Username: userRepo.Username,
		UserRole: userRepo.UserRole}
}
