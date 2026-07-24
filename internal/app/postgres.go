package app

import (
	assignmentDomain "github.com/Koshsky/erp/api/internal/assignment/domain"
	assignmentRepo "github.com/Koshsky/erp/api/internal/assignment/repository"
	milestoneDomain "github.com/Koshsky/erp/api/internal/milestone/domain"
	milestoneRepoPkg "github.com/Koshsky/erp/api/internal/milestone/repository"
	processDomain "github.com/Koshsky/erp/api/internal/process/domain"
	processRepo "github.com/Koshsky/erp/api/internal/process/repository"
	projectDomain "github.com/Koshsky/erp/api/internal/project/domain"
	projectRepo "github.com/Koshsky/erp/api/internal/project/repository"
	resourceDomain "github.com/Koshsky/erp/api/internal/resource/domain"
	resourceRepo "github.com/Koshsky/erp/api/internal/resource/repository"
	schedulingDomain "github.com/Koshsky/erp/api/internal/scheduling/domain"
	schedulingRepo "github.com/Koshsky/erp/api/internal/scheduling/repository"
	taskDomain "github.com/Koshsky/erp/api/internal/task/domain"
	taskRepo "github.com/Koshsky/erp/api/internal/task/repository"
	userDomain "github.com/Koshsky/erp/api/internal/user/domain"
	userRepo "github.com/Koshsky/erp/api/internal/user/repository"
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
