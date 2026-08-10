package delivery

import "github.com/gin-gonic/gin"

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup, loginGuard gin.HandlerFunc) {
	r := router.Group("/auth")
	{
		r.POST("/register", h.Register)
		r.POST("/login", loginGuard, h.Login)
		r.POST("/refresh", h.RefreshToken)
	}
}
