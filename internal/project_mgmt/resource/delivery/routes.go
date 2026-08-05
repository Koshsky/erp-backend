package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *ResourceHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/resource")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("", h.ListResources)
		r.GET("/:id", h.FindResource)
		r.POST("", h.CreateResource)
		r.PUT("/:id", h.UpdateResource)
		r.DELETE("/:id", h.DeleteResource)
	}
}
