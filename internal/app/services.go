package app

import (
	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/delivery"
	resourceDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
	userDomain "github.com/Koshsky/erp-backend/internal/user/domain"
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
	_ UserService       = (*userDomain.UserService)(nil)
	_ TaskService       = (*taskDomain.TaskService)(nil)
	_ ResourceService   = (*resourceDomain.ResourceService)(nil)
	_ ProjectService    = (*projectDomain.ProjectService)(nil)
	_ ProcessService    = (*processDomain.ProcessService)(nil)
	_ MilestoneService  = (*milestoneDomain.MilestoneService)(nil)
	_ AssignmentService = (*assignmentDomain.AssignmentService)(nil)
)
