package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-gorm-session-cookie/middleware"
	"github.com/mandarinkb/go-gorm-session-cookie/repository"
	"github.com/mandarinkb/go-gorm-session-cookie/service"
)

type userHandler struct {
	userServ service.UserService
}

func NewUserHandler(userServ service.UserService) userHandler {
	return userHandler{userServ}
}
func (h userHandler) Login(c *gin.Context) {
	var reqBody repository.User
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userServ.Authenticate(reqBody.Username, reqBody.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	middleware.CreateSessionCookie(c, reqBody.Username)
	c.IndentedJSON(http.StatusOK, user)
}
func (h userHandler) LogOut(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "log out success"})
}
func (h userHandler) Read(c *gin.Context) {
	users, err := h.userServ.Read()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}
func (h userHandler) ReadId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userServ.ReadById(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
func (h userHandler) Create(c *gin.Context) {
	var reqBody repository.User
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = h.userServ.Create(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "created success"})
}
func (h userHandler) Update(c *gin.Context) {
	var reqBody repository.User
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = h.userServ.Update(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "update success"})
}
func (h userHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = h.userServ.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusNoContent, gin.H{"message": "deleted success"})
}
