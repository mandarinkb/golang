package middleware

import (
	"errors"
	"time"

	"github.com/valyala/fasthttp"
)

func cookie(ctx *fasthttp.RequestCtx) error {
	// กรณี path login ไม่ต้องใช้ cookie
	if string(ctx.Path()) == "/login" {
		return nil
	}
	// กรณีเกิด CORS จะเรียกใช้ method OPTION
	// (fix bug) ต้องทำการดักไม่ให้ใช้ method OPTION มิฉะนั้นจะ Create Update Delete ไม่ได้
	method := string(ctx.Method())
	if method != "OPTIONS" {
		cookieValue := string(ctx.Request.Header.Cookie("access_token"))
		if cookieValue == "" {
			return errors.New("cookie session is expired")
		}
	}
	return nil
}

func SetCookie(ctx *fasthttp.RequestCtx, value string) {
	var c fasthttp.Cookie
	c.SetKey("access_token")
	c.SetValue(value)
	c.SetExpire(time.Now().Add(1 * time.Minute))
	c.SetHTTPOnly(true)
	ctx.Response.Header.SetCookie(&c)
}
