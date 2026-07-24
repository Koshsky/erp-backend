package delivery

import "github.com/gin-gonic/gin"

func (h *ProcessHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/process", h.ListProcesses)
	router.GET("/process/:id", h.GetProcess)
	router.POST("/process", h.CreateProcess)
	router.PUT("/process/:id", h.UpdateProcess)
	router.DELETE("/process/:id", h.DeleteProcess)
}

