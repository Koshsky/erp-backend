// @title			Enterprise Resource Planning
// @version		1.0
// @description	For managing the enterprise's universal resources
// @termsOfService	http://swagger.io/terms/

// @contact.name	Shmonov Matvey
// @contact.url	https://t.me/Koshsky
// @contact.email	shmonov.mv@gmail.com

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				"Provide JWT token in the format: Bearer {token}"

// @externalDocs.description	ERP documentation (placeholder)
// @externalDocs.url			https://swagger.io/resources/open-api/

// Package app wires the HTTP routes.
package server

import (
	"github.com/gin-gonic/gin"
)

// registerRoutes registers every module's routes: public ones on the api
// group, protected ones behind RequireAuth.
func (a *App) registerRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	for _, m := range a.modules {
		m.RegisterPublicRoutes(api)
	}

	protected := api.Group("")
	protected.Use(a.authMw.RequireAuth())
	for _, m := range a.modules {
		m.RegisterProtectedRoutes(protected)
	}
}
