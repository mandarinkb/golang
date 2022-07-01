package middleware

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func CORE(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Println(string(ctx.Path()))

		cors(ctx)
		cookie(ctx)
		next(ctx)
	}
}
