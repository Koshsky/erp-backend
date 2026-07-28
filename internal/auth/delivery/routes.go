package delivery

import "github.com/gin-gonic/gin"

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.POST("/change-password", h.ChangePassword)
		auth.POST("/refresh", h.RefreshToken)
	}
}
