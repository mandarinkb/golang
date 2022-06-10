package middleware

import "github.com/gin-gonic/gin"

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ข้อสำคัญต้องระบุ "Access-Control-Allow-Origin จาก domain ที่ร้องขอข้อมูล จะใส่ * ไม่ได้
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5501") //http://127.0.0.1:5501
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
