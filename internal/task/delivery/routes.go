package delivery

import "github.com/gin-gonic/gin"

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/task", h.ListTasks)
	router.GET("/task/:id", h.GetTask)
	router.POST("/task", h.CreateTask)
	router.PUT("/task/:id", h.UpdateTask)
	router.DELETE("/task/:id", h.DeleteTask)
}
