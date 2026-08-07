package delivery

import (
	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/access"
)

func (h *EmployeeHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("")
	r.Use(access.DirectorReadOnly())
	{
		r.GET("/resources/:id/employees", h.ListEmployeesByResource)
		r.POST("/resources/:id/employees", h.CreateEmployee)
	}
	e := router.Group("/employees")
	e.Use(access.DirectorReadOnly())
	{
		e.GET("", h.ListEmployees)
		e.GET("/:id", h.FindEmployee)
		e.PUT("/:id", h.UpdateEmployee)
		e.DELETE("/:id", h.DeleteEmployee)
		e.GET("/:id/days", h.ListDays)
		e.PUT("/:id/days", h.SetDays)
		e.DELETE("/:id/days", h.DeleteDays)
	}
}
