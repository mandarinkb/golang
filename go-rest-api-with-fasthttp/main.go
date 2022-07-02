package main

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/mandarinkb/go-rest-api-fasthttp/handlers"
	"github.com/mandarinkb/go-rest-api-fasthttp/middleware"
	"github.com/mandarinkb/go-rest-api-fasthttp/repository"
	"github.com/mandarinkb/go-rest-api-fasthttp/service"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:mandarinkb@tcp(127.0.0.1)/TEST?charset=utf8"
	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial)
	if err != nil {
		panic(err)
	}
	// จะสร้าง table ขึ้นมาถ้ามีอยู่แล้วจะไม่สร้างซ้ำ
	db.AutoMigrate(repository.User{})

	userRepo := repository.NewUserORM(db)
	userServ := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userServ)

	router := router.New()
	router.POST("/login", userHandler.Authenticate)
	router.POST("/logout", userHandler.LogOut)
	router.GET("/users", userHandler.Read)
	router.GET("/users/{id}", userHandler.ReadById)
	router.POST("/users", userHandler.Create)
	router.PUT("/users", userHandler.Update)
	router.DELETE("/users/{id}", userHandler.Delete)

	fmt.Println("service is start port:8989")
	fasthttp.ListenAndServe(":8989", middleware.CORE(router.Handler))
}
