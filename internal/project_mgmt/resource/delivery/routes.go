package delivery

import "github.com/gin-gonic/gin"

func (h *ResourceHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/resource", h.ListResources)
	router.GET("/resource/:id", h.GetResource)
	router.POST("/resource", h.CreateResource)
	router.PUT("/resource/:id", h.UpdateResource)
	router.DELETE("/resource/:id", h.DeleteResource)
}
