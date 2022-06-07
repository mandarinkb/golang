package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	expiresAt := time.Now().Add(1 * time.Minute)
	sessionToken := uuid.NewString()
	middleware.Sessions[sessionToken] = middleware.Session{
		Username: reqBody.Username,
		Expire:   expiresAt,
	}
	fmt.Println(middleware.Sessions)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true})
}

func (h userHandler) LogOut(c *gin.Context) {
	fmt.Println(middleware.Sessions)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "log out success"})
}
func (h userHandler) Refresh(c *gin.Context) {
	fmt.Println(middleware.UserSession.Username)
	expiresAt := time.Now().Add(1 * time.Minute)
	sessionToken := uuid.NewString()
	middleware.Sessions[sessionToken] = middleware.Session{
		Username: middleware.UserSession.Username,
		Expire:   expiresAt,
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true})
}

func (h userHandler) Read(c *gin.Context) {
	users, err := h.userServ.Read()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}
