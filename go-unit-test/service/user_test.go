package service_test

import (
	"testing"

	"github.com/mandarinkb/go-unit-test/repository"
	"github.com/mandarinkb/go-unit-test/service"
	"github.com/stretchr/testify/assert"
)

var userReq = service.UserReq{
	UserId:   1,
	Username: "aaa",
	Password: "1234",
	UserRole: "admin",
}

func TestNewUserServ(t *testing.T) {
	// t.Run("", func(t *testing.T) {})
	newUserRepo := repository.NewUser()
	newUserServ := service.NewUser(newUserRepo)
	// if newUserServ == nil {
	// 	t.Errorf("got %v", newUserServ)
	// }
	assert.NotEqual(t, nil, newUserServ)
}
func TestGetUser(t *testing.T) {
	// t.Run("", func(t *testing.T) {})
	newUserRepo := repository.NewUser()
	gUserServ := service.NewUser(newUserRepo).GetUser()
	assert.NotEqual(t, nil, gUserServ)
}
func TestCreateUser(t *testing.T) {
	// t.Run("", func(t *testing.T) {})
	newUserRepo := repository.NewUser()
	cUserServ := service.NewUser(newUserRepo).CreateUser(userReq)
	assert.NotEqual(t, nil, cUserServ)
}
func TestUpdateUser(t *testing.T) {
	// t.Run("", func(t *testing.T) {})
	newUserRepo := repository.NewUser()
	uUserServ := service.NewUser(newUserRepo).UpdateUser(userReq)
	assert.NotEqual(t, nil, uUserServ)
}

func TestDeleteUser(t *testing.T) {
	// t.Run("", func(t *testing.T) {})
	newUserRepo := repository.NewUser()
	dUserServ := service.NewUser(newUserRepo).DeleteUser(1)
	assert.NotEqual(t, nil, dUserServ)
}
