//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"

	authDelivery "github.com/Koshsky/erp-backend/internal/auth/delivery"
	authService "github.com/Koshsky/erp-backend/internal/auth/service"
	"github.com/Koshsky/erp-backend/internal/config"
	"github.com/Koshsky/erp-backend/internal/database"
	"github.com/Koshsky/erp-backend/internal/logger"
	authMw "github.com/Koshsky/erp-backend/internal/middleware/auth"
	rbacMW "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	planningDelivery "github.com/Koshsky/erp-backend/internal/planning/delivery"
	planningRepo "github.com/Koshsky/erp-backend/internal/planning/repository"
	planningService "github.com/Koshsky/erp-backend/internal/planning/service"
	"github.com/Koshsky/erp-backend/internal/policies"
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
	"github.com/Koshsky/erp-backend/internal/security/jwt"
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

// InitializeApp builds the whole application dependency graph.
func InitializeApp() (*App, error) {
	wire.Build(
		config.ProvideConfig,
		logger.ProvideLogger,
		database.ProvidePostgresDB,
		config.ProvidePostgresConfig,
		config.ProvideJWTConfig,
		jwt.ProvideJWTService,
		authMw.ProvideAuthMiddleware,
		rbacMW.ProvideMiddleware,
		policies.ProvideAll,
		ProvideRBACData,

		userRepo.ProvideUserRepository,
		planningRepo.ProvidePlanningRepository,
		projectRepo.ProvideProjectRepository,
		processRepo.ProvideProcessRepository,
		taskRepo.ProvideTaskRepository,
		milestoneRepo.ProvideMilestoneRepository,
		assignmentRepo.ProvideAssignmentRepository,
		resourceRepo.ProvideResourceRepository,
		employeeRepo.ProvideEmployeeRepository,
		stateRepo.ProvideStateRepository,
		calendarRepo.ProvideCalendarRepository,

		userService.ProvideUserService,
		planningService.ProvidePlanningService,
		authService.ProvideAuthService,
		projectService.ProvideProjectService,
		processService.ProvideProcessService,
		taskService.ProvideTaskService,
		milestoneService.ProvideMilestoneService,
		assignmentService.ProvideAssignmentService,
		resourceService.ProvideResourceService,
		employeeService.ProvideEmployeeService,
		stateService.ProvideStateService,
		calendarService.ProvideCalendarService,

		userDelivery.ProvideUserHandler,
		planningDelivery.ProvidePlanningHandler,
		authDelivery.ProvideAuthHandler,
		projectDelivery.ProvideProjectHandler,
		processDelivery.ProvideProcessHandler,
		taskDelivery.ProvideTaskHandler,
		milestoneDelivery.ProvideMilestoneHandler,
		assignmentDelivery.ProvideAssignmentHandler,
		resourceDelivery.ProvideResourceHandler,
		employeeDelivery.ProvideEmployeeHandler,
		stateDelivery.ProvideStateHandler,
		calendarDelivery.ProvideCalendarHandler,

		ProvideAuthHandler,
		ProvideProtectedHandlers,
		New,
	)
	return nil, nil
}
