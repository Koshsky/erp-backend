package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *StateHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/states")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("", h.ListStates)
		r.GET("/:id", h.FindState)
		r.POST("", h.CreateState)
		r.PUT("/:id", h.UpdateState)
		r.DELETE("/:id", h.DeleteState)
	}
}
