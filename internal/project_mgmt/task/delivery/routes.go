package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *TaskHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/task")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("", h.ListTasks)
		r.GET("/:id", h.FindTask)
		r.POST("", h.CreateTask)
		r.PUT("/:id", h.UpdateTask)
		r.DELETE("/:id", h.DeleteTask)
	}
}
