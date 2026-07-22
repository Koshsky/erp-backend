package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type AssignmentRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewAssignmentRepository(logger *slog.Logger, db *sqlc.Queries) *AssignmentRepository {
	return &AssignmentRepository{logger: logger, db: db}
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

func (r *AssignmentRepository) GetAssignment(ctx context.Context, id int64) (*domain.Assignment, error) {
	row, err := r.db.GetAssignment(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapAssignment(row)
	return &mapped, nil
}

func (r *AssignmentRepository) UpdateAssignment(ctx context.Context, assignment domain.Assignment) (*domain.Assignment, error) {
	row, err := r.db.UpdateAssignment(ctx, sqlc.UpdateAssignmentParams{
		ID:       assignment.ID,
		Quantity: int32(assignment.Quantity),
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

func (r *AssignmentRepository) ListAssignmentsByTaskID(ctx context.Context, taskID int64) ([]domain.Assignment, error) {
	rows, err := r.db.ListAssignmentsByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	assignments := make([]domain.Assignment, 0, len(rows))
	for _, row := range rows {
		assignments = append(assignments, mapAssignment(row))
	}

	return assignments, nil
}
