package delivery

import "github.com/gin-gonic/gin"

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/task")
	{
		r.GET("", h.ListTasks)
		r.GET("/:id", h.FindTask)
		r.POST("", h.CreateTask)
		r.PUT("/:id", h.UpdateTask)
		r.DELETE("/:id", h.DeleteTask)
	}
}
