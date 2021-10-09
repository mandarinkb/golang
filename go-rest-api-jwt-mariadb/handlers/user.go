package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/service"
)

type userHandler struct {
	userSrv service.UserService
}

func NewUserHandler(userSrv service.UserService) userHandler {
	return userHandler{userSrv: userSrv}
}

func (h userHandler) Authenticate(c *gin.Context) {
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	token, err := h.userSrv.Authenticate(reqBody.Username, reqBody.Password)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, token)
}

func (h userHandler) ReadUsers(c *gin.Context) {
	users, err := h.userSrv.Read()
	if err != nil {
		fmt.Println(err)
		// gin.H{} สามารถใส่ข้อความแล้วจะ return json ให้
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	// return ค่า http status และค่า json obj หรือ json array ออกให้
	c.IndentedJSON(http.StatusOK, users)
}

func (h userHandler) ReadUserByID(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userSrv.ReadById(id)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func (h userHandler) CreateUsers(c *gin.Context) {
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	var reqBody service.UserRequest
	// แปลงค่าจาก body payload เป็น struct
	err := c.BindJSON(&reqBody)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userSrv.Create(reqBody)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, user)
}

func (h userHandler) UpdateUsers(c *gin.Context) {
	var reqBody service.UserRequest
	err := c.BindJSON(&reqBody)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	user, err := h.userSrv.Update(reqBody)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func (h userHandler) DeleteUsers(c *gin.Context) {
	idStr := c.Param("id")
	// แปลงเป็น int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	err = h.userSrv.Delete(id)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "delete user success"})
}
