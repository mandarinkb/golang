package repository

type User struct {
	UserId   int    `db:"USER_ID"`
	Username string `db:"USERNAME"`
	Password string `db:"PASSWORD"`
	UserRole string `db:"USER_ROLE"`
}
type UserRepository interface {
	GetUser() []User
	CreateUser(user User) []User
	UpdateUser(user User) []User
	DeleteUser(id int) []User
}

func NewUser() UserRepository {
	return User{
		UserId:   1,
		Username: "mandarinkb",
		Password: "1234",
		UserRole: "admin",
	}
}

func (u User) GetUser() (users []User) {
	users = append(users, u)
	return users
}
func (u User) CreateUser(user User) (users []User) {
	users = append(users, u)
	users = append(users, user)
	return users
}
func (u User) UpdateUser(user User) (users []User) {
	uUser := User{
		UserId:   user.UserId,
		Username: user.Username,
		Password: user.Password,
		UserRole: user.UserRole,
	}
	users = append(users, uUser)
	return users
}
func (u User) DeleteUser(id int) (users []User) {
	users = append(users, u)
	return users
}
