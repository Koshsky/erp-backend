package auth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyRole   = "role"
	ContextKeyUserID = "user_id"
)

type AuthMiddleware struct {
	logger        *slog.Logger
	jwtManager    *jwt.JWTService
	DefaultRole   string
	DefaultUserID int64
}

func NewAuthMiddleware(logger *slog.Logger, jwtManager *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:        logger,
		jwtManager:    jwtManager,
		DefaultRole:   "ДП",
		DefaultUserID: 1,
	}
}

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		reqCtx := c.Request.Context()

		role := m.DefaultRole
		userID := m.DefaultUserID

		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := m.jwtManager.ValidateAccessToken(tokenString)
			if err != nil {
				m.logger.Warn("invalid access token", "error", err)
			} else {
				role = claims.Role
				userID = claims.UserID
			}
		}

		reqCtx = context.WithValue(reqCtx, ContextKeyRole, role)
		reqCtx = context.WithValue(reqCtx, ContextKeyUserID, userID)
		fmt.Println()
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		reqCtx := c.Request.Context()
		reqCtx = context.WithValue(reqCtx, ContextKeyUserID, claims.UserID)
		reqCtx = context.WithValue(reqCtx, ContextKeyRole, claims.Role)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}

// GetRole извлекает роль из контекста запроса (для использования в репозиториях/сервисах)
func GetRole(ctx context.Context) string {
	role, ok := ctx.Value(ContextKeyRole).(string)
	if !ok {
		return ""
	}
	return role
}

// GetUserID извлекает user_id из контекста запроса (для использования в репозиториях/сервисах)
func GetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(ContextKeyUserID).(int64)
	if !ok {
		return 0
	}
	return userID
}
