package delivery

import "github.com/gin-gonic/gin"

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/project")
	{
		r.GET("", h.ListProjects)
		r.GET("/:id", h.FindProject)
		r.POST("", h.CreateProject)
		r.PUT("/:id", h.UpdateProject)
		r.DELETE("/:id", h.DeleteProject)
	}
}
