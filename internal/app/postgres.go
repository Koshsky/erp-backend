package app

import (
	assignmentDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	milestoneRepoPkg "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	processDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	processRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	projectDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	projectRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	resourceDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	resourceRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/resource/repository"
	taskDomain "github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	taskRepo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	schedulingRepo "github.com/Koshsky/erp-backend/internal/scheduling/repository"
	schedulingService "github.com/Koshsky/erp-backend/internal/scheduling/service"
	userDomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
)

type (
	SchedulingRepository = schedulingService.RepositoryInterface
	UserRepository       = userDomain.RepositoryInterface
	TaskRepository       = taskDomain.RepositoryInterface
	ResourceRepository   = resourceDomain.RepositoryInterface
	ProjectRepository    = projectDomain.RepositoryInterface
	ProcessRepository    = processDomain.RepositoryInterface
	MilestoneRepository  = milestoneDomain.RepositoryInterface
	AssignmentRepository = assignmentDomain.RepositoryInterface
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
