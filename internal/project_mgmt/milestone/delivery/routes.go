package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *MilestoneHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/milestone")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("", h.ListMilestones)
		r.GET("/:id", h.FindMilestone)
		r.POST("", h.CreateMilestone)
		r.PUT("/:id", h.UpdateMilestone)
		r.DELETE("/:id", h.DeleteMilestone)
	}
}
