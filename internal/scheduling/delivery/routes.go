package delivery

import "github.com/gin-gonic/gin"

func (h *SchedulingHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/scheduling/projects", h.GetProjectScheduling)
	router.GET("/scheduling/processes", h.GetProcessScheduling)
	router.GET("/scheduling/tasks", h.GetTaskScheduling)
}
