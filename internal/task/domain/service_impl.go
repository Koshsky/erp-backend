package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/task/dto"
)

type TaskService struct {
	logger     *slog.Logger
	repository RepositoryInterface
	mapper     *TaskMapper
	validator  *TaskValidator
}

func NewTaskService(logger *slog.Logger, repository RepositoryInterface) *TaskService {
	return &TaskService{
		logger:     logger,
		repository: repository,
		mapper:     &TaskMapper{},
		validator:  &TaskValidator{},
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	if err := s.validator.ValidateTask(req.ProcessID, req.Title, req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	task := s.mapper.ToDomainFromCreate(req)
	createdTask, err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(createdTask), nil
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*dto.TaskResponse, error) {
	task, err := s.repository.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return s.mapper.ToDTO(task), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.repository.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("Task not found")
	}

	s.mapper.ApplyUpdateToDomain(task, req)
	if err := s.validator.ValidateTask(task.ProcessID, task.Title, task.StartDate, task.EndDate); err != nil {
		return nil, err
	}
	updatedTask, err := s.repository.UpdateTask(ctx, *task)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updatedTask), nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	return s.repository.DeleteTask(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context) ([]dto.TaskResponse, error) {
	rows, err := s.repository.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
