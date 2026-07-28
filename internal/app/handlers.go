package app

import (
	"github.com/gin-gonic/gin"

	authDelivery "github.com/Koshsky/erp-backend/internal/auth/delivery"
	authDomain "github.com/Koshsky/erp-backend/internal/auth/domain"
	authRepo "github.com/Koshsky/erp-backend/internal/auth/repository"

	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	milestoneRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/delivery"
	resourceDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	resourceRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/repository"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	schedulingDelivery "github.com/Koshsky/erp-backend/internal/scheduling/delivery"
	schedulingDomain "github.com/Koshsky/erp-backend/internal/scheduling/domain"
	schedulingRepo "github.com/Koshsky/erp-backend/internal/scheduling/repository"
	"github.com/Koshsky/erp-backend/internal/security/password"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
	userDomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
)

type RouteRegistrar interface {
	RegisterRoutes(router *gin.RouterGroup)
}

var (
	_ RouteRegistrar = (*schedulingDelivery.SchedulingHandler)(nil)
	_ RouteRegistrar = (*userDelivery.UserHandler)(nil)
	_ RouteRegistrar = (*taskDelivery.TaskHandler)(nil)
	_ RouteRegistrar = (*resourceDelivery.ResourceHandler)(nil)
	_ RouteRegistrar = (*projectDelivery.ProjectHandler)(nil)
	_ RouteRegistrar = (*processDelivery.ProcessHandler)(nil)
	_ RouteRegistrar = (*milestoneDelivery.MilestoneHandler)(nil)
	_ RouteRegistrar = (*assignmentDelivery.AssignmentHandler)(nil)
)

func (a *App) registerRoutes(router *gin.Engine) {
	// --- Auth ---
	authQueries := authRepo.NewAuthRepository(a.pool)
	authHasher := password.NewBcryptHasher()
	authSvc := authDomain.NewAuthService(authQueries, authHasher, a.jwtManager)
	authHandler := authDelivery.NewAuthHandler(authSvc)

	// --- Scheduling ---
	schedulingQueries := schedulingRepo.NewSchedulingRepository(a.logger, a.pool)
	schedulingSvc := schedulingDomain.NewSchedulingService(a.logger, schedulingQueries)
	schedulingHandler := schedulingDelivery.NewSchedulingHandler(a.logger, schedulingSvc)

	// --- User ---
	userQueries := userRepo.NewUserRepository(a.logger, a.pool)
	userHasher := password.NewBcryptHasher()
	userSvc := userDomain.NewUserService(a.logger, userQueries, userHasher)
	userHandler := userDelivery.NewUserHandler(a.logger, userSvc)

	// --- Task ---
	taskQueries := taskRepo.NewTaskRepository(a.logger, a.pool)
	taskSvc := taskDomain.NewTaskService(a.logger, taskQueries)
	taskHandler := taskDelivery.NewTaskHandler(a.logger, taskSvc)

	// --- Resource ---
	resourceQueries := resourceRepo.NewResourceRepository(a.logger, a.pool)
	resourceSvc := resourceDomain.NewResourceService(a.logger, resourceQueries)
	resourceHandler := resourceDelivery.NewResourceHandler(a.logger, resourceSvc)

	// --- Project ---
	projectQueries := projectRepo.NewProjectRepository(a.logger, a.pool)
	projectSvc := projectDomain.NewProjectService(a.logger, projectQueries)
	projectHandler := projectDelivery.NewProjectHandler(a.logger, projectSvc)

	// --- Process ---
	processQueries := processRepo.NewProcessRepository(a.logger, a.pool)
	processSvc := processDomain.NewProcessService(a.logger, processQueries)
	processHandler := processDelivery.NewProcessHandler(a.logger, processSvc)

	// --- Milestone ---
	milestoneQueries := milestoneRepo.NewMilestoneRepository(a.logger, a.pool)
	milestoneSvc := milestoneDomain.NewMilestoneService(a.logger, milestoneQueries)
	milestoneHandler := milestoneDelivery.NewMilestoneHandler(a.logger, milestoneSvc)

	// --- Assignment ---
	assignmentQueries := assignmentRepo.NewAssignmentRepository(a.logger, a.pool)
	assignmentSvc := assignmentDomain.NewAssignmentService(a.logger, assignmentQueries)
	assignmentHandler := assignmentDelivery.NewAssignmentHandler(a.logger, assignmentSvc)

	api := router.Group("/api/v1")

	authHandler.RegisterRoutes(api)

	protected := api.Group("")
	protected.Use(a.authMw.RequireAuth())
	{
		schedulingHandler.RegisterRoutes(protected)
		userHandler.RegisterRoutes(protected)
		taskHandler.RegisterRoutes(protected)
		resourceHandler.RegisterRoutes(protected)
		projectHandler.RegisterRoutes(protected)
		processHandler.RegisterRoutes(protected)
		milestoneHandler.RegisterRoutes(protected)
		assignmentHandler.RegisterRoutes(protected)
	}
}
