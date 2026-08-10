package server

import (
	assignmentDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/delivery"
	assignmentService "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
	milestoneDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/delivery"
	milestoneService "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/service"
	processDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/process/delivery"
	processService "github.com/Koshsky/erp-backend/internal/project_mgmt/process/service"
	projectDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/project/delivery"
	projectService "github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
	taskDelivery "github.com/Koshsky/erp-backend/internal/project_mgmt/task/delivery"
	taskService "github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
	calendarDelivery "github.com/Koshsky/erp-backend/internal/timesheet/calendar/delivery"
	calendarService "github.com/Koshsky/erp-backend/internal/timesheet/calendar/service"
	employeeDelivery "github.com/Koshsky/erp-backend/internal/timesheet/employee/delivery"
	employeeService "github.com/Koshsky/erp-backend/internal/timesheet/employee/service"
	resourceDelivery "github.com/Koshsky/erp-backend/internal/timesheet/resource/delivery"
	resourceService "github.com/Koshsky/erp-backend/internal/timesheet/resource/service"
	stateDelivery "github.com/Koshsky/erp-backend/internal/timesheet/state/delivery"
	stateService "github.com/Koshsky/erp-backend/internal/timesheet/state/service"
	userDelivery "github.com/Koshsky/erp-backend/internal/user/delivery"
	userService "github.com/Koshsky/erp-backend/internal/user/service"
)

type (
	UserService       = userDelivery.UserService
	TaskService       = taskDelivery.TaskService
	ResourceService   = resourceDelivery.ResourceService
	ProjectService    = projectDelivery.ProjectService
	ProcessService    = processDelivery.ProcessService
	MilestoneService  = milestoneDelivery.MilestoneService
	AssignmentService = assignmentDelivery.AssignmentService
	StateService      = stateDelivery.StateService
	EmployeeService   = employeeDelivery.EmployeeService
	CalendarService   = calendarDelivery.CalendarService
)

// COMPILATION CHECK.
var (
	_ UserService       = (*userService.UserService)(nil)
	_ TaskService       = (*taskService.TaskService)(nil)
	_ ResourceService   = (*resourceService.ResourceService)(nil)
	_ ProjectService    = (*projectService.ProjectService)(nil)
	_ ProcessService    = (*processService.ProcessService)(nil)
	_ MilestoneService  = (*milestoneService.MilestoneService)(nil)
	_ AssignmentService = (*assignmentService.AssignmentService)(nil)
	_ StateService      = (*stateService.StateService)(nil)
	_ EmployeeService   = (*employeeService.EmployeeService)(nil)
	_ CalendarService   = (*calendarService.CalendarService)(nil)
)
