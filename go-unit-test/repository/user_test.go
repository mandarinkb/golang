package repository_test

import (
	"testing"

	"github.com/mandarinkb/go-unit-test/repository"
	"github.com/stretchr/testify/assert"
)

var user = repository.User{
	UserId:   1,
	Username: "mandarinkb",
	Password: "1234",
	UserRole: "admin",
}

// ปกติแล้วส่วนของ repository จะไม่ต้องเขียน unit test
func TestNewUser(t *testing.T) {
	// t.Run("", func(t *testing.T) {})
	newUser := repository.NewUser()
	// if newUser == nil {
	// 	t.Errorf("got %v", newUser)
	// }
	assert.NotEqual(t, nil, newUser)
}
func TestGetUser(t *testing.T) {
	getUser := repository.NewUser().GetUser()
	assert.NotEqual(t, nil, getUser)
}
func TestCreateUser(t *testing.T) {
	createUser := repository.NewUser().CreateUser(user)
	assert.NotEqual(t, nil, createUser)
}
func TestUpdateUser(t *testing.T) {
	updateUser := repository.NewUser().UpdateUser(user)
	assert.NotEqual(t, nil, updateUser)
}
func TestDeleteUser(t *testing.T) {
	deleteUser := repository.NewUser().DeleteUser(1)
	assert.NotEqual(t, nil, deleteUser)
}
