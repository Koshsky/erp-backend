package app

import (
	assignmentDomain "github.com/Koshsky/erp-backend/internal/assignment/domain"
	assignmentRepo "github.com/Koshsky/erp-backend/internal/assignment/repository"
	milestoneDomain "github.com/Koshsky/erp-backend/internal/milestone/domain"
	milestoneRepoPkg "github.com/Koshsky/erp-backend/internal/milestone/repository"
	processDomain "github.com/Koshsky/erp-backend/internal/process/domain"
	processRepo "github.com/Koshsky/erp-backend/internal/process/repository"
	projectDomain "github.com/Koshsky/erp-backend/internal/project/domain"
	projectRepo "github.com/Koshsky/erp-backend/internal/project/repository"
	resourceDomain "github.com/Koshsky/erp-backend/internal/resource/domain"
	resourceRepo "github.com/Koshsky/erp-backend/internal/resource/repository"
	schedulingDomain "github.com/Koshsky/erp-backend/internal/scheduling/domain"
	schedulingRepo "github.com/Koshsky/erp-backend/internal/scheduling/repository"
	taskDomain "github.com/Koshsky/erp-backend/internal/task/domain"
	taskRepo "github.com/Koshsky/erp-backend/internal/task/repository"
	userDomain "github.com/Koshsky/erp-backend/internal/user/domain"
	userRepo "github.com/Koshsky/erp-backend/internal/user/repository"
)

type (
	SchedulingRepository = schedulingDomain.Repository
	UserRepository       = userDomain.RepositoryInterface
	TaskRepository       = taskDomain.Repository
	ResourceRepository   = resourceDomain.RepositoryInterface
	ProjectRepository    = projectDomain.RepositoryInterface
	ProcessRepository    = processDomain.RepositoryInterface
	MilestoneRepository  = milestoneDomain.RepositoryInterface
	AssignmentRepository = assignmentDomain.RepositoryInterface
)

var (
	_ SchedulingRepository = (*schedulingRepo.SchedulingRepository)(nil)
	_ UserRepository       = (*userRepo.Repository)(nil)
	_ TaskRepository       = (*taskRepo.TaskRepository)(nil)
	_ ResourceRepository   = (*resourceRepo.ResourceRepository)(nil)
	_ ProjectRepository    = (*projectRepo.ProjectRepository)(nil)
	_ ProcessRepository    = (*processRepo.ProcessRepository)(nil)
	_ MilestoneRepository  = (*milestoneRepoPkg.MilestoneRepository)(nil)
	_ AssignmentRepository = (*assignmentRepo.AssignmentRepository)(nil)
)
