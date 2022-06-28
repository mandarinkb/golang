package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// เก็บ accessSessions RefreshSessions ในตัวแปร สามารถดัดแปลงเก็บที่อื่นได้เช่น redis
var (
	accessSessions = map[string]Session{}
	session        Session
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
		// ไม่ต้องตรวจสอบ cookie ตอน login
		if path == "/login" {
			c.Next()
			return
		}
		// รับ cookie
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		// ตรวจสอบ cookie ที่เคยเก็บไว้ในระบบ
		var exists bool
		session, exists = accessSessions[accessToken]
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "session cookie not found in store"})
			c.Abort()
			return
		}
		// ตรวจสอบหมดอายุ cookie
		if session.isExpired() {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "session cookie expires"})
			c.Abort()
			return
		}
		if path == "/logout" {
			delete(accessSessions, accessToken)
		}
		c.Next()
	}
}

func CreateSessionCookie(c *gin.Context, username string) {
	atExpiresAt := time.Now().Add(30 * time.Minute)
	accessToken := uuid.NewString()
	// เก็บ AccessSessions RefreshSessions ในตัวแปร สามารถดัดแปลงเก็บที่อื่นได้เช่น redis
	accessSessions[accessToken] = Session{
		Username: username,
		Expire:   atExpiresAt,
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  atExpiresAt,
		HttpOnly: true})
}
