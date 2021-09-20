package service

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
	Read() ([]UserResponse, error)
	ReadById(id int) (*UserResponse, error)
	Create(user UserRequest) (*UserResponse, error)
	Update(user UserRequest) (*UserResponse, error)
	Delete(id int) error
}
