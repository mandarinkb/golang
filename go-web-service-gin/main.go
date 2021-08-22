package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	USER_ID  int    `json:"userId"`
	USERNAME string `json:"username"`
	PASSWORD string `json:"password"`
	ROLE     string `json:"role"`
}

var users = []User{
	{USER_ID: 1, USERNAME: "test", PASSWORD: "test-password", ROLE: "admin"},
	{USER_ID: 2, USERNAME: "test2", PASSWORD: "test2-password", ROLE: "assistant"},
	{USER_ID: 3, USERNAME: "test3", PASSWORD: "test3-password", ROLE: "admin"},
	{USER_ID: 4, USERNAME: "test4", PASSWORD: "test4-password", ROLE: "assistant"},
}

var user = User{USER_ID: 1, USERNAME: "test", PASSWORD: "test-password", ROLE: "admin"}

func readUsers(c *gin.Context) {
	// query := c.Query("q") // shortcut for c.Request.URL.Query().Get("q")
	// fmt.Println(query)
	c.IndentedJSON(http.StatusOK, users) // return ค่า http status และค่า json obj หรือ json array ออกให้เลย
}

func readUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // แปลงเป็น int
	if err != nil {
		fmt.Println(err)
	}
	// ค้นหา id  ที่ตรงกับ mock data
	for _, user := range users {
		if user.USER_ID == id {
			c.IndentedJSON(http.StatusOK, user)
			return //ถ้าพบก็ออกจาก for loop ทันที
		}
	}
	// กรณี id ที่ค้นหาไม่พบ
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"}) // gin.H{} สามารถใส่ json ได้เลย
}

func createUsers(c *gin.Context) {
	var reqBody User
	// request body จะถูกเรียกเมื่อ ใช้ ฟังก์ชัน BindJSON โดยค่าใน body จะตรงกับ reqBody
	// reqBody ชื่อ filed ในนี้จะต้องตรงกับ ค่า key ใน body ด้วย
	err := c.BindJSON(&reqBody) // แปลงค่าจาก body payload เป็น struct
	if err != nil {
		return
	}
	// Add the new user to the slice.
	users = append(users, reqBody)
	c.IndentedJSON(http.StatusCreated, reqBody)
}

func updateUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // แปลงเป็น int
	if err != nil {
		fmt.Println(err)
	}

	var reqBody User
	var newUser User
	var sliceUser []User         // ไว้ อัพเดท slice User
	err2 := c.BindJSON(&reqBody) // แปลงค่าจาก body payload เป็น struct
	if err2 != nil {
		return
	}
	// ตรวจเช็ค id ที่รับเข้ามาก่อน
	haveId := false
	for _, user := range users {
		if user.USER_ID == id {
			haveId = true
		}
	}
	// ค้นหา id  ที่ตรงกับ mock data
	if haveId {
		for i, user := range users {
			if user.USER_ID == id {
				newUser = User{
					USER_ID:  users[i].USER_ID,
					USERNAME: reqBody.USERNAME,
					PASSWORD: reqBody.PASSWORD,
					ROLE:     reqBody.ROLE,
				}
				sliceUser = append(sliceUser, newUser)

			} else {
				newUser = User{
					USER_ID:  users[i].USER_ID,
					USERNAME: users[i].USERNAME,
					PASSWORD: users[i].PASSWORD,
					ROLE:     users[i].ROLE,
				}
				sliceUser = append(sliceUser, newUser)
			}
		}
		users = sliceUser
		c.IndentedJSON(http.StatusOK, reqBody)
		return //ถ้าพบก็ออกจาก for loop ทันที
	}
	// กรณี id ที่ค้นหาไม่พบ
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func deleteUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr) // แปลงเป็น int
	if err != nil {
		fmt.Println(err)
	}
	var newUser User
	var sliceUser []User // ไว้ อัพเดท slice User
	// ตรวจเช็ค id ที่รับเข้ามาก่อน
	haveId := false
	for _, user := range users {
		if user.USER_ID == id {
			haveId = true
		}
	}
	// ค้นหา id  ที่ตรงกับ mock data
	if haveId {
		for i, user := range users {
			if user.USER_ID != id {
				newUser = User{
					USER_ID:  users[i].USER_ID,
					USERNAME: users[i].USERNAME,
					PASSWORD: users[i].PASSWORD,
					ROLE:     users[i].ROLE,
				}
				sliceUser = append(sliceUser, newUser)
			}
		}
		users = sliceUser
		c.IndentedJSON(http.StatusOK, gin.H{"message": "delete user success"})
		return //ถ้าพบก็ออกจาก for loop ทันที
	}
	// กรณี id ที่ค้นหาไม่พบ
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func main() {
	fmt.Println("Listening and serving HTTP on :8080")

	// set release mode
	// using env:   export GIN_MODE=release
	// using code:  gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS gin's middleware Default() allows all origins
	// router.Use(cors.Default())

	// Simple grouping routes: v1
	// v1 := router.Group("/v1")
	// {
	// 	v1.GET("/users", readUsers)
	// 	v1.GET("/users/:id", readUserByID)
	// 	v1.POST("/users", createUsers)
	// 	v1.PUT("/users/:id", updateUsers)
	// 	v1.DELETE("/users/:id", deleteUsers)
	// }

	router.GET("/users", readUsers)
	router.GET("/users/:id", readUserByID)
	router.POST("/users", createUsers)
	router.PUT("/users/:id", updateUsers)
	router.DELETE("/users/:id", deleteUsers)
	router.Run(":8080")
}
