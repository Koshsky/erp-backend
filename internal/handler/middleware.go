package handler

import (
	"github.com/gin-gonic/gin"
)

// AuthMiddleware — заглушка для авторизации.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("role", "ДП")
		c.Set("user_id", int64(1))
		c.Next()
	}
}
