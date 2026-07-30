package auth

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/ctx"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
)

type Middleware struct {
	logger        *slog.Logger
	jwtManager    *jwt.Service
	DefaultRole   string
	DefaultUserID int64
}

func NewMiddleware(logger *slog.Logger, jwtManager *jwt.Service) *Middleware {
	return &Middleware{
		logger:        logger,
		jwtManager:    jwtManager,
		DefaultRole:   "unauthenticated",
		DefaultUserID: 1,
	}
}

func (m *Middleware) Middleware() gin.HandlerFunc {
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

		reqCtx = ctx.SetRole(reqCtx, role)
		reqCtx = ctx.SetUserID(reqCtx, userID)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}

func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			return
		}

		reqCtx := c.Request.Context()
		reqCtx = ctx.SetUserID(reqCtx, claims.UserID)
		reqCtx = ctx.SetRole(reqCtx, claims.Role)
		c.Request = c.Request.WithContext(reqCtx)

		c.Next()
	}
}
