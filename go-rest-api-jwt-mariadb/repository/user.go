package repository

type User struct {
	USER_ID   int    `db:"USER_ID"`
	USERNAME  string `db:"USERNAME"`
	PASSWORD  string `db:"PASSWORD"`
	USER_ROLE string `db:"USER_ROLE"`
}
type UserRepository interface {
	Read() ([]User, error)
	ReadById(id int) (*User, error)
	Create(user User) error
	Update(user User) error
	Delete(id int) error
}
