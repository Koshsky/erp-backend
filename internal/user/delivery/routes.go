package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, changePasswordGuard gin.HandlerFunc) {
	// Single /user route set: account/self-service, picker pool and CRUD collection.
	// - GET /user       — scoped list with pagination (admin sees all, vp — own subordinates).
	// - GET /user/all   — unscoped pool for owner pickers (user.picker).
	// - /user/{id}      — user/worker card and management.
	// - /user/{id}/days — employee timesheet.
	r := router.Group("/user")
	{
		r.GET("/all", h.mw.Check("user.picker"), h.ListAllUsers)
		// CRUD collection.
		r.GET("", h.mw.Check("worker.list"), h.ListUsers)
		r.POST("", h.mw.Check("worker.create"), h.CreateUser)
		r.GET("/:id", h.mw.Check("worker.view"), h.FindUser)
		r.PUT("/:id", h.mw.Check("worker.update"), h.UpdateUser)
		r.PUT("/:id/manager", h.mw.Check("worker.update"), h.UpdateManager)
		r.DELETE("/:id", h.mw.Check("worker.delete"), h.DeleteUser)
		r.POST("/:id/reset-password", h.mw.Check("worker.update"), h.ResetPassword)
		// Employee timesheet.
		r.GET("/:id/days", h.mw.Check("worker.view"), h.ListDays)
		r.PUT("/:id/days", h.mw.Check("worker.update"), h.SetDays)
		r.DELETE("/:id/days", h.mw.Check("worker.delete"), h.DeleteDays)
		// Self-service (password change) — with a per-user limit and uniform delay.
		r.POST("/change-password", changePasswordGuard, h.ChangePassword)
	}
}
