package service

import "github.com/mandarinkb/go-rest-api-jwt-mariadb/middleware"

type UserRequest struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"user_role"`
}
type UserResponse struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	UserRole string `json:"user_role"`
}

type UserService interface {
	Authenticate(username string, password string) (*middleware.TokenResponse, error)
	// VerifyToken(token string) (bool, error)
	Read() ([]UserResponse, error)
	ReadById(id int) (*UserResponse, error)
	Create(user UserRequest) (*UserResponse, error)
	Update(user UserRequest) (*UserResponse, error)
	Delete(id int) error
}
