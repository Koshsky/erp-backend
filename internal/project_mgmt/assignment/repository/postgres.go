//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository/sqlc"
)

type AssignmentRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewAssignmentRepository builds the AssignmentRepository repository.
func NewAssignmentRepository(logger *slog.Logger, pool *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *AssignmentRepository) CreateAssignment(
	ctx context.Context,
	taskID, resourceID int64,
	quantity int,
) (*sqlc.Assignment, error) {
	row, err := r.db.CreateAssignment(ctx, sqlc.CreateAssignmentParams{
		TaskID:     taskID,
		ResourceID: resourceID,
		Quantity:   int64(quantity),
	})
	if err != nil {
		// Idempotent create: for an already-active (task_id, resource_id) pair,
		// INSERT ... ON CONFLICT DO NOTHING inserts no row and RETURNING comes
		// back empty. Return the already existing record.
		if errors.Is(err, pgx.ErrNoRows) {
			return r.FindAssignmentByKey(ctx, taskID, resourceID)
		}
		return nil, err
	}

	return &row, nil
}

// FindAssignmentByKey returns the assignment for the active (task_id, resource_id) pair.
func (r *AssignmentRepository) FindAssignmentByKey(
	ctx context.Context,
	taskID, resourceID int64,
) (*sqlc.Assignment, error) {
	row, err := r.db.FindAssignmentByKey(ctx, sqlc.FindAssignmentByKeyParams{
		TaskID:     taskID,
		ResourceID: resourceID,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AssignmentRepository) FindAssignment(ctx context.Context, id int64) (*sqlc.Assignment, error) {
	row, err := r.db.FindAssignment(ctx, id)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *AssignmentRepository) UpdateAssignment(
	ctx context.Context,
	assignment sqlc.Assignment,
	quantity int,
) (*sqlc.Assignment, error) {
	row, err := r.db.UpdateAssignment(ctx, sqlc.UpdateAssignmentParams{
		AssignmentID: assignment.ID,
		TaskID:       assignment.TaskID,
		ResourceID:   assignment.ResourceID,
		Quantity:     int64(quantity),
	})
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *AssignmentRepository) DeleteAssignment(ctx context.Context, id int64) error {
	return r.db.DeleteAssignment(ctx, id)
}

func (r *AssignmentRepository) ListAssignments(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
	limit, offset int,
) ([]sqlc.Assignment, error) {
	return r.db.ListAssigments(ctx, sqlc.ListAssigmentsParams{
		ScopeView:  viewScope,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
}

func (r *AssignmentRepository) CountAssignments(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
) (int64, error) {
	return r.db.CountAssignments(
		ctx,
		sqlc.CountAssignmentsParams{
			ScopeView: viewScope,
			UserID:    userID,
			OwnerID:   ownerID,
		},
	)
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *AssignmentRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: row.ProjectOwner, ProcessOwner: row.ProcessOwner, Owner: row.OwnerID}, nil
}
