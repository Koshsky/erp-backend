package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
	"github.com/Koshsky/erp/api/internal/service/mapper"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProcessRepository interface {
	CreateProcess(ctx context.Context, process domain.Process) (*domain.Process, error)
	GetProcess(ctx context.Context, id int64) (*domain.Process, error)
	UpdateProcess(ctx context.Context, process domain.Process) (*domain.Process, error)
	DeleteProcess(ctx context.Context, id int64) error
	ListProcessesByProjectID(ctx context.Context, projectID int64) ([]domain.Process, error)
	GetProcessTimeline(ctx context.Context, id int64) (*sqlc.GetProcessWithProjectRow, []sqlc.GetProcessTasksWithResourcesRow, error)
}

type ProcessService struct {
	logger     *slog.Logger
	repository ProcessRepository
	mapper     *mapper.ProcessMapper
	validator  *Validator
}

func NewProcessService(logger *slog.Logger, repository ProcessRepository, validator *Validator) *ProcessService {
	return &ProcessService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewProcessMapper(),
		validator:  validator,
	}
}

func (s *ProcessService) CreateProcess(ctx context.Context, req dto.CreateProcessRequest) (*dto.ProcessResponse, error) {
	if err := s.validator.ValidateProcess(req.ProjectID, req.Title, req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	created, err := s.repository.CreateProcess(ctx, s.mapper.ToDomainFromCreate(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(created), nil
}

func (s *ProcessService) GetProcess(ctx context.Context, id int64) (*dto.ProcessResponse, error) {
	process, err := s.repository.GetProcess(ctx, id)
	if err != nil {
		return nil, err
	}
	if process == nil {
		return nil, ErrProcessNotFound
	}
	return s.mapper.ToDTO(process), nil
}

func (s *ProcessService) UpdateProcess(ctx context.Context, id int64, req dto.UpdateProcessRequest) (*dto.ProcessResponse, error) {
	process, err := s.repository.GetProcess(ctx, id)
	if err != nil {
		return nil, err
	}
	if process == nil {
		return nil, ErrProcessNotFound
	}
	s.mapper.ApplyUpdateToDomain(process, req)
	if err := s.validator.ValidateProcess(process.ProjectID, process.Title, process.StartDate, process.EndDate); err != nil {
		return nil, err
	}
	updated, err := s.repository.UpdateProcess(ctx, *process)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updated), nil
}

func (s *ProcessService) DeleteProcess(ctx context.Context, id int64) error {
	return s.repository.DeleteProcess(ctx, id)
}

func (s *ProcessService) ListProcessesByProjectID(ctx context.Context, projectID int64) ([]dto.ProcessResponse, error) {
	domainProcesses, err := s.repository.ListProcessesByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(domainProcesses), nil
}

func (s *ProcessService) GetProcessTimeline(ctx context.Context, id int64) (*dto.TimelineResponse, error) {
	processRow, taskRows, err := s.repository.GetProcessTimeline(ctx, id)
	if err != nil {
		return nil, err
	}

	type ri struct {
		title string
		code  string
		qty   int
	}

	resourceMap := make(map[int64]*ri)
	taskMap := make(map[int64]*dto.TimelineTaskWithResources)
	var taskOrder []int64

	for _, row := range taskRows {
		tid := row.TaskID
		if _, ok := taskMap[tid]; !ok {
			taskMap[tid] = &dto.TimelineTaskWithResources{
				ID:        tid,
				Title:     row.TaskTitle,
				StartDate: fromPgDate(row.TaskStartDate),
				EndDate:   fromPgDate(row.TaskEndDate),
				Resources: []dto.TimelineTaskResource{},
			}
			taskOrder = append(taskOrder, tid)
		}

		if row.ResourceID.Valid {
			rid := row.ResourceID.Int64
			taskMap[tid].Resources = append(taskMap[tid].Resources, dto.TimelineTaskResource{
				ResourceID: rid,
				Quantity:   int(row.AssignmentQuantity.Int32),
				Code:       row.ResourceCode.String,
			})

			if _, ok := resourceMap[rid]; !ok {
				resourceMap[rid] = &ri{
					title: row.ResourceTitle.String,
					code:  row.ResourceCode.String,
					qty:   int(row.ResourceQuantity.Int32),
				}
			}
		}
	}

	tasks := make([]dto.TimelineTaskWithResources, len(taskOrder))
	for i, tid := range taskOrder {
		tasks[i] = *taskMap[tid]
	}

	resources := make([]dto.TimelineResource, 0, len(resourceMap))
	for id, info := range resourceMap {
		resources = append(resources, dto.TimelineResource{
			ID:       id,
			Code:     info.code,
			Title:    info.title,
			Quantity: info.qty,
		})
	}

	return &dto.TimelineResponse{
		Processes: []dto.TimelineProcessItem{
			{
				ID:          processRow.ProcessID,
				Title:       processRow.ProcessTitle,
				StartDate:   fromPgDate(processRow.ProcessStartDate),
				EndDate:     fromPgDate(processRow.ProcessEndDate),
				ProjectCode: processRow.ProjectCode,
				Tasks:       tasks,
			},
		},
		Resources: resources,
	}, nil
}

func fromPgDate(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}
