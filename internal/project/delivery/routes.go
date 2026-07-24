package delivery

import "github.com/gin-gonic/gin"

func (h *ProjectHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/project", h.ListProjects)
	router.GET("/project/:id", h.GetProject)
	router.POST("/project", h.CreateProject)
	router.PUT("/project/:id", h.UpdateProject)
	router.DELETE("/project/:id", h.DeleteProject)
}

