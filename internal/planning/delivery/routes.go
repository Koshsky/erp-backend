package delivery

import "github.com/gin-gonic/gin"

func (h *PlanningHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("planning")
	{
		r.GET("/projects", h.GetProjectPlanning)
		r.GET("/processes", h.GetProcessPlanning)
		r.GET("/tasks", h.GetTaskPlanning)
	}
}
