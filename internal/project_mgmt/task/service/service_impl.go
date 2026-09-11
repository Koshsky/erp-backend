package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/order"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type TaskService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository TaskRepository
	mapper     *TaskMapper
	validator  *TaskValidator
}

// NewTaskService builds the TaskService service.
func NewTaskService(logger *slog.Logger, tracer *tracingpkg.Tracer, r *repo.TaskRepository) *TaskService {
	return &TaskService{
		logger:     logger,
		tracer:     tracer,
		repository: r,
		mapper:     &TaskMapper{},
		validator:  &TaskValidator{},
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	ctx, end := s.tracer.Start(ctx, "task.CreateTask")
	defer end(nil)

	task := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateTask(&task); err != nil {
		return nil, err
	}

	// Subtask (operation) semantics:
	//  - the parent must be a top-level task (a subtask cannot be a parent);
	//  - process_id always equals the parent's process (RBAC unchanged);
	//  - dates inherit the parent's bounds (subtask spans the parent's
	//    interval; the frontend never renders subtask dates).
	if err := s.applyParent(ctx, &task, req.ParentID); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

// applyParent resolves the parent for a subtask request: inherits the
// parent's process and dates and rejects a parent that is itself a subtask.
func (s *TaskService) applyParent(ctx context.Context, task *domain.Task, parentID *int64) error {
	if parentID == nil {
		return nil
	}
	parent, err := s.repository.FindTask(ctx, *parentID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return errors.NewValidationError("родительская задача не найдена")
		}
		return err
	}
	if parent == nil {
		return errors.NewValidationError("родительская задача не найдена")
	}
	if parent.ParentID != nil {
		return errors.NewValidationError("подзадача не может иметь собственные подзадачи")
	}
	task.ProcessID = parent.ProcessID
	task.StartDate = parent.StartDate
	task.EndDate = parent.EndDate
	return nil
}

func (s *TaskService) FindTask(ctx context.Context, id int64) (*dto.TaskResponse, error) {
	ctx, end := s.tracer.Start(ctx, "task.FindTask")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "task.UpdateTask")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "task.DeleteTask")
	defer end(nil)

	task, err := s.repository.FindTask(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // idempotent delete: already deleted — not an error
		}
		return err
	}
	if task == nil {
		return nil // idempotent delete
	}

	return s.repository.DeleteTask(ctx, id)
}

func (s *TaskService) ListTasks(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
	limit, offset int,
) ([]dto.TaskResponse, int64, error) {
	ctx, end := s.tracer.Start(ctx, "task.ListTasks")
	defer end(nil)

	rows, err := s.repository.ListTasks(ctx, userID, viewScope, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountTasks(ctx, userID, viewScope, ownerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}

// ReorderTasks applies a new order to all active tasks of a process: the
// request carries the complete ordered id list, the server validates it covers
// the whole group and rewrites sort_order by list position.
func (s *TaskService) ReorderTasks(
	ctx context.Context,
	req dto.ReorderTaskRequest,
) error {
	ctx, end := s.tracer.Start(ctx, "task.ReorderTasks")
	defer end(nil)

	if len(req.IDs) == 0 {
		return errors.NewValidationError("список id задач пуст")
	}
	if err := order.RejectDuplicateIDs("задач", req.IDs); err != nil {
		return err
	}

	current, err := s.repository.ListTaskIDsByProcess(ctx, req.ProcessID)
	if err != nil {
		return err
	}
	if !order.SameIDSet(req.IDs, current) {
		return errors.NewValidationError(
			"список должен содержать все активные задачи процесса без изменений состава",
		)
	}

	return s.repository.ReorderTasks(ctx, req.IDs)
}
