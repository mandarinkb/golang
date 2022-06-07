package service

type UserRequest struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type UserResponse struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	UserRole string `json:"role"`
}

type UserService interface {
	Authenticate(username string, password string) error
	Read() ([]UserResponse, error)
	ReadById(id int) (*UserResponse, error)
	Create(user UserRequest) (*UserResponse, error)
	Update(user UserRequest) (*UserResponse, error)
	Delete(id int) (*UserResponse, error)
}
