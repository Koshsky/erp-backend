package app

import (
	"github.com/gin-gonic/gin"

	authDelivery "github.com/Koshsky/erp-backend/internal/auth/delivery"
	authService "github.com/Koshsky/erp-backend/internal/auth/service"

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
	resourceDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/delivery"
	resourceRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/repository"
	resourceService "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/service"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	taskService "github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
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

	// --- Planning ---
	planningQueries := planningRepo.NewPlanningRepository(a.logger, a.pool)
	planningSvc := planningService.NewPlanningService(a.logger, planningQueries)
	planningHandler := planningDelivery.NewPlanningHandler(a.logger, planningSvc)

	// --- Task ---
	taskQueries := taskRepo.NewTaskRepository(a.logger, a.pool)
	taskSvc := taskService.NewTaskService(a.logger, taskQueries)
	taskHandler := taskDelivery.NewTaskHandler(a.logger, taskSvc)

	// --- Resource ---
	resourceQueries := resourceRepo.NewResourceRepository(a.logger, a.pool)
	resourceSvc := resourceService.NewResourceService(a.logger, resourceQueries)
	resourceHandler := resourceDelivery.NewResourceHandler(a.logger, resourceSvc)

	// --- Project ---
	projectQueries := projectRepo.NewProjectRepository(a.logger, a.pool)
	projectSvc := projectService.NewProjectService(a.logger, projectQueries)
	projectHandler := projectDelivery.NewProjectHandler(a.logger, projectSvc)

	// --- Process ---
	processQueries := processRepo.NewProcessRepository(a.logger, a.pool)
	processSvc := processService.NewProcessService(a.logger, processQueries)
	processHandler := processDelivery.NewProcessHandler(a.logger, processSvc)

	// --- Milestone ---
	milestoneQueries := milestoneRepo.NewMilestoneRepository(a.logger, a.pool)
	milestoneSvc := milestoneService.NewMilestoneService(a.logger, milestoneQueries)
	milestoneHandler := milestoneDelivery.NewMilestoneHandler(a.logger, milestoneSvc)

	// --- Assignment ---
	assignmentQueries := assignmentRepo.NewAssignmentRepository(a.logger, a.pool)
	assignmentSvc := assignmentService.NewAssignmentService(a.logger, assignmentQueries)
	assignmentHandler := assignmentDelivery.NewAssignmentHandler(a.logger, assignmentSvc)

	api := router.Group("/api/v1")

	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(a.authMw.RequireAuth())
	{
		planningHandler.RegisterRoutes(protected)
		userHandler.RegisterRoutes(protected)
		taskHandler.RegisterRoutes(protected)
		resourceHandler.RegisterRoutes(protected)
		projectHandler.RegisterRoutes(protected)
		processHandler.RegisterRoutes(protected)
		milestoneHandler.RegisterRoutes(protected)
		assignmentHandler.RegisterRoutes(protected)
	}
}
