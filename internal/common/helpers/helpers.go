package helpers

import (
	"errors"

	"github.com/Koshsky/erp-backend/internal/common/ctx"

	"github.com/gin-gonic/gin"
)

// GetUser извлекает пользователя из контекста
func GetUser(c *gin.Context) (ctx.UserContext, error) {
	val, exists := c.Get(ctx.KeyUser)
	if !exists {
		return ctx.UserContext{}, errors.New("user not found in context")
	}

	user, ok := val.(ctx.UserContext)
	if !ok {
		return ctx.UserContext{}, errors.New("invalid user context type")
	}

	return user, nil
}

// MustGetUser - паникует если нет пользователя
func MustGetUser(c *gin.Context) ctx.UserContext {
	user, err := GetUser(c)
	if err != nil {
		panic(err)
	}
	return user
}

// GetUserID возвращает ID пользователя
func GetUserID(c *gin.Context) (int64, error) {
	user, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetUserRole возвращает роль пользователя
func GetUserRole(c *gin.Context) (string, error) {
	user, err := GetUser(c)
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

// IsAdmin проверяет, является ли пользователь админом
func IsAdmin(c *gin.Context) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}
	return user.Role == "admin"
}

// HasRole проверяет наличие роли
func HasRole(c *gin.Context, allowedRoles ...string) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}

	for _, allowed := range allowedRoles {
		if user.Role == allowed {
			return true
		}
	}
	return false
}

// GetRequestID получает ID запроса (если есть)
func GetRequestID(c *gin.Context) string {
	val, exists := c.Get("request_id")
	if !exists {
		return ""
	}
	requestID, ok := val.(string)
	if !ok {
		return ""
	}
	return requestID
}
