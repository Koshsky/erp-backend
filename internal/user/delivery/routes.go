package delivery

import "github.com/gin-gonic/gin"

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Аккаунт/самообслуживание и пул пользователей для пикеров владельцев.
	// GET /user — только для ролей с планировочными страницами (user.picker);
	// карточка по id — через защищённый /users/:id (worker.view).
	r := router.Group("/user")
	{
		r.GET("", h.mw.Check("user.picker"), h.ListAllUsers)
		r.POST("/change-password", h.ChangePassword)
	}

	// Управление пользователями/рабочими (коллекция).
	u := router.Group("/users")
	{
		u.GET("", h.mw.Check("worker.list"), h.ListUsers)
		u.POST("", h.mw.Check("worker.create"), h.CreateUser)
		u.GET("/:id", h.mw.Check("worker.view"), h.FindUser)
		u.PUT("/:id", h.mw.Check("worker.update"), h.UpdateUser)
		u.PUT("/:id/manager", h.mw.Check("worker.update"), h.UpdateManager)
		u.DELETE("/:id", h.mw.Check("worker.delete"), h.DeleteUser)
		u.POST("/:id/reset-password", h.mw.Check("worker.update"), h.ResetPassword)
		u.GET("/:id/days", h.mw.Check("worker.view"), h.ListDays)
		u.PUT("/:id/days", h.mw.Check("worker.update"), h.SetDays)
		u.DELETE("/:id/days", h.mw.Check("worker.delete"), h.DeleteDays)
	}
}
