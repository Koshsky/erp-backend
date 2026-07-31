package delivery

import "github.com/gin-gonic/gin"

func (h *ResourceHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/resource")
	{
		r.GET("", h.ListResources)
		r.GET("/:id", h.FindResource)
		r.POST("", h.CreateResource)
		r.PUT("/:id", h.UpdateResource)
		r.DELETE("/:id", h.DeleteResource)
	}
}
