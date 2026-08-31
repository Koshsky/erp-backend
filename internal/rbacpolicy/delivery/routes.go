package delivery

import "github.com/gin-gonic/gin"

// RegisterRoutes регистрирует admin-контур /rbac (весь за проверкой
// rbac.manage — виртуальный ресурс rbac_config, только admin).
func (h *RBACHandler) RegisterRoutes(router *gin.RouterGroup) {
	r := router.Group("/rbac")
	r.Use(h.mw.Check("rbac.manage"))
	{
		r.GET("/roles", h.ListRoles)
		r.POST("/roles", h.CreateRole)
		r.PUT("/roles/:name", h.UpdateRole)
		r.DELETE("/roles/:name", h.DeleteRole)
		r.GET("/rules", h.ListRules)
		r.PUT("/rules", h.UpsertRule)
		r.DELETE("/rules/:id", h.DeleteRule)
		r.GET("/policies", h.ListRoutePolicies)
		r.PUT("/policies", h.UpsertRoutePolicy)
		r.DELETE("/policies/:name", h.DeleteRoutePolicy)
		r.GET("/kinds", h.Kinds)
		r.GET("/matrix", h.Matrix)
		r.POST("/reset", h.Reset)
		r.GET("/explain", h.Explain)
	}
}
