package app

import (
	planningRepo "github.com/Koshsky/erp-backend/internal/planning/repository"
	planningService "github.com/Koshsky/erp-backend/internal/planning/service"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	assignmentService "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
	milestoneRepoPkg "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	milestoneService "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/service"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	processService "github.com/Koshsky/erp-backend/internal/project_mgmt/process/service"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	projectService "github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	taskService "github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
	calendarRepo "github.com/Koshsky/erp-backend/internal/timesheet/calendar/repository"
	calendarService "github.com/Koshsky/erp-backend/internal/timesheet/calendar/service"
	employeeRepo "github.com/Koshsky/erp-backend/internal/timesheet/employee/repository"
	employeeService "github.com/Koshsky/erp-backend/internal/timesheet/employee/service"
	resourceRepo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"
	resourceService "github.com/Koshsky/erp-backend/internal/timesheet/resource/service"
	stateRepo "github.com/Koshsky/erp-backend/internal/timesheet/state/repository"
	stateService "github.com/Koshsky/erp-backend/internal/timesheet/state/service"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
	userService "github.com/Koshsky/erp-backend/internal/user/service"
)

type (
	PlanningRepository   = planningService.PlanningRepository
	UserRepository       = userService.UserRepository
	TaskRepository       = taskService.TaskRepository
	ResourceRepository   = resourceService.ResourceRepository
	ProjectRepository    = projectService.ProjectRepository
	ProcessRepository    = processService.ProcessRepository
	MilestoneRepository  = milestoneService.MilestoneRepository
	AssignmentRepository = assignmentService.AssignmentRepository
	StateRepository      = stateService.StateRepository
	EmployeeRepository   = employeeService.EmployeeRepository
	CalendarRepository   = calendarService.CalendarRepository
)

var (
	_ PlanningRepository   = (*planningRepo.PlanningRepository)(nil)
	_ UserRepository       = (*userRepo.UserRepository)(nil)
	_ TaskRepository       = (*taskRepo.TaskRepository)(nil)
	_ ResourceRepository   = (*resourceRepo.ResourceRepository)(nil)
	_ ProjectRepository    = (*projectRepo.ProjectRepository)(nil)
	_ ProcessRepository    = (*processRepo.ProcessRepository)(nil)
	_ MilestoneRepository  = (*milestoneRepoPkg.MilestoneRepository)(nil)
	_ AssignmentRepository = (*assignmentRepo.AssignmentRepository)(nil)
	_ StateRepository      = (*stateRepo.StateRepository)(nil)
	_ EmployeeRepository   = (*employeeRepo.EmployeeRepository)(nil)
	_ CalendarRepository   = (*calendarRepo.CalendarRepository)(nil)
)
