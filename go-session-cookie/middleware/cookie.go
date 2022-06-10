package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// เก็บ AccessSessions RefreshSessions ในตัวแปร สามารถดัดแปลงเก็บที่อื่นได้เช่น redis
var (
	AccessSessions  = map[string]Session{}
	RefreshSessions = map[string]Session{}
	ASession        Session
	RSession        Session
	Username        string
)

type Session struct {
	Username        string    `json:"username"`
	Expire          time.Time `json:"expire"`
	RefAccessToken  string    `json:"refAccessToken"`
	RefRefreshToken string    `json:"refRefreshToken"`
}

func (s Session) isExpired() bool {
	return s.Expire.Before(time.Now())
}

func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		switch path {
		case "/login":
			c.Next()
			return
		case "/refresh":
			rToken, err := c.Cookie("refresh_token")
			if err != nil {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				c.Abort()
				return
			}
			// ตรวจสอบ cookie ที่เคยเก็บไว้ในระบบ
			var exists bool
			RSession, exists = RefreshSessions[rToken]
			if !exists {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "refresh cookie not found in store"})
				c.Abort()
				return
			}
			delete(AccessSessions, RSession.RefAccessToken)
			delete(RefreshSessions, rToken)
			c.Next()
			return
		}
		// รับ cookie
		aToken, err := c.Cookie("access_token")
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}
		// ตรวจสอบ cookie ที่เคยเก็บไว้ในระบบ
		var exists bool
		ASession, exists = AccessSessions[aToken]
		if !exists {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "session cookie not found in store"})
			c.Abort()
			return
		}
		// ตรวจสอบหมดอายุ cookie
		if ASession.isExpired() {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "session cookie expires"})
			c.Abort()
			return
		}
		if path == "/logout" {
			delete(AccessSessions, aToken)
			delete(RefreshSessions, ASession.RefRefreshToken)
		}
		c.Next()
	}
}
func CreateSessionCookie(c *gin.Context, username string) {
	atExpiresAt := time.Now().Add(1 * time.Minute)
	rtExpiresAt := time.Now().Add(60 * time.Minute)
	accessToken := uuid.NewString()
	refreshToken := uuid.NewString()
	//access cookie
	AccessSessions[accessToken] = Session{
		Username:        username,
		Expire:          atExpiresAt,
		RefRefreshToken: refreshToken,
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  atExpiresAt,
		HttpOnly: true})

	//refresh cookie
	RefreshSessions[refreshToken] = Session{
		Username:       username,
		Expire:         rtExpiresAt,
		RefAccessToken: accessToken,
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  rtExpiresAt,
		HttpOnly: true})
}
