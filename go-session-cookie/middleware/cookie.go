package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// เก็บ session ในตัวแปร สามารถดัดแปลงเก็บที่อื่นได้เช่น redis
var Sessions = map[string]Session{}
var (
	UserSession  Session
	sessionToken string
)

type Session struct {
	Username string    `json:"username"`
	Expire   time.Time `json:"expire"`
}

func (s Session) isExpired() bool {
	return s.Expire.Before(time.Now())
}

func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path != "/login" {
			var err error
			// รับ cookie
			sessionToken, err = c.Cookie("session_token")
			if err != nil {
				// กรณีเรียก path  /refresh แต่ cookie หมดอายุให้ขอ refresh cookie ใหม่
				if path == "/refresh" {
					c.Next()
					return
				}
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				c.Abort()
				return
			}
			// ตรวจสอบ cookie ที่เคยเก็บไว้ในระบบ
			var exists bool
			UserSession, exists = Sessions[sessionToken]
			if !exists {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "session cookie not found in store"})
				c.Abort()
				return
			}
			// ตรวจสอบหมดอายุ cookie
			if UserSession.isExpired() {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "session cookie expires"})
				c.Abort()
				return
			}
			if path == "/logout" {
				delete(Sessions, sessionToken)
			}
			c.Next()
		}
		c.Next()
	}
}
