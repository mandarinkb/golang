package middleware

import (
	"encoding/json"
	"time"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

type response struct {
	Message string `json:"message"`
}

func message(msg string) []byte {
	message := response{
		Message: msg,
	}
	resp, _ := json.Marshal(message)
	return resp
}
func Cookie() routing.Handler {
	return func(c *routing.Context) error {
		if string(c.Path()) == "/login" {
			c.Next()
			return nil
		}
		cookieValue := string(c.Request.Header.Cookie("access_token"))
		if cookieValue == "" {
			c.Response.Header.Set("Content-Type", "application/json")
			c.SetStatusCode(401)
			c.Write(message("cookie session is expired"))
			c.Abort()
			return nil
		}
		c.Next()
		return nil
	}
}

func SetCookie(c *routing.Context, value string) {
	var fc fasthttp.Cookie
	fc.SetKey("access_token")
	fc.SetValue(value)
	fc.SetExpire(time.Now().Add(1 * time.Minute))
	fc.SetHTTPOnly(true)
	c.Response.Header.SetCookie(&fc)
}
