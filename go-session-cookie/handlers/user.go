package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-session-cookie/middleware"
	"github.com/mandarinkb/go-session-cookie/service"
)

type userHandler struct {
	userServ service.UserService
}

func NewUserHandler(userServ service.UserService) userHandler {
	return userHandler{userServ}
}

func (h userHandler) Login(c *gin.Context) {
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err = h.userServ.Authenticate(reqBody.Username, reqBody.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	middleware.CreateSessionCookie(c, reqBody.Username)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "authenticate success"})
}

func (h userHandler) LogOut(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "log out success"})
}
func (h userHandler) Refresh(c *gin.Context) {
	middleware.CreateSessionCookie(c, middleware.RSession.Username)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "create new cookie success"})
}

func (h userHandler) Read(c *gin.Context) {
	users, err := h.userServ.Read()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func (h userHandler) ReadById(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
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
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userServ.Create(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
func (h userHandler) Update(c *gin.Context) {
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userServ.Update(reqBody)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
func (h userHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userServ.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}
