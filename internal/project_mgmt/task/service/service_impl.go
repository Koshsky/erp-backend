package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type TaskService struct {
	logger     *slog.Logger
	repository TaskRepository
	mapper     *TaskMapper
	validator  *TaskValidator
}

// NewTaskService builds the TaskService service.
func NewTaskService(logger *slog.Logger, r *repo.TaskRepository) *TaskService {
	return &TaskService{
		logger:     logger,
		repository: r,
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
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrTaskNotFound
		}
		return nil, err
	}
	if task == nil {
		return nil, errors.ErrTaskNotFound
	}
	return s.mapper.ToDTO(task), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, req dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.repository.FindTask(ctx, id)
	if err != nil || task == nil {
		return nil, errors.ErrTaskNotFound
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
		return errors.ErrTaskNotFound
	}

	return s.repository.DeleteTask(ctx, id)
}

func (s *TaskService) ListTasks(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
	limit, offset int,
) ([]dto.TaskResponse, int64, error) {
	rows, err := s.repository.ListTasks(ctx, userID, role, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountTasks(ctx, userID, role, ownerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}
