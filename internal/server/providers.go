package server

import (
	"github.com/Koshsky/erp-backend/internal/auth"
	rbacMW "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/planning"
	"github.com/Koshsky/erp-backend/internal/project_mgmt"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	milestoneRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	"github.com/Koshsky/erp-backend/internal/timesheet"
	resourceRepo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"
	"github.com/Koshsky/erp-backend/internal/user"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
)

// ProvideRBACData assembles the owner resolvers for the policy engine. The
// repositories satisfy rbac.OwnerResolver (asserted in each repository).
func ProvideRBACData(
	project *projectRepo.ProjectRepository,
	process *processRepo.ProcessRepository,
	task *taskRepo.TaskRepository,
	milestone *milestoneRepo.MilestoneRepository,
	assignment *assignmentRepo.AssignmentRepository,
	resource *resourceRepo.ResourceRepository,
	user *userRepo.UserRepository,
) rbacMW.Data {
	return rbacMW.Data{
		ProjectOwners:    project.OwnerChain,
		ProcessOwners:    process.OwnerChain,
		TaskOwners:       task.OwnerChain,
		MilestoneOwners:  milestone.OwnerChain,
		AssignmentOwners: assignment.OwnerChain,
		ResourceOwners:   resource.OwnerChain,
		WorkerOwners:     user.OwnerChain,
	}
}

// ProvideModules collects the domain modules of the application.
func ProvideModules(
	auth auth.Module,
	user user.Module,
	planning planning.Module,
	project project_mgmt.Module,
	timesheet timesheet.Module,
) []Module {
	return []Module{auth, user, planning, project, timesheet}
}
