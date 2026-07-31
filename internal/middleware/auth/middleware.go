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
		// 1. Проверяем наличие заголовка
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "authorization header required")
			c.Abort() // Останавливаем выполнение
			return
		}

		// 2. Проверяем формат
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, "invalid authorization format, expected Bearer token")
			c.Abort()
			return
		}

		// 3. Извлекаем токен
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "token is empty")
			c.Abort()
			return
		}

		// 4. Валидируем токен
		claims, err := m.jwtManager.ValidateAccessToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// 5. Сохраняем пользователя в контексте (одним объектом!)
		user := ctx.UserContext{
			ID:    claims.UserID,
			Role:  claims.Role,
			Email: claims.Email, // если есть
		}
		c.Set("user", user)

		// 6. Пропускаем запрос дальше
		c.Next()
	}
}
