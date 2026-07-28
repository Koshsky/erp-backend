package app

import (
	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentService "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneService "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/service"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processService "github.com/Koshsky/erp-backend/internal/project_mgmt/process/service"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectService "github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/delivery"
	resourceService "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/service"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskService "github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
	userService "github.com/Koshsky/erp-backend/internal/user/service"
)

// TYPE ALIASES
type (
	UserService       = userDelivery.UserService
	TaskService       = taskDelivery.TaskService
	ResourceService   = resourceDelivery.ResourceService
	ProjectService    = projectDelivery.ProjectService
	ProcessService    = processDelivery.ProcessService
	MilestoneService  = milestoneDelivery.MilestoneService
	AssignmentService = assignmentDelivery.AssignmentService
)

// COMPILATION CHECK
var (
	_ UserService       = (*userService.UserService)(nil)
	_ TaskService       = (*taskService.TaskService)(nil)
	_ ResourceService   = (*resourceService.ResourceService)(nil)
	_ ProjectService    = (*projectService.ProjectService)(nil)
	_ ProcessService    = (*processService.ProcessService)(nil)
	_ MilestoneService  = (*milestoneService.MilestoneService)(nil)
	_ AssignmentService = (*assignmentService.AssignmentService)(nil)
)
