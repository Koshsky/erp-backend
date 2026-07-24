package domain

import (
	"context"
	"log/slog"
)

type SchedulingService struct {
	logger     *slog.Logger
	repository Repository
}

func NewSchedulingService(logger *slog.Logger, repository Repository) *SchedulingService {
	return &SchedulingService{
		logger:     logger,
		repository: repository,
	}
}

func (s *SchedulingService) GetProjectScheduling(ctx context.Context) (*ProjectScheduling, error) {
	scheduling, err := s.repository.GetProjectScheduling(ctx)
	if err != nil {
		return nil, err
	}
	return scheduling, nil
}

func (s *SchedulingService) GetProcessScheduling(ctx context.Context) (*ProcessScheduling, error) {
	scheduling, err := s.repository.GetProcessScheduling(ctx)
	if err != nil {
		return nil, err
	}
	return scheduling, nil
}

func (s *SchedulingService) GetTaskScheduling(ctx context.Context) (*TaskScheduling, error) {
	scheduling, err := s.repository.GetTaskScheduling(ctx)
	if err != nil {
		return nil, err
	}
	return scheduling, nil
}
