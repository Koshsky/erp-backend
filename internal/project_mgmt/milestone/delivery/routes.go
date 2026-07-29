package delivery

import "github.com/gin-gonic/gin"

func (h *MilestoneHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/milestone")
	{
		r.GET("/", h.ListMilestones)
		r.GET("/:id", h.FindMilestone)
		r.POST("/", h.CreateMilestone)
		r.PUT("/:id", h.UpdateMilestone)
		r.DELETE("/:id", h.DeleteMilestone)
	}
}
