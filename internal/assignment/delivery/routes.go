package delivery

import "github.com/gin-gonic/gin"

func (h *AssignmentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/assignment", h.ListAssignments)
	router.GET("/assignment/:id", h.GetAssignment)
	router.POST("/assignment", h.CreateAssignment)
	router.PUT("/assignment/:id", h.UpdateAssignment)
	router.DELETE("/assignment/:id", h.DeleteAssignment)
}