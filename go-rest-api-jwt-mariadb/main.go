package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/database"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/handlers"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/repository"
	"github.com/mandarinkb/go-rest-api-jwt-mariadb/service"
)

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "mandarinkb", // no password set
		DB:       0,            // use default DB
	})

	exp := time.Now().Add(time.Minute * 20).Unix()
	at := time.Unix(exp, 0) //converting Unix to UTC(to Time object)
	now := time.Now()

	accessUuid := uuid.New()
	fmt.Println(accessUuid.String())

	err := rdb.Set(ctx, accessUuid.String(), "41sarero2584", at.Sub(now)).Err()
	if err != nil {
		fmt.Println(err)
	}

	// err := rdb.Del(ctx, "5467b86b-7622-475b-94be-6ee9496ade31")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// val, err := rdb.Get(ctx, "8b26b0b2-7b58-4c7f-bac1-5f2aeb1941e6").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("key: ", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
	// Output: key value
	// key2 does not exist
}
func main() {
	db, err := database.NewDatabase().Conn()
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
	// using code:  gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS gin's middleware Default() allows all origins
	// router.Use(cors.Default())

	// ใช้ middleware ที่เขียนขึ้นมาเอง
	// โดยจะต้องแนบ token มาถึงจะเรียกใช้งาน api ได้
	router.Use(handlers.NewJwtHandler().JWTAuth)

	// Simple grouping routes: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/token/refresh", handlers.NewJwtHandler().JWTRefresh)
		v1.GET("/products", productHandler.SearchProduct)
		v1.GET("/products-pagination", productHandler.PaginationProduct)
		v1.POST("/authenticate", userHandler.Authenticate)
		v1.GET("/users", userHandler.ReadUsers)
		v1.GET("/users/:id", userHandler.ReadUserByID)
		v1.POST("/users", userHandler.CreateUsers)
		v1.PUT("/users", userHandler.UpdateUsers)
		v1.DELETE("/users/:id", userHandler.DeleteUsers)
	}
	router.Run(":8080")
}
