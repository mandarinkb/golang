package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/handlers"
	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/middleware"
	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/repository"
	"github.com/mandarinkb/go-rest-api-with-fasthttp-and-fasthttp-routing/service"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func main() {
	db, err := sql.Open("mysql", "root:mandarinkb@tcp(127.0.0.1)/TEST?charset=utf8")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	userServ := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userServ)

	router := routing.New()
	router.Use(middleware.CORS())
	router.Use(middleware.Cookie())

	router.Post("/login", userHandler.Authenticate)
	router.Post("/logout", userHandler.Authenticate)
	router.Get("/users", userHandler.Read)
	router.Get("/users/<id>", userHandler.ReadById)
	router.Post("/users", userHandler.Create)
	router.Put("/users", userHandler.Update)
	router.Delete("/users/<id>", userHandler.Delete)

	fmt.Println("service run port 8989")
	fasthttp.ListenAndServe(":8989", router.HandleRequest)
}
