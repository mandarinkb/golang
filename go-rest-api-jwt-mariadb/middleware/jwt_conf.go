package middleware

import (
	"github.com/gin-gonic/gin"
)

type permitPathConfig struct {
	c *gin.Context
}

func NewPermitPathConfig(c *gin.Context) permitPathConfig {
	return permitPathConfig{c: c}
}

// ไว้ใช้อนุญาต path ที่ไม่ต้องผ่าน การ Authenticate
func (p permitPathConfig) Path(paths ...string) bool {
	isPath := false
	for _, path := range paths {
		// ตรวจสอบ path จาก url และ path ที่อนุญาต
		if p.c.Request.URL.Path == path {
			isPath = true
		}
	}
	if isPath {
		return true
	} else {
		return false
	}
}
