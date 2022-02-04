package service

import "github.com/mandarinkb/go-api-project-final/middleware"

type UserRequest struct {
	UserId             int    `json:"userId"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	Role               string `json:"role"`
	NewPassword        string `json:"newPassword"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}
type UserResponse struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	UserRole string `json:"role"`
}

type UserService interface {
	Authenticate(username string, password string) (*middleware.TokenResponse, error)
	Read() ([]UserResponse, error)
	ReadById(id int) (*UserResponse, error)
	Create(user UserRequest) (*UserResponse, error)
	Update(user UserRequest) (*UserResponse, error)
	Delete(id int) (string, error)
}
