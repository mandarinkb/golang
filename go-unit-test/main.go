package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-unit-test/handler"
	"github.com/mandarinkb/go-unit-test/repository"
	"github.com/mandarinkb/go-unit-test/service"
)

func main() {
	userRepo := repository.NewUser()
	userServ := service.NewUser(userRepo)
	userHandler := handler.NewUserHandler(userServ)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/users", userHandler.GetUser)
	router.POST("/users", userHandler.CreateUser)
	router.PUT("/users", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
	router.Run(":8080")
}
