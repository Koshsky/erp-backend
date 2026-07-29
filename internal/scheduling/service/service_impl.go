package domain

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/scheduling/dto"
)

type SchedulingService struct {
	logger     *slog.Logger
	repository SchedulingRepository
}

func NewSchedulingService(logger *slog.Logger, repository SchedulingRepository) *SchedulingService {
	return &SchedulingService{
		logger:     logger,
		repository: repository,
	}
}

func (s *SchedulingService) GetProjectScheduling(ctx context.Context, userID int64, role string) (*dto.ProjectScheduling, error) {
	projects, err := s.repository.ListProjects(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectScheduling{
		Projects: projects,
	}, nil
}

func (s *SchedulingService) GetProcessScheduling(ctx context.Context, userID int64, role string) (*dto.ProcessScheduling, error) {
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

	return &dto.ProcessScheduling{
		Projects:  projects,
		Processes: processes,
	}, nil
}

func (s *SchedulingService) GetTaskScheduling(ctx context.Context, userID int64, role string) (*dto.TaskScheduling, error) {
	processes, err := s.repository.ListProcesses(ctx, userID, role)
	if err != nil {
		return nil, err
	}

	processIDs := make([]int64, len(processes))
	for i, process := range processes {
		processIDs[i] = process.ID
	}
	milestones, err := s.repository.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	tasks, err := s.repository.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}

	taskIDs := make([]int64, 0, len(tasks))
	for _, group := range tasks {
		for _, task := range group {
			taskIDs = append(taskIDs, task.ID)
		}
	}
	assignments, err := s.repository.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}

	resources, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.TaskScheduling{
		Processes:   processes,
		Milestones:  milestones,
		Tasks:       tasks,
		Assignments: assignments,
		Resources:   resources,
	}, nil
}
