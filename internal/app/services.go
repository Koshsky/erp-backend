package app

import (
	assignmentDelivery "github.com/Koshsky/erp/api/internal/assignment/delivery"
	assignmentDomain "github.com/Koshsky/erp/api/internal/assignment/domain"
	milestoneDelivery "github.com/Koshsky/erp/api/internal/milestone/delivery"
	milestoneDomain "github.com/Koshsky/erp/api/internal/milestone/domain"
	processDelivery "github.com/Koshsky/erp/api/internal/process/delivery"
	processDomain "github.com/Koshsky/erp/api/internal/process/domain"
	projectDelivery "github.com/Koshsky/erp/api/internal/project/delivery"
	projectDomain "github.com/Koshsky/erp/api/internal/project/domain"
	resourceDelivery "github.com/Koshsky/erp/api/internal/resource/delivery"
	resourceDomain "github.com/Koshsky/erp/api/internal/resource/domain"
	taskDelivery "github.com/Koshsky/erp/api/internal/task/delivery"
	taskDomain "github.com/Koshsky/erp/api/internal/task/domain"
	userDelivery "github.com/Koshsky/erp/api/internal/user/delivery"
	userDomain "github.com/Koshsky/erp/api/internal/user/domain"
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
