package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/user")
	{
		r.GET("/", h.ListUsers)
		r.GET("/:id", h.GetUser)
		r.POST("/", h.CreateUser)
		r.PUT("/:id", h.UpdateUser)
		r.DELETE("/:id", h.DeleteUser)
	}
}
