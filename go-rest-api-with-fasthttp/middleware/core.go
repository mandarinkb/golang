package middleware

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type response struct {
	Message string
}

func message(msg string) []byte {
	message := response{
		Message: msg,
	}
	resp, _ := json.Marshal(message)
	return resp
}
func CORE(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		cors(ctx)
		if err := cookie(ctx); err != nil {
			ctx.Response.Header.Set("Content-Type", "application/json")
			ctx.SetStatusCode(401)
			ctx.Write(message("unauth"))
			return
		}
		next(ctx)
	}
}
