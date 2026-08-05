package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *AssignmentHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/assignment")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("", h.ListAssignments)
		r.GET("/:id", h.FindAssignment)
		r.POST("", h.CreateAssignment)
		r.PUT("/:id", h.UpdateAssignment)
		r.DELETE("/:id", h.DeleteAssignment)
	}
}
