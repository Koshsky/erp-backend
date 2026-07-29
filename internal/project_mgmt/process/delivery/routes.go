package delivery

import "github.com/gin-gonic/gin"

func (h *ProcessHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/process")
	{
		r.GET("/", h.ListProcesses)
		r.GET("/:id", h.FindProcess)
		r.POST("/", h.CreateProcess)
		r.PUT("/:id", h.UpdateProcess)
		r.DELETE("/:id", h.DeleteProcess)
	}
}
