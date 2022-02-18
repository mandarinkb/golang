package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-unit-test/service"
)

type userHandler struct {
	userServ service.UserService
}

func NewUserHandler(userServ service.UserService) userHandler {
	return userHandler{userServ}
}

func (h userHandler) GetUser(c *gin.Context) {
	userServ := h.userServ.GetUser()
	c.IndentedJSON(http.StatusOK, userServ)
}

func (h userHandler) CreateUser(c *gin.Context) {
	var req service.UserReq
	c.BindJSON(&req)
	userServ := h.userServ.CreateUser(req)
	c.IndentedJSON(http.StatusCreated, userServ)
}
func (h userHandler) UpdateUser(c *gin.Context) {
	var req service.UserReq
	c.BindJSON(&req)
	userServ := h.userServ.UpdateUser(req)
	c.IndentedJSON(http.StatusOK, userServ)

}
func (h userHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userServ := h.userServ.DeleteUser(id)
	c.IndentedJSON(http.StatusOK, userServ)
}
