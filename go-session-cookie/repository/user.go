package repository

type User struct {
	UserId   int    `db:"USER_ID"`
	Username string `db:"USERNAME"`
	Password string `db:"PASSWORD"`
	UserRole string `db:"USER_ROLE"`
}
type UserRepository interface {
	Authenticate(username string) (*User, error)
	Read() ([]User, error)
	ReadById(id int) (*User, error)
	Create(user User) error
	Update(user User) error
	Delete(id int) error
}
