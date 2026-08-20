//	@title			Enterprise Resource Planning
//	@version		1.0
//	@description	For managing the enterprise's universal resources
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Shmonov Matvey
//	@contact.url	https://t.me/Koshsky
//	@contact.email	shmonov.mv@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				"Provide JWT token in the format: Bearer {token}"

//	@externalDocs.description	ERP documentation (placeholder)
//	@externalDocs.url			https://swagger.io/resources/open-api/

// Package server wires the HTTP routes and assembles the application.
package server

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/ratelimit"
	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/internal/userctx"
)

// registerRoutes registers every module's routes: public ones on the api
// group behind a per-IP limiter, protected ones on a sibling group behind
// RequireAuth and a per-user limiter. Protected routes are a sibling group
// (not a child of api) so they do NOT inherit the public per-IP limiter: a
// heavy authenticated user behind a shared NAT must not drain a common IP
// bucket and block its neighbors. Instead every authenticated user gets its
// own token bucket keyed by user id.
func (a *App) registerRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	api.Use(ratelimit.FromConfig(a.cfg.RateLimit, a.logger))

	// Health: liveness-probe без авторизации. Фронтенд проверяет им
	// доступность бэкенда перед офлайн-синхронизацией (вместо статики).
	api.GET("/health", a.health)

	for _, m := range a.modules {
		m.RegisterPublicRoutes(api)
	}

	protected := router.Group("/api/v1")
	protected.Use(a.authMw.RequireAuth())
	protected.Use(ratelimit.FromConfigKeyed(a.cfg.UserRateLimit, a.userKey, a.logger))
	for _, m := range a.modules {
		m.RegisterProtectedRoutes(protected)
	}
}

// userKey selects the token bucket identity for an authenticated request: the
// JWT user id. It runs after RequireAuth, which guarantees the user context is
// present; on any unexpected miss it falls back to a constant key so a
// misbehaving request shares one bucket instead of allocating unbounded ones.
func (a *App) userKey(c *gin.Context) string {
	id, err := userctx.GetUserID(c)
	if err != nil {
		return "anonymous"
	}
	return strconv.FormatInt(id, 10)
}

// health отвечает 200 с состоянием статус-пробы «ok».
func (a *App) health(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}
