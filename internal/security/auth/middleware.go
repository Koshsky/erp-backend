package auth

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyRole   = "role"
	ContextKeyUserID = "user_id"
)

type AuthMiddleware struct {
	DefaultRole   string
	DefaultUserID int64
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		DefaultRole:   "ДП",
		DefaultUserID: 1,
	}
}

// Middleware устанавливает роль и user_id в контекст запроса.
// В будущем будет извлекать их из JWT claims.
func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: парсить JWT, извлекать role и user_id из claims
		// claims, err := parseJWT(token)
		// c.Set(ContextKeyRole, claims.Role)
		// c.Set(ContextKeyUserID, claims.UserID)

		// Значения по умолчанию, пока нет JWT
		ctx := context.WithValue(c.Request.Context(), ContextKeyRole, "ДП")
		ctx = context.WithValue(ctx, ContextKeyUserID, int64(1))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// GetRole извлекает роль из контекста gin
func GetRole(c *gin.Context) string {
	role, _ := c.Get(ContextKeyRole)
	if role == nil {
		return "admin"
	}
	return role.(string)
}

// GetUserID извлекает user_id из контекста gin
func GetUserID(c *gin.Context) int64 {
	userID, _ := c.Get(ContextKeyUserID)
	if userID == nil {
		return 1
	}
	return userID.(int64)
}
