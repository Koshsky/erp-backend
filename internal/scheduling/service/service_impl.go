package domain

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/scheduling/domain"
	"github.com/Koshsky/erp-backend/internal/scheduling/dto"
)

type SchedulingService struct {
	logger     *slog.Logger
	repository RepositoryInterface
}

func NewSchedulingService(logger *slog.Logger, repository RepositoryInterface) *SchedulingService {
	return &SchedulingService{
		logger:     logger,
		repository: repository,
	}
}

func (s *SchedulingService) GetProjectScheduling(ctx context.Context) (*dto.ProjectScheduling, error) {
	projects, err := s.repository.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	return &dto.ProjectScheduling{
		Projects: projects,
	}, nil
}

func (s *SchedulingService) GetProcessScheduling(ctx context.Context) (*dto.ProcessScheduling, error) {
	projects, err := s.repository.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	scheduling := &dto.ProcessScheduling{
		Projects: make(map[int64]domain.Project),
	}
	projectIDs := make([]int64, len(projects))
	for i, project := range projects {
		projectIDs[i] = project.ID
		scheduling.Projects[project.ID] = project
	}
	processes, err := s.repository.ListProcessesByProjectID(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	scheduling.Processes = processes
	return scheduling, nil
}

func (s *SchedulingService) GetTaskScheduling(ctx context.Context) (*dto.TaskScheduling, error) {
	scheduling, err := s.repository.GetTaskScheduling(ctx)
	if err != nil {
		return nil, err
	}
	return scheduling, nil
}
