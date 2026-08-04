//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/common/nullable"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository/sqlc"
)

type ProcessRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewProcessRepository(logger *slog.Logger, pool *pgxpool.Pool) *ProcessRepository {
	return &ProcessRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *ProcessRepository) CreateProcess(ctx context.Context, process domain.Process) (*domain.Process, error) {
	row, err := r.db.CreateProcess(ctx, sqlc.CreateProcessParams{
		ProjectID: process.ProjectID,
		Title:     process.Title,
		StartDate: process.StartDate,
		EndDate:   process.EndDate,
		OwnerID:   nullable.ToInt8(process.OwnerID),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapProcess(row)
	return &mapped, nil
}

func (r *ProcessRepository) FindProcess(ctx context.Context, id int64) (*domain.Process, error) {
	row, err := r.db.FindProcess(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapProcess(row)
	return &mapped, nil
}

func (r *ProcessRepository) UpdateProcess(ctx context.Context, process domain.Process) (*domain.Process, error) {
	row, err := r.db.UpdateProcess(ctx, sqlc.UpdateProcessParams{
		ProcessID: process.ID,
		OwnerID:   nullable.ToInt8(process.OwnerID),
		ProjectID: process.ProjectID,
		Title:     process.Title,
		StartDate: process.StartDate,
		EndDate:   process.EndDate,
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

func (r *ProcessRepository) ListProcesss(ctx context.Context) ([]domain.Process, error) {
	rows, err := r.db.ListProcesss(ctx)
	if err != nil {
		return nil, err
	}
	processes := make([]domain.Process, 0, len(rows))
	for _, row := range rows {
		processes = append(processes, mapProcess(row))
	}
	return processes, nil
}

func mapProcess(row sqlc.Process) domain.Process {
	return domain.Process{
		ID:        row.ID,
		OwnerID:   nullable.Int64Ptr(row.OwnerID),
		ProjectID: row.ProjectID,
		Title:     row.Title,
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
	}
}
