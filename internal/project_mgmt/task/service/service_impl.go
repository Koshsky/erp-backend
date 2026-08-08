package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type TaskService struct {
	logger     *slog.Logger
	repository TaskRepository
	mapper     *TaskMapper
	validator  *TaskValidator
}

func NewTaskService(logger *slog.Logger, repository TaskRepository) *TaskService {
	return &TaskService{
		logger:     logger,
		repository: repository,
		mapper:     &TaskMapper{},
		validator:  &TaskValidator{},
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	task := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateTask(&task); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *TaskService) FindTask(ctx context.Context, id int64) (*dto.TaskResponse, error) {
	task, err := s.repository.FindTask(ctx, id)
	if err != nil {
		if validator.IsNotFoundError(err) {
			return nil, validator.ErrTaskNotFound
		}
		return nil, err
	}
	if task == nil {
		return nil, validator.ErrTaskNotFound
	}
	return s.mapper.ToDTO(task), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.repository.FindTask(ctx, id)
	if err != nil || task == nil {
		return nil, validator.ErrTaskNotFound
	}

	s.mapper.ApplyUpdateToDomain(task, req)
	if err = s.validator.ValidateTask(task); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateTask(ctx, *task)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	task, err := s.repository.FindTask(ctx, id)
	if err != nil || task == nil {
		return validator.ErrTaskNotFound
	}

	return s.repository.DeleteTask(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context) ([]dto.TaskResponse, error) {
	rows, err := s.repository.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
