package middleware

import (
	routing "github.com/qiangxue/fasthttp-routing"
)

var (
	corsAllowCredentials = "true"
	corsAllowHeaders     = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
	corsAllowMethods     = "POST, OPTIONS, GET, PUT, DELETE"
	corsAllowOrigin      = "*" //http://127.0.0.1:5501
)

func CORS() routing.Handler {
	return func(c *routing.Context) error {
		c.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		c.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		c.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		c.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)

		method := string(c.Method())
		if method == "OPTIONS" {
			c.SetStatusCode(204)
			c.Abort()
			return nil
		}
		c.Next()
		return nil
	}
}
