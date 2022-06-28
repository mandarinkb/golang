package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-gorm-session-cookie/handlers"
	"github.com/mandarinkb/go-gorm-session-cookie/middleware"
	"github.com/mandarinkb/go-gorm-session-cookie/repository"
	"github.com/mandarinkb/go-gorm-session-cookie/service"
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

	userdb := repository.NewUserORM(db)
	userServ := service.NewUserService(userdb)
	userHandler := handlers.NewUserHandler(userServ)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Cookie())
	router.POST("/login", userHandler.Login)
	router.POST("/logout", userHandler.LogOut)
	router.GET("/users", userHandler.Read)
	router.GET("/users/:id", userHandler.ReadId)
	router.POST("/users", userHandler.Create)
	router.PUT("/users", userHandler.Update)
	router.DELETE("/users/:id", userHandler.Delete)
	router.Run(":8989")

}
