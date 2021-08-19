package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// type user struct {
// 	USER_ID  int    `json:"userId"`
// 	USERNAME string `json:"username"`
// 	PASSWORD string `json:"password"`
// 	ROLE     string `json:"role"`
// }

// var users = []user{
// 	{USER_ID: 1, USERNAME: "test", PASSWORD: "test-password", ROLE: "admin"},
// 	{USER_ID: 2, USERNAME: "test1", PASSWORD: "test1-password", ROLE: "assistant"},
// 	{USER_ID: 3, USERNAME: "test2", PASSWORD: "test2-password", ROLE: "admin"},
// 	{USER_ID: 4, USERNAME: "test3", PASSWORD: "test3-password", ROLE: "assistant"},
// }

// // getUsers responds with the list of all users as JSON.
// func getUsers(c *gin.Context) {
// 	c.IndentedJSON(http.StatusOK, users)
// }

// // postUsers adds an user from JSON received in the request body.
// func postUsers(c *gin.Context) {
// 	var newUser user

// 	// Call BindJSON to bind the received JSON to
// 	// newUser.
// 	err := c.BindJSON(&newUser)
// 	if err != nil {
// 		return
// 	}

// 	// Add the new user to the slice.
// 	users = append(users, newUser)
// 	c.IndentedJSON(http.StatusCreated, newUser)
// }

// // getUserByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getUserByID(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Loop through the list of users, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range users {
// 		if a.USER_ID == id {
// 			c.IndentedJSON(http.StatusOK, a)
// 			return
// 		}
// 	}
// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
// }

//

/*****************************
***  for without framwork  ***
*****************************/

type user struct {
	USER_ID  int    `json:"userId"`
	USERNAME string `json:"username"`
	PASSWORD string `json:"password"`
	ROLE     string `json:"role"`
}

var users = []user{
	{USER_ID: 1, USERNAME: "test", PASSWORD: "test-password", ROLE: "admin"},
	{USER_ID: 2, USERNAME: "test1", PASSWORD: "test1-password", ROLE: "assistant"},
	{USER_ID: 3, USERNAME: "test2", PASSWORD: "test2-password", ROLE: "admin"},
	{USER_ID: 4, USERNAME: "test3", PASSWORD: "test3-password", ROLE: "assistant"},
}

type respose struct {
	Message string `json:"message"`
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}
func homePage(w http.ResponseWriter, r *http.Request) {
	res := respose{
		Message: "Hwllo rest api",
	}
	json.NewEncoder(w).Encode(res)
}
func handleRequest() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", getAllUsers)
	http.ListenAndServe(":8080", nil)
}
func main() {
	// router := gin.Default()
	// router.GET("/users", getUsers)
	// router.GET("/users/:id", getUserByID)
	// router.POST("/users", postUsers)
	// router.Run("localhost:8080")
	fmt.Println("service running at port 8080")
	handleRequest()
}
