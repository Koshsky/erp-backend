package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error)
	GetTask(ctx context.Context, id int64) (*domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	ListTasksByProcessID(ctx context.Context, processID int64) ([]domain.Task, error)
}

type TaskService struct {
	logger     *slog.Logger
	repository TaskRepository
	mapper     *mapper.TaskMapper
	validator  *Validator
}

func NewTaskService(logger *slog.Logger, repository TaskRepository, validator *Validator) *TaskService {
	return &TaskService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewTaskMapper(),
		validator:  validator,
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
		return nil, ErrTaskNotFound
	}
	return s.mapper.ToDTO(task), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.repository.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
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

func (s *TaskService) ListTasksByProcessID(ctx context.Context, processID int64) ([]dto.TaskResponse, error) {
	tasks, err := s.repository.ListTasksByProcessID(ctx, processID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(tasks), nil
}
