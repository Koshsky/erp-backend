package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/planning/repository"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
)

type PlanningService struct {
	logger     *slog.Logger
	repository PlanningRepository
}

// NewPlanningService builds the PlanningService service.
func NewPlanningService(logger *slog.Logger, r *repo.PlanningRepository) *PlanningService {
	return &PlanningService{
		logger:     logger,
		repository: r,
	}
}

func (s *PlanningService) GetProjectPlanning(
	ctx context.Context,
	userID int64,
	role string,
) (*dto.ProjectPlanning, error) {
	projects, err := s.repository.ListProjects(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectPlanning{
		Projects: projects,
	}, nil
}

func (s *PlanningService) GetProcessPlanning(
	ctx context.Context,
	userID int64,
	role string,
) (*dto.ProcessPlanning, error) {
	projects, err := s.repository.ListProjects(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	projectIDs := make([]int64, len(projects))
	for i, project := range projects {
		projectIDs[i] = project.ID
	}
	processes, err := s.repository.ListProcessesByProjectIDs(ctx, projectIDs)
	if err != nil {
		return nil, err
	}

	planning := dto.ProcessPlanning{
		Projects: make([]dto.DetailedProject, len(projects)),
	}

	for i, p := range projects {
		planning.Projects[i] = dto.DetailedProject{
			Project:   p,
			Processes: processes[p.ID],
		}
	}

	return &planning, nil
}

func (s *PlanningService) GetTaskPlanning(
	ctx context.Context,
	userID int64,
	role string,
) (*dto.TaskPlanning, error) {
	processes, err := s.loadProcesses(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	if len(processes) == 0 {
		return &dto.TaskPlanning{
			Processes: []dto.DetailedProcess{},
		}, nil
	}

	milestones, tasks, assignments, resourcesMap, err := s.loadAllData(ctx, processes)
	if err != nil {
		return nil, err
	}

	return s.buildPlanning(processes, milestones, tasks, assignments, resourcesMap), nil
}
