package middleware

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func cookie(ctx *fasthttp.RequestCtx) {
	ctx.Request.Header.Cookie("access_token")
	cookieValue := string(ctx.Request.Header.Cookie("access_token"))
	if cookieValue == "" {
		fmt.Println("cookie is expired")
	} else {
		fmt.Println(cookieValue)
	}

}
