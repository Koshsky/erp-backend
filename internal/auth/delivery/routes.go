package delivery

import "github.com/gin-gonic/gin"

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup, loginGuard gin.HandlerFunc) {
	r := router.Group("/auth")
	{
		r.POST("/login", loginGuard, h.Login)
		r.POST("/refresh", h.RefreshToken)
		r.POST("/logout", h.Logout)
	}
}
