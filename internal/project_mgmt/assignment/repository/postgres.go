//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository/sqlc"
)

type AssignmentRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewAssignmentRepository(logger *slog.Logger, pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *AssignmentRepository) CreateAssignment(ctx context.Context, assignment domain.Assignment) (*domain.Assignment, error) {
	row, err := r.db.CreateAssignment(ctx, sqlc.CreateAssignmentParams{
		TaskID:     assignment.TaskID,
		ResourceID: assignment.ResourceID,
		Quantity:   int32(assignment.Quantity),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapAssignment(row)
	return &mapped, nil
}

func (r *AssignmentRepository) FindAssignment(ctx context.Context, id int64) (*domain.Assignment, error) {
	row, err := r.db.FindAssignment(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapAssignment(row)
	return &mapped, nil
}

func (r *AssignmentRepository) UpdateAssignment(ctx context.Context, assignment domain.Assignment) (*domain.Assignment, error) {
	row, err := r.db.UpdateAssignment(ctx, sqlc.UpdateAssignmentParams{
		AssignmentID: assignment.ID,
		Quantity:     int32(assignment.Quantity),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapAssignment(row)
	return &mapped, nil
}

func (r *AssignmentRepository) DeleteAssignment(ctx context.Context, id int64) error {
	return r.db.DeleteAssignment(ctx, id)
}

func (r *AssignmentRepository) ListAssignments(ctx context.Context) ([]domain.Assignment, error) {
	rows, err := r.db.ListAssigments(ctx)
	if err != nil {
		return nil, err
	}
	assignments := make([]domain.Assignment, 0, len(rows))
	for _, row := range rows {
		assignments = append(assignments, mapAssignment(row))
	}
	return assignments, nil
}

func mapAssignment(row sqlc.Assignment) domain.Assignment {
	return domain.Assignment{
		ID:         row.ID,
		TaskID:     row.TaskID,
		ResourceID: row.ResourceID,
		Quantity:   int(row.Quantity),
	}
}
