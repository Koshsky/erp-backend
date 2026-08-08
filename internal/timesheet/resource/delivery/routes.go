package delivery

import (
	"github.com/gin-gonic/gin"
)

func (h *ResourceHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/resources")
	{
		r.GET("", h.ListResources)
		r.GET("/:id", h.mw.Check("resource.view"), h.FindResource)
		r.POST("", h.mw.Check("resource.create"), h.CreateResource)
		r.PUT("/:id", h.mw.Check("resource.update"), h.UpdateResource)
		r.DELETE("/:id", h.mw.Check("resource.delete"), h.DeleteResource)
	}
}
