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
	logger     *slog.Logger
	jwtManager *jwt.Service
}

func NewMiddleware(logger *slog.Logger, jwtManager *jwt.Service) *Middleware {
	return &Middleware{
		logger:     logger,
		jwtManager: jwtManager,
	}
}

// RequireAuth verifies the JWT token and sets the user context.
func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check for the presence of the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header required")
			c.Abort() // Stop the execution
			return
		}

		// 2. Check the format
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, "invalid authorization format, expected Bearer token")
			c.Abort()
			return
		}

		// 3. Extract the token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "token is empty")
			c.Abort()
			return
		}

		// 4. Validate the token
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// 5. Store the user in the context (as a single object!)
		user := ctx.UserContext{
			ID:    claims.UserID,
			Role:  claims.Role,
			Email: claims.Email, // if present
		}
		c.Set("user", user)

		// 6. Pass the request through
		c.Next()
	}
}
