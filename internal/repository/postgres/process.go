package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type ProcessRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewProcessRepository(logger *slog.Logger, db *sqlc.Queries) *ProcessRepository {
	return &ProcessRepository{logger: logger, db: db}
}

func (r *ProcessRepository) CreateProcess(ctx context.Context, process domain.Process) (*domain.Process, error) {
	row, err := r.db.CreateProcess(ctx, sqlc.CreateProcessParams{
		ProjectID: process.ProjectID,
		Title:     process.Title,
		StartDate: toDate(process.StartDate),
		EndDate:   toDate(process.EndDate),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapProcess(row)
	return &mapped, nil
}

func (r *ProcessRepository) GetProcess(ctx context.Context, id int64) (*domain.Process, error) {
	row, err := r.db.GetProcess(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapProcess(row)
	return &mapped, nil
}

func (r *ProcessRepository) UpdateProcess(ctx context.Context, process domain.Process) (*domain.Process, error) {
	row, err := r.db.UpdateProcess(ctx, sqlc.UpdateProcessParams{
		ID:        process.ID,
		Title:     process.Title,
		StartDate: toDate(process.StartDate),
		EndDate:   toDate(process.EndDate),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapProcess(row)
	return &mapped, nil
}

func (r *ProcessRepository) DeleteProcess(ctx context.Context, id int64) error {
	return r.db.DeleteProcess(ctx, id)
}

func (r *ProcessRepository) ListProcessesByProjectID(ctx context.Context, projectID int64) ([]domain.Process, error) {
	rows, err := r.db.ListProcessesByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	processes := make([]domain.Process, 0, len(rows))
	for _, row := range rows {
		processes = append(processes, mapProcess(row))
	}

	return processes, nil
}

func (r *ProcessRepository) GetProcessTimeline(ctx context.Context, id int64) (*sqlc.GetProcessWithProjectRow, []sqlc.GetProcessTasksWithResourcesRow, error) {
	processRow, err := r.db.GetProcessWithProject(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	taskRows, err := r.db.GetProcessTasksWithResources(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return &processRow, taskRows, nil
}

