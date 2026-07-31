package helpers

import (
	"errors"
	"slices"

	"github.com/Koshsky/erp-backend/internal/common/ctx"

	"github.com/gin-gonic/gin"
)

// GetUser extracts the user context from the Gin context.
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

// MustGetUser panics if the user is not found in the context.
func MustGetUser(c *gin.Context) ctx.UserContext {
	user, err := GetUser(c)
	if err != nil {
		panic(err)
	}
	return user
}

// GetUserID returns the user ID from the context.
func GetUserID(c *gin.Context) (int64, error) {
	user, err := GetUser(c)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetUserRole returns the user role from the context.
func GetUserRole(c *gin.Context) (string, error) {
	user, err := GetUser(c)
	if err != nil {
		return "", err
	}
	return user.Role, nil
}

// IsAdmin verifies if the user is an admin.
func IsAdmin(c *gin.Context) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}
	return user.Role == "admin"
}

// HasRole verifies if the user has any of the allowed roles.
func HasRole(c *gin.Context, allowedRoles ...string) bool {
	user, err := GetUser(c)
	if err != nil {
		return false
	}

	return slices.Contains(allowedRoles, user.Role)
}

// GetRequestID returns the request ID from the context.
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
