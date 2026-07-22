package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (a *App) runSwaggerServer(router *gin.Engine) {
	router.GET("/api/openapi.yaml", func(c *gin.Context) {
		specPath := "docs/openapi.yaml"

		data, err := os.ReadFile(specPath)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Data(http.StatusOK, "application/x-yaml", data)
	})

	// Swagger UI
	router.GET("/api/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/api/openapi.yaml"),
		httpSwagger.DocExpansion("none"),
	)))
}
