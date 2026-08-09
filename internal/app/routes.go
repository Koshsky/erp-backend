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
package app

import (
	"github.com/gin-gonic/gin"

	authDelivery "github.com/Koshsky/erp-backend/internal/auth/delivery"
	planningDelivery "github.com/Koshsky/erp-backend/internal/planning/delivery"
	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	calendarDelivery "github.com/Koshsky/erp-backend/internal/timesheet/calendar/delivery"
	employeeDelivery "github.com/Koshsky/erp-backend/internal/timesheet/employee/delivery"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/timesheet/resource/delivery"
	stateDelivery "github.com/Koshsky/erp-backend/internal/timesheet/state/delivery"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
)

type RouteRegistrar interface {
	RegisterRoutes(router *gin.RouterGroup)
}

var (
	_ RouteRegistrar = (*authDelivery.AuthHandler)(nil)
	_ RouteRegistrar = (*planningDelivery.PlanningHandler)(nil)
	_ RouteRegistrar = (*userDelivery.UserHandler)(nil)
	_ RouteRegistrar = (*taskDelivery.TaskHandler)(nil)
	_ RouteRegistrar = (*resourceDelivery.ResourceHandler)(nil)
	_ RouteRegistrar = (*projectDelivery.ProjectHandler)(nil)
	_ RouteRegistrar = (*processDelivery.ProcessHandler)(nil)
	_ RouteRegistrar = (*milestoneDelivery.MilestoneHandler)(nil)
	_ RouteRegistrar = (*assignmentDelivery.AssignmentHandler)(nil)
	_ RouteRegistrar = (*stateDelivery.StateHandler)(nil)
	_ RouteRegistrar = (*employeeDelivery.EmployeeHandler)(nil)
	_ RouteRegistrar = (*calendarDelivery.CalendarHandler)(nil)
)

// registerRoutes registers the injected handlers: auth on the public group,
// the rest behind RequireAuth.
func (a *App) registerRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")

	a.authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(a.authMw.RequireAuth())
	for _, h := range a.protectedHandlers {
		h.RegisterRoutes(protected)
	}
}
