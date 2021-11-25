package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/database"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/handlers"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/middleware"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/service"
)

func main() {
	db, err := database.Conn()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	productRepo := repository.NewProductRepo(db)
	productSrv := service.NewProductServ(productRepo)
	productHandler := handlers.NewProductHandler(productSrv)

	userRepo := repository.NewUserRepo(db)
	userSrv := service.NewUserServ(userRepo)
	userHandler := handlers.NewUserHandler(userSrv)

	// set release mode
	// using env:   export GIN_MODE=release
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// ใช้ middleware CORS ที่เขียนขึ้นมาเอง
	router.Use(middleware.CORS())

	// ใช้ middleware ที่เขียนขึ้นมาเอง
	// โดยจะต้องแนบ token มาถึงจะเรียกใช้งาน api ได้
	router.Use(middleware.JWTAuth)

	// Simple grouping routes: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/token/refresh", middleware.JWTRefresh)
		v1.GET("/products", productHandler.SearchProduct)
		v1.GET("/products-pagination", productHandler.PaginationProduct)
		v1.POST("/authenticate", userHandler.Authenticate)
		v1.POST("/logout", userHandler.Logout)
		v1.GET("/users", userHandler.ReadUsers)
		v1.GET("/users/:id", userHandler.ReadUserByID)
		v1.POST("/users", userHandler.CreateUsers)
		v1.PUT("/users", userHandler.UpdateUsers)
		v1.DELETE("/users/:id", userHandler.DeleteUsers)
	}
	router.Run(":8080")
}
