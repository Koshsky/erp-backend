package delivery

import "github.com/gin-gonic/gin"

func (h *AssignmentHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/assignment")
	{
		r.GET("/", h.ListAssignments)
		r.GET("/:id", h.FindAssignment)
		r.POST("/", h.CreateAssignment)
		r.PUT("/:id", h.UpdateAssignment)
		r.DELETE("/:id", h.DeleteAssignment)
	}
}
