package delivery

import "github.com/gin-gonic/gin"

func (h *SchedulingHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("scheduling")
	{
		r.GET("/projects", h.GetProjectScheduling)
		r.GET("/processes", h.GetProcessScheduling)
		r.GET("/tasks", h.GetTaskScheduling)
	}
}
