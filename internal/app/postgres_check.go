package app

import (
	assignmentService "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	milestoneService "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/service"
	milestoneRepoPkg "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	processService "github.com/Koshsky/erp-backend/internal/project_mgmt/process/service"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	projectService "github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	resourceService "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/service"
	resourceRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/repository"
	taskService "github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	schedulingRepo "github.com/Koshsky/erp-backend/internal/scheduling/repository"
	schedulingService "github.com/Koshsky/erp-backend/internal/scheduling/service"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
	userService "github.com/Koshsky/erp-backend/internal/user/service"
)

type (
	SchedulingRepository = schedulingService.RepositoryInterface
	UserRepository       = userService.RepositoryInterface
	TaskRepository       = taskService.RepositoryInterface
	ResourceRepository   = resourceService.RepositoryInterface
	ProjectRepository    = projectService.RepositoryInterface
	ProcessRepository    = processService.RepositoryInterface
	MilestoneRepository  = milestoneService.RepositoryInterface
	AssignmentRepository = assignmentService.RepositoryInterface
)

var (
	_ SchedulingRepository = (*schedulingRepo.SchedulingRepository)(nil)
	_ UserRepository       = (*userRepo.UserRepository)(nil)
	_ TaskRepository       = (*taskRepo.TaskRepository)(nil)
	_ ResourceRepository   = (*resourceRepo.ResourceRepository)(nil)
	_ ProjectRepository    = (*projectRepo.ProjectRepository)(nil)
	_ ProcessRepository    = (*processRepo.ProcessRepository)(nil)
	_ MilestoneRepository  = (*milestoneRepoPkg.MilestoneRepository)(nil)
	_ AssignmentRepository = (*assignmentRepo.AssignmentRepository)(nil)
)
