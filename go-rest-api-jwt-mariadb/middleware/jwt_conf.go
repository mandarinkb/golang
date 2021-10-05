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

func (p permitPathConfig) Path(paths ...string) bool {
	isPath := false
	for _, path := range paths {
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
