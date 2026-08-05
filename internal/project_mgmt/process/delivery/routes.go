package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *ProcessHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/process")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("", h.ListProcesses)
		r.GET("/:id", h.FindProcess)
		r.POST("", h.CreateProcess)
		r.PUT("/:id", h.UpdateProcess)
		r.DELETE("/:id", h.DeleteProcess)
	}
}
