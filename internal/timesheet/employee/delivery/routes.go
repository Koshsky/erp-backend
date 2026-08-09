package delivery

import (
	"github.com/gin-gonic/gin"
)

func (h *EmployeeHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/timesheet")
	{
		r.GET("/resources/:id/employees", h.ListEmployeesByResource)
		r.POST("/resources/:id/employees", h.mw.Check("employee.create"), h.CreateEmployee)
	}
	e := router.Group("/timesheet/employees")
	{
		e.GET("", h.ListEmployees)
		e.GET("/:id", h.mw.Check("employee.view"), h.FindEmployee)
		e.PUT("/:id", h.mw.Check("employee.update"), h.UpdateEmployee)
		e.DELETE("/:id", h.mw.Check("employee.delete"), h.DeleteEmployee)
		e.GET("/:id/days", h.mw.Check("employee.view"), h.ListDays)
		e.PUT("/:id/days", h.mw.Check("employee.update"), h.SetDays)
		e.DELETE("/:id/days", h.mw.Check("employee.delete"), h.DeleteDays)
	}
}
