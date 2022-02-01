package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mandarinkb/go-api-project-final/database"
	"github.com/mandarinkb/go-api-project-final/handlers"
	"github.com/mandarinkb/go-api-project-final/middleware"
	"github.com/mandarinkb/go-api-project-final/repository"
	"github.com/mandarinkb/go-api-project-final/service"
)

func main() {
	db, err := database.Conn()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	userSrv := service.NewUserServ(userRepo)
	userHandler := handlers.NewUserHandler(userSrv)

	scheduleRepo := repository.NewScheduleRepo(db)
	scheduleServ := service.NewScheduleService(scheduleRepo)
	scheduleHandler := handlers.NewScheduleHandler(scheduleServ)

	swDbRepo := repository.NewSwitchDBRepo(db)
	swDbServ := service.NewSwDatabaseService(swDbRepo)
	swDbHandler := handlers.NewSwitchDabaseHandler(swDbServ)

	webRepo := repository.NewWebRepo(db)
	webServ := service.NewWebService(webRepo)
	webHandler := handlers.NewWebHandler(webServ)

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
		v1.POST("/authenticate", userHandler.Authenticate)
		v1.POST("/logout", userHandler.Logout)
		v1.GET("/users", userHandler.ReadUsers)
		v1.GET("/users/:id", userHandler.ReadUserByID)
		v1.POST("/users", userHandler.CreateUsers)
		v1.PUT("/users", userHandler.UpdateUsers)
		v1.DELETE("/users/:id", userHandler.DeleteUsers)

		v1.GET("/schedule", scheduleHandler.Read)
		v1.GET("/schedule/:id", scheduleHandler.ReadById)
		v1.POST("/schedule", scheduleHandler.Create)
		v1.PUT("/schedule", scheduleHandler.Update)
		v1.DELETE("/schedule/:id", scheduleHandler.Delete)

		v1.GET("/switch-database", swDbHandler.Read)
		v1.GET("/switch-database/:id", swDbHandler.ReadById)
		v1.POST("/switch-database", swDbHandler.Create)
		v1.PUT("/switch-database", swDbHandler.Update)
		v1.PUT("/switch-database-status", swDbHandler.UpdateStatus)
		v1.DELETE("/switch-database/:id", swDbHandler.Delete)

		v1.GET("/web", webHandler.Read)
		v1.GET("/web/:id", webHandler.ReadById)
		v1.POST("/web", webHandler.Create)
		v1.PUT("/web", webHandler.Update)
		v1.PUT("/web-status", webHandler.UpdateStatus)
		v1.DELETE("/web/:id", webHandler.Delete)
	}
	router.Run(":8080")
}
