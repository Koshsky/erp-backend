package auth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Koshsky/erp-backend/internal/common/ctx"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/gin-gonic/gin"
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
		DefaultRole:   "unauthenticated",
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

		reqCtx = context.WithValue(reqCtx, ctx.ContextKeyRole, role)
		reqCtx = context.WithValue(reqCtx, ctx.ContextKeyUserID, userID)
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
		reqCtx = context.WithValue(reqCtx, ctx.ContextKeyUserID, claims.UserID)
		reqCtx = context.WithValue(reqCtx, ctx.ContextKeyRole, claims.Role)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}
