package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fasthttp/router"
	"github.com/mandarinkb/go-rest-api-fasthttp/middleware"
	"github.com/mandarinkb/go-rest-api-fasthttp/repository"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type response struct {
	Message string
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func read(ctx *fasthttp.RequestCtx) {
	r := response{
		Message: "hello",
	}
	resp, err := json.Marshal(r)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(200)
	ctx.Write(resp)
}
func create(ctx *fasthttp.RequestCtx) {
	var reqBody user
	err := json.Unmarshal(ctx.PostBody(), &reqBody)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	fmt.Println(reqBody)

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("access_token")
	c.SetValue("cookie-value")
	c.SetExpire(time.Now().Add(1 * time.Minute))
	c.SetHTTPOnly(true)
	ctx.Response.Header.SetCookie(&c)
}

func main() {
	dsn := "root:mandarinkb@tcp(127.0.0.1)/TEST?charset=utf8"
	dial := mysql.Open(dsn)
	db, err := gorm.Open(dial)
	if err != nil {
		panic(err)
	}
	// จะสร้าง table ขึ้นมาถ้ามีอยู่แล้วจะไม่สร้างซ้ำ
	db.AutoMigrate(repository.User{})
	router := router.New()
	router.GET("/test", read)
	router.POST("/test", create)
	log.Fatal(fasthttp.ListenAndServe(":8080", middleware.CORE(router.Handler)))
}
