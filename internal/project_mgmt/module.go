// Package project_mgmt wires the project management module's providers and routes.
package project_mgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
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
)

// ProviderSet aggregates the project management module's dependencies.
var ProviderSet = wire.NewSet(
	projectRepo.NewProjectRepository,
	projectService.NewProjectService,
	projectDelivery.NewProjectHandler,

	processRepo.NewProcessRepository,
	processService.NewProcessService,
	processDelivery.NewProcessHandler,

	taskRepo.NewTaskRepository,
	taskService.NewTaskService,
	taskDelivery.NewTaskHandler,

	milestoneRepo.NewMilestoneRepository,
	milestoneService.NewMilestoneService,
	milestoneDelivery.NewMilestoneHandler,

	assignmentRepo.NewAssignmentRepository,
	service.NewAssignmentService,
	delivery.NewAssignmentHandler,

	ProvideModule,
)

// Module registers the project management module's routes (all protected).
type Module struct {
	task       *taskDelivery.TaskHandler
	project    *projectDelivery.ProjectHandler
	process    *processDelivery.ProcessHandler
	milestone  *milestoneDelivery.MilestoneHandler
	assignment *delivery.AssignmentHandler
}

// ProvideModule builds the project management module.
func ProvideModule(
	task *taskDelivery.TaskHandler,
	project *projectDelivery.ProjectHandler,
	process *processDelivery.ProcessHandler,
	milestone *milestoneDelivery.MilestoneHandler,
	assignment *delivery.AssignmentHandler,
) Module {
	return Module{
		task:       task,
		project:    project,
		process:    process,
		milestone:  milestone,
		assignment: assignment,
	}
}

// RegisterPublicRoutes is a no-op: the module has no public routes.
func (m Module) RegisterPublicRoutes(r *gin.RouterGroup) {
}

// RegisterProtectedRoutes registers the module's routes behind authentication.
func (m Module) RegisterProtectedRoutes(r *gin.RouterGroup) {
	m.task.RegisterRoutes(r)
	m.project.RegisterRoutes(r)
	m.process.RegisterRoutes(r)
	m.milestone.RegisterRoutes(r)
	m.assignment.RegisterRoutes(r)
}
