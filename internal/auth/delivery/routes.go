package delivery

import "github.com/gin-gonic/gin"

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/auth")
	{
		r.POST("/register", h.Register)
		r.POST("/login", h.Login)
		r.POST("/logout", h.Logout)
		r.POST("/change-password", h.ChangePassword)
		r.POST("/refresh", h.RefreshToken)
	}
}
