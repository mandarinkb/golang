package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-unit-test/handler"
	"github.com/mandarinkb/go-unit-test/repository"
	"github.com/mandarinkb/go-unit-test/service"
	"github.com/stretchr/testify/assert"
)

var reqBody = service.UserReq{
	UserId:   2,
	Username: "mit",
	Password: "1234",
	UserRole: "admin",
}

func TestNewUserHandler(t *testing.T) {
	newUserRepo := repository.NewUser()
	newUserServ := service.NewUser(newUserRepo)
	newUserHandler := handler.NewUserHandler(newUserServ)
	assert.NotEqual(t, nil, newUserHandler)
}

func TestGetUser(t *testing.T) {
	newUserRepo := repository.NewUser()
	newUserServ := service.NewUser(newUserRepo)
	newUserHandler := handler.NewUserHandler(newUserServ)

	router := gin.Default()
	router.GET("/users", newUserHandler.GetUser)

	reqWriter := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users", nil)
	defer req.Body.Close()

	router.ServeHTTP(reqWriter, req)
	assert.Equal(t, http.StatusOK, reqWriter.Code)

}
func TestCreateUser(t *testing.T) {
	json, _ := json.Marshal(reqBody)

	newUserRepo := repository.NewUser()
	newUserServ := service.NewUser(newUserRepo)
	newUserHandler := handler.NewUserHandler(newUserServ)

	router := gin.Default()
	router.POST("/users", newUserHandler.CreateUser)

	resWriter := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(json))
	defer req.Body.Close()

	router.ServeHTTP(resWriter, req)
	assert.Equal(t, http.StatusCreated, resWriter.Code)

}

func TestUpdateUser(t *testing.T) {
	json, _ := json.Marshal(reqBody)

	newUserRepo := repository.NewUser()
	newUserServ := service.NewUser(newUserRepo)
	newUserHandler := handler.NewUserHandler(newUserServ)

	router := gin.Default()
	router.PUT("/users", newUserHandler.UpdateUser)

	resWriter := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/users", bytes.NewReader(json))
	defer req.Body.Close()

	router.ServeHTTP(resWriter, req)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}
func TestDeleteUser(t *testing.T) {
	newUserRepo := repository.NewUser()
	newUserServ := service.NewUser(newUserRepo)
	newUserHandler := handler.NewUserHandler(newUserServ)

	router := gin.Default()
	router.DELETE("/users/:id", newUserHandler.DeleteUser)

	resWriter := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/users/1", nil)
	defer req.Body.Close()

	router.ServeHTTP(resWriter, req)
	assert.Equal(t, http.StatusOK, resWriter.Code)
}
