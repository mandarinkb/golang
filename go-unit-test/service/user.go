package service

import "github.com/mandarinkb/go-unit-test/repository"

type User struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	UserRole string `json:"userRole"`
}
type UserReq struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"userRole"`
}
type userService struct {
	userRepo repository.UserRepository
}
type UserService interface {
	GetUser() []User
	CreateUser(user UserReq) []User
	UpdateUser(user UserReq) []User
	DeleteUser(id int) []User
}

func NewUser(user repository.UserRepository) UserService {
	return userService{user}
}

func (s userService) GetUser() (users []User) {
	userRepo := s.userRepo.GetUser()
	for _, row := range userRepo {
		users = append(users, mapRepoToServ(row))
	}
	return users
}
func (s userService) CreateUser(user UserReq) (users []User) {
	userRepo := s.userRepo.CreateUser(mapReqToRepo(user))
	for _, row := range userRepo {
		users = append(users, mapRepoToServ(row))
	}
	return users
}
func (s userService) UpdateUser(user UserReq) (users []User) {
	userRepo := s.userRepo.UpdateUser(mapReqToRepo(user))
	for _, row := range userRepo {
		users = append(users, mapRepoToServ(row))
	}
	return users
}
func (s userService) DeleteUser(id int) (users []User) {
	userRepo := s.userRepo.DeleteUser(id)
	for _, row := range userRepo {
		users = append(users, mapRepoToServ(row))
	}
	return users
}
func mapRepoToServ(userRepo repository.User) User {
	return User{
		UserId:   userRepo.UserId,
		Username: userRepo.Username,
		UserRole: userRepo.UserRole,
	}
}
func mapReqToRepo(userReq UserReq) repository.User {
	return repository.User{
		UserId:   userReq.UserId,
		Username: userReq.Username,
		Password: userReq.Password,
		UserRole: userReq.UserRole,
	}
}
