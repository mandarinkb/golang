package repository

type User struct {
	UserId   int    `db:"USER_ID" json:"userId"`
	Username string `db:"USERNAME" json:"username"`
	Password string `db:"PASSWORD" json:"password"`
	UserRole string `db:"USER_ROLE" json:"role"`
}
type UserRepository interface {
	Authenticate(username string) (*User, error)
	Read() ([]User, error)
	ReadById(id int) (*User, error)
	Create(user User) error
	Update(user User) error
	Delete(id int) error
}
