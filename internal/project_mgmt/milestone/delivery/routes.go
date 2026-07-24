package delivery

import "github.com/gin-gonic/gin"

func (h *MilestoneHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/milestone", h.ListMilestones)
	router.GET("/milestone/:id", h.GetMilestone)
	router.POST("/milestone", h.CreateMilestone)
	router.PUT("/milestone/:id", h.UpdateMilestone)
	router.DELETE("/milestone/:id", h.DeleteMilestone)
}
