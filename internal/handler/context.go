package handler

import (
	"context"

	"github.com/gin-gonic/gin"
)

// withAuthContext возвращает context.Context, обогащённый значениями role и user_id из gin.Context.
// Это необходимо, так как c.Request.Context() не содержит значений, установленных через c.Set().
func withAuthContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if role, exists := c.Get("role"); exists {
		ctx = context.WithValue(ctx, "role", role)
	}
	if userID, exists := c.Get("user_id"); exists {
		ctx = context.WithValue(ctx, "user_id", userID)
	}
	return ctx
}
