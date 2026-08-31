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
		// CRUD collection. An employee IS a system user (the users table), so
		// profile mutations are gated by the user-admin rights (user_admin.*,
		// admin by default; grantable via the RBAC matrix), not by worker.*:
		// editing employees happens only through the admin users section.
		r.GET("", h.mw.Check("worker.list"), h.ListUsers)
		r.POST("", h.mw.Check("user_admin.create"), h.CreateUser)
		r.GET("/:id", h.mw.Check("worker.view"), h.FindUser)
		r.PUT("/:id", h.mw.Check("user_admin.update"), h.UpdateUser)
		r.PUT("/:id/manager", h.mw.Check("user_admin.update"), h.UpdateManager)
		r.DELETE("/:id", h.mw.Check("user_admin.delete"), h.DeleteUser)
		r.POST("/:id/reset-password", h.mw.Check("user_admin.update"), h.ResetPassword)
		// Employee timesheet.
		r.GET("/:id/days", h.mw.Check("worker.view"), h.ListDays)
		r.PUT("/:id/days", h.mw.Check("worker.update"), h.SetDays)
		r.DELETE("/:id/days", h.mw.Check("worker.delete"), h.DeleteDays)
		// Self-service (password change) — with a per-user limit and uniform delay.
		r.POST("/change-password", changePasswordGuard, h.ChangePassword)
	}
}
