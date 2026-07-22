package service

import (
	"log/slog"

	"github.com/Koshsky/erp/api/internal/security/password"
)

type Repository interface {
	ProjectRepository
	ProcessRepository
	MilestoneRepository
	TaskRepository
	ResourceRepository
	AssignmentRepository
	UserRepository
	Close()
}

type Service struct {
	logger *slog.Logger

	*ProjectService
	*ProcessService
	*MilestoneService
	*TaskService
	*ResourceService
	*AssignmentService
	*UserService
}

func New(logger *slog.Logger, repository Repository) *Service {
	validator := NewValidator()
	hasher := password.NewBcryptHasher()

	return &Service{
		logger:            logger,
		ProjectService:    NewProjectService(logger, repository, validator),
		ProcessService:    NewProcessService(logger, repository, validator),
		MilestoneService:  NewMilestoneService(logger, repository, validator),
		TaskService:       NewTaskService(logger, repository, validator),
		ResourceService:   NewResourceService(logger, repository, validator),
		AssignmentService: NewAssignmentService(logger, repository, validator),
		UserService:       NewUserService(logger, repository, validator, hasher),
	}
}
