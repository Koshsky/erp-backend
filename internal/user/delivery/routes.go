package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/user")
	{
		r.GET("", h.ListUsers)
		r.GET("/:id", h.FindUser)
		r.GET("/:id/subordinates", h.ListSubordinates)
		r.PUT("/:id", h.UpdateUser)
		r.POST("/change-password", h.ChangePassword)
		r.DELETE("/:id", h.DeleteUser)
	}
}
