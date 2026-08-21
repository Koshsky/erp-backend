package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Единый контур /user: аккаунт/самообслуживание, пул пикеров и CRUD-коллекция.
	// - GET /user       — скоупированный список с пагинацией (admin всё, vp — свои подчинённые).
	// - GET /user/all   — нескоупированный пул для пикеров владельцев (user.picker).
	// - /user/{id}      — карточка и управление пользователем/рабочим.
	// - /user/{id}/days — табель сотрудника.
	r := router.Group("/user")
	{
		r.GET("/all", h.mw.Check("user.picker"), h.ListAllUsers)
		// CRUD коллекции.
		r.GET("", h.mw.Check("worker.list"), h.ListUsers)
		r.POST("", h.mw.Check("worker.create"), h.CreateUser)
		r.GET("/:id", h.mw.Check("worker.view"), h.FindUser)
		r.PUT("/:id", h.mw.Check("worker.update"), h.UpdateUser)
		r.PUT("/:id/manager", h.mw.Check("worker.update"), h.UpdateManager)
		r.DELETE("/:id", h.mw.Check("worker.delete"), h.DeleteUser)
		r.POST("/:id/reset-password", h.mw.Check("worker.update"), h.ResetPassword)
		// Табель сотрудника.
		r.GET("/:id/days", h.mw.Check("worker.view"), h.ListDays)
		r.PUT("/:id/days", h.mw.Check("worker.update"), h.SetDays)
		r.DELETE("/:id/days", h.mw.Check("worker.delete"), h.DeleteDays)
		// Самообслуживание.
		r.POST("/change-password", h.ChangePassword)
	}
}
