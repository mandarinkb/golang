package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-session-cookie/handlers"
	"github.com/mandarinkb/go-session-cookie/middleware"
	"github.com/mandarinkb/go-session-cookie/repository"
	"github.com/mandarinkb/go-session-cookie/service"
)

func main() {
	db, err := sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST?charset=utf8")
	if err != nil {
		panic(err)
	}
	userdb := repository.NewUserRepo(db)
	userServ := service.NewUserServ(userdb)
	userHandler := handlers.NewUserHandler(userServ)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Cookie())
	router.POST("/login", userHandler.Login)
	router.POST("/logout", userHandler.LogOut)
	router.POST("/refresh", userHandler.Refresh)
	router.GET("/users", userHandler.Read)
	router.Run(":8989")
}
