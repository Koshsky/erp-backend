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
	authService "github.com/Koshsky/erp-backend/internal/auth/service"

	rbacMW "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/middleware/rbac/policies"
	planningDelivery "github.com/Koshsky/erp-backend/internal/planning/delivery"
	planningRepo "github.com/Koshsky/erp-backend/internal/planning/repository"
	planningService "github.com/Koshsky/erp-backend/internal/planning/service"
	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	assignmentService "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	milestoneService "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/service"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	processService "github.com/Koshsky/erp-backend/internal/project_mgmt/process/service"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	projectService "github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	taskService "github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
	calendarDelivery "github.com/Koshsky/erp-backend/internal/timesheet/calendar/delivery"
	calendarRepo "github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository"
	calendarService "github.com/Koshsky/erp-backend/internal/timesheet/calendar/service"
	employeeDelivery "github.com/Koshsky/erp-backend/internal/timesheet/employee/delivery"
	employeeRepo "github.com/Koshsky/erp-backend/internal/timesheet/employee/repository"
	employeeService "github.com/Koshsky/erp-backend/internal/timesheet/employee/service"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/timesheet/resource/delivery"
	resourceRepo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"
	resourceService "github.com/Koshsky/erp-backend/internal/timesheet/resource/service"
	stateDelivery "github.com/Koshsky/erp-backend/internal/timesheet/state/delivery"
	stateRepo "github.com/Koshsky/erp-backend/internal/timesheet/state/repository"
	stateService "github.com/Koshsky/erp-backend/internal/timesheet/state/service"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
	userService "github.com/Koshsky/erp-backend/internal/user/service"
)

type RouteRegistrar interface {
	RegisterRoutes(router *gin.RouterGroup)
}

var (
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

func (a *App) registerRoutes(router *gin.Engine) {
	// --- Shared dependencies ---
	userRepo := userRepo.NewUserRepository(a.logger, a.pool)

	// --- User ---
	userSvc := userService.NewUserService(a.logger, userRepo)
	userHandler := userDelivery.NewUserHandler(a.logger, userSvc)

	// --- Auth (uses userSvc for user data) ---
	authSvc := authService.NewAuthService(userSvc, a.jwtManager)
	authHandler := authDelivery.NewAuthHandler(a.logger, authSvc)

	// --- Planning (Gantt: projects/processes/tasks/milestones/assignments) ---
	planningQueries := planningRepo.NewPlanningRepository(a.logger, a.pool)
	planningSvc := planningService.NewPlanningService(a.logger, planningQueries)
	planningHandler := planningDelivery.NewPlanningHandler(a.logger, planningSvc)

	// --- Repositories (services + RBAC resolvers) ---
	projectQueries := projectRepo.NewProjectRepository(a.logger, a.pool)
	processQueries := processRepo.NewProcessRepository(a.logger, a.pool)
	taskQueries := taskRepo.NewTaskRepository(a.logger, a.pool)
	milestoneQueries := milestoneRepo.NewMilestoneRepository(a.logger, a.pool)
	assignmentQueries := assignmentRepo.NewAssignmentRepository(a.logger, a.pool)
	resourceQueries := resourceRepo.NewResourceRepository(a.logger, a.pool)
	employeeQueries := employeeRepo.NewEmployeeRepository(a.logger, a.pool)

	// --- RBAC: policy engine (Data — owner resolvers from repositories) ---
	rbacMiddleware := rbacMW.New(a.logger, rbacMW.Data{
		ProjectOwners:    projectQueries.OwnerChain,
		ProcessOwners:    processQueries.OwnerChain,
		TaskOwners:       taskQueries.OwnerChain,
		MilestoneOwners:  milestoneQueries.OwnerChain,
		AssignmentOwners: assignmentQueries.OwnerChain,
		ResourceOwners:   resourceQueries.OwnerChain,
		EmployeeOwners:   employeeQueries.OwnerChain,
	}, policies.All())

	// --- Task ---
	taskSvc := taskService.NewTaskService(a.logger, taskQueries)
	taskHandler := taskDelivery.NewTaskHandler(a.logger, taskSvc, rbacMiddleware)

	// --- Project ---
	projectSvc := projectService.NewProjectService(a.logger, projectQueries)
	projectHandler := projectDelivery.NewProjectHandler(a.logger, projectSvc, rbacMiddleware)

	// --- Process ---
	processSvc := processService.NewProcessService(a.logger, processQueries)
	processHandler := processDelivery.NewProcessHandler(a.logger, processSvc, rbacMiddleware)

	// --- Milestone ---
	milestoneSvc := milestoneService.NewMilestoneService(a.logger, milestoneQueries)
	milestoneHandler := milestoneDelivery.NewMilestoneHandler(a.logger, milestoneSvc, rbacMiddleware)

	// --- Assignment ---
	assignmentSvc := assignmentService.NewAssignmentService(a.logger, assignmentQueries)
	assignmentHandler := assignmentDelivery.NewAssignmentHandler(a.logger, assignmentSvc, rbacMiddleware)

	api := router.Group("/api/v1")

	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(a.authMw.RequireAuth())
	{
		planningHandler.RegisterRoutes(protected)
		userHandler.RegisterRoutes(protected)
		taskHandler.RegisterRoutes(protected)
		projectHandler.RegisterRoutes(protected)
		processHandler.RegisterRoutes(protected)
		milestoneHandler.RegisterRoutes(protected)
		assignmentHandler.RegisterRoutes(protected)

		a.registerTimesheet(protected.Group("/timesheet"), rbacMiddleware, resourceQueries, employeeQueries)
	}
}

// registerTimesheet wires the timesheet domains (resources, states, employees,
// calendar) and registers their routes under a common prefix.
func (a *App) registerTimesheet(
	router *gin.RouterGroup,
	rbacMiddleware *rbacMW.Middleware,
	resourceQueries *resourceRepo.ResourceRepository,
	employeeQueries *employeeRepo.EmployeeRepository,
) {
	resourceSvc := resourceService.NewResourceService(a.logger, resourceQueries)
	resourceHandler := resourceDelivery.NewResourceHandler(a.logger, resourceSvc, rbacMiddleware)

	stateQueries := stateRepo.NewStateRepository(a.logger, a.pool)
	stateSvc := stateService.NewStateService(a.logger, stateQueries)
	stateHandler := stateDelivery.NewStateHandler(a.logger, stateSvc, rbacMiddleware)

	employeeSvc := employeeService.NewEmployeeService(a.logger, employeeQueries)
	employeeHandler := employeeDelivery.NewEmployeeHandler(a.logger, employeeSvc, rbacMiddleware)

	calendarQueries := calendarRepo.NewCalendarRepository(a.logger, a.pool)
	calendarSvc := calendarService.NewCalendarService(a.logger, calendarQueries)
	calendarHandler := calendarDelivery.NewCalendarHandler(a.logger, calendarSvc, rbacMiddleware)

	resourceHandler.RegisterRoutes(router)
	stateHandler.RegisterRoutes(router)
	employeeHandler.RegisterRoutes(router)
	calendarHandler.RegisterRoutes(router)
}
