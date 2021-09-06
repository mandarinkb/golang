package service

type UserRequest struct {
	USER_ID   int    `json:"user_id"`
	USERNAME  string `json:"username"`
	PASSWORD  string `json:"password"`
	USER_ROLE string `json:"user_role"`
}
type UserResponse struct {
	USER_ID   int    `json:"user_id"`
	USERNAME  string `json:"username"`
	USER_ROLE string `json:"user_role"`
}

type UserService interface {
	Read() ([]UserResponse, error)
	ReadById(id int) (*UserResponse, error)
	Create(user UserRequest) error
	Update(user UserRequest) error
	Delete(id int) error
}
