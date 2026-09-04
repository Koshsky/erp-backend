package auth

import (
	"context"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// PrincipalResolver resolves the caller's current effective rights by user id
// (admin bypass, assigned preset, per-user overrides). Implemented by the
// rbacpolicy PolicyStore — an in-memory snapshot, never a per-request DB call.
type PrincipalResolver interface {
	EffectiveUser(ctx context.Context, userID int64) (userctx.UserContext, error)
}

type Middleware struct {
	logger     *slog.Logger
	jwtManager *jwt.Service
	resolver   PrincipalResolver
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

		// 5. Resolve the caller's current rights from the RBAC snapshot (so
		// permission changes apply immediately, without waiting for token
		// expiry). A missing entry falls back to the default deny principal.
		user, err := m.resolver.EffectiveUser(c.Request.Context(), claims.UserID)
		if err != nil {
			if m.logger != nil {
				m.logger.ErrorContext(
					c.Request.Context(),
					"auth: не удалось получить права пользователя",
					"user_id",
					claims.UserID,
					"error",
					err,
				)
			}
			user = userctx.UserContext{ID: claims.UserID}
		}
		user.ID = claims.UserID
		user.Email = claims.Email

		// 6. Store the user in the context (as a single object!)
		c.Set("user", user)
		tracingpkg.SetUserOnSpan(c, user.ID, user.Preset)

		// 7. Pass the request through
		c.Next()
	}
}
