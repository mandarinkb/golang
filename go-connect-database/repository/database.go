package repository

type User struct {
	USER_ID  int
	USERNAME string
	PASSWORD string
	ROLE     string
}
type database interface {
	MariadbRead() ([]User, error)
	MariadbReadById(id int) (*User, error)
	MariadbCreate(user User) error
	MariadbUpdate(user User) error
	MariadbDelete(id int) error
}
