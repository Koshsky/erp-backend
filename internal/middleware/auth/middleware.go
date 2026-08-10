package auth

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type Middleware struct {
	logger     *slog.Logger
	jwtManager *jwt.Service
}

// RequireAuth verifies the JWT token and sets the user context.
func (m *Middleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Check for the presence of the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, errors.CodeUnauthorized, "authorization header required")
			c.Abort() // Stop the execution
			return
		}

		// 2. Check the format
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, errors.CodeInvalidToken, "invalid authorization format, expected Bearer token")
			c.Abort()
			return
		}

		// 3. Extract the token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, errors.CodeInvalidToken, "token is empty")
			c.Abort()
			return
		}

		// 4. Validate the token
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, errors.CodeInvalidToken, "invalid or expired token")
			c.Abort()
			return
		}

		// 5. Store the user in the context (as a single object!)
		user := userctx.UserContext{
			ID:    claims.UserID,
			Role:  claims.Role,
			Email: claims.Email, // if present
		}
		c.Set("user", user)

		// 6. Pass the request through
		c.Next()
	}
}
