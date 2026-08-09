package app

import (
	authDelivery "github.com/Koshsky/erp-backend/internal/auth/delivery"
	rbacMW "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	planningDelivery "github.com/Koshsky/erp-backend/internal/planning/delivery"
	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	calendarDelivery "github.com/Koshsky/erp-backend/internal/timesheet/calendar/delivery"
	employeeDelivery "github.com/Koshsky/erp-backend/internal/timesheet/employee/delivery"
	employeeRepo "github.com/Koshsky/erp-backend/internal/timesheet/employee/repository"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/timesheet/resource/delivery"
	resourceRepo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"
	stateDelivery "github.com/Koshsky/erp-backend/internal/timesheet/state/delivery"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
)

// ProvideAuthHandler exposes the auth handler as a route registrar.
func ProvideAuthHandler(h *authDelivery.AuthHandler) RouteRegistrar {
	return h
}

// ProvideRBACData assembles the owner resolvers for the policy engine.
func ProvideRBACData(
	project *projectRepo.ProjectRepository,
	process *processRepo.ProcessRepository,
	task *taskRepo.TaskRepository,
	milestone *milestoneRepo.MilestoneRepository,
	assignment *assignmentRepo.AssignmentRepository,
	resource *resourceRepo.ResourceRepository,
	employee *employeeRepo.EmployeeRepository,
) rbacMW.Data {
	return rbacMW.Data{
		ProjectOwners:    project.OwnerChain,
		ProcessOwners:    process.OwnerChain,
		TaskOwners:       task.OwnerChain,
		MilestoneOwners:  milestone.OwnerChain,
		AssignmentOwners: assignment.OwnerChain,
		ResourceOwners:   resource.OwnerChain,
		EmployeeOwners:   employee.OwnerChain,
	}
}

// ProtectedHandlers collects the handlers registered behind RequireAuth.
type ProtectedHandlers []RouteRegistrar

// ProvideProtectedHandlers collects the handlers registered behind RequireAuth.
func ProvideProtectedHandlers(
	planning *planningDelivery.PlanningHandler,
	user *userDelivery.UserHandler,
	task *taskDelivery.TaskHandler,
	project *projectDelivery.ProjectHandler,
	process *processDelivery.ProcessHandler,
	milestone *milestoneDelivery.MilestoneHandler,
	assignment *assignmentDelivery.AssignmentHandler,
	resource *resourceDelivery.ResourceHandler,
	state *stateDelivery.StateHandler,
	employee *employeeDelivery.EmployeeHandler,
	calendar *calendarDelivery.CalendarHandler,
) ProtectedHandlers {
	return ProtectedHandlers{planning, user, task, project, process, milestone, assignment, resource, state, employee, calendar}
}
