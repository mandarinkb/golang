package middleware

import (
	"github.com/valyala/fasthttp"
)

var (
	corsAllowCredentials = "true"
	corsAllowHeaders     = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
	corsAllowMethods     = "POST, OPTIONS, GET, PUT, DELETE"
	corsAllowOrigin      = "http://127.0.0.1:5501" //http://127.0.0.1:5501
)

func cors(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
	ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
	ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
}
