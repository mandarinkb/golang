package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type permitPathConfig struct {
	c *gin.Context
}

func NewPermitPathConfig(c *gin.Context) permitPathConfig {
	return permitPathConfig{c: c}
}

// ไว้ใช้อนุญาต path ที่ไม่ต้องผ่าน การ Authenticate
// สามารถกำหนด path แบบ /.../** ได้  เช่น /example/** คือ ขึ้นต้น /example/  ด้านหลังเป็น endpoint อะไรก็ได้
// ตัวอย่าง /example/1  /example/users
func (p permitPathConfig) Path(paths ...string) bool {
	isPath := false
	for _, path := range paths {
		// ตรวจสอบ path จาก url และ path ที่อนุญาต
		if p.c.Request.URL.Path == path {
			isPath = true
		}

		// ตรวจสอบ "/**" มีอยู่ใน path ที่กำหนดหรือไม่
		isMatch := strings.Contains(path, "/**")
		if isMatch {
			// ตัดท้าย  "/**" ออก
			trimPath := strings.TrimSuffix(path, "/**")

			// ตรวจสอบ trim path มีตัวอักษรอยู่ใน path จาก request url หรือไม่
			isContainsReqPath := strings.Contains(p.c.Request.URL.Path, trimPath)
			if isContainsReqPath {
				isPath = true
			}
		}
	}

	if isPath {
		return true
	} else {
		return false
	}
}
