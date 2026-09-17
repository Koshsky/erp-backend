//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	errapi "github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type ProcessRepository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	db     *sqlc.Queries
}

// NewProcessRepository builds the ProcessRepository repository.
func NewProcessRepository(logger *slog.Logger, pool *pgxpool.Pool) *ProcessRepository {
	return &ProcessRepository{
		logger: logger,
		pool:   pool,
		db:     sqlc.New(pool),
	}
}

func (r *ProcessRepository) CreateProcess(
	ctx context.Context,
	projectID int64,
	title string,
	color *string,
	ownerID *int64,
	startDate, endDate time.Time,
) (*sqlc.Process, error) {
	row, err := r.db.CreateProcess(ctx, sqlc.CreateProcessParams{
		ProjectID: projectID,
		Title:     title,
		Color:     nullable.ToString(color),
		StartDate: startDate,
		EndDate:   endDate,
		OwnerID:   nullable.ToInt8(ownerID),
	})
	if err != nil {
		return nil, errapi.MapPgConstraint(errapi.FromPgInvalidParam(err))
	}

	return &row, nil
}

func (r *ProcessRepository) FindProcess(ctx context.Context, id int64) (*sqlc.Process, error) {
	row, err := r.db.FindProcess(ctx, id)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *ProcessRepository) UpdateProcess(ctx context.Context, process sqlc.Process) (*sqlc.Process, error) {
	row, err := r.db.UpdateProcess(ctx, sqlc.UpdateProcessParams{
		ProcessID: process.ID,
		OwnerID:   process.OwnerID,
		ProjectID: process.ProjectID,
		Title:     process.Title,
		Color:     process.Color,
		StartDate: process.StartDate,
		EndDate:   process.EndDate,
	})
	if err != nil {
		return nil, errapi.MapPgConstraint(errapi.FromPgInvalidParam(err))
	}

	return &row, nil
}

func (r *ProcessRepository) DeleteProcess(ctx context.Context, id int64) error {
	return r.db.DeleteProcess(ctx, id)
}

func (r *ProcessRepository) ListProcesss(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
	limit, offset int,
) ([]sqlc.Process, error) {
	return r.db.ListProcesss(ctx, sqlc.ListProcesssParams{
		ScopeView:  viewScope,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
}

func (r *ProcessRepository) CountProcesses(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
) (int64, error) {
	return r.db.CountProcesses(
		ctx,
		sqlc.CountProcessesParams{
			ScopeView: viewScope,
			UserID:    userID,
			OwnerID:   ownerID,
		},
	)
}

// ListProcessIDsByProject returns the active process ids of a project in their
// display order — to validate a reorder request covers the whole group.
func (r *ProcessRepository) ListProcessIDsByProject(ctx context.Context, projectID int64) ([]int64, error) {
	return r.db.ListProcessIdsByProject(ctx, projectID)
}

// ReorderProcesses rewrites the sort_order of the given process ids by list
// position (1-based) in one transaction. The caller validates that the ids
// cover the whole group. The two-phase UPDATE parks the rows on offset slots
// first, because a single-statement value swap would transiently violate the
// partial unique index (project_id, sort_order).
func (r *ProcessRepository) ReorderProcesses(ctx context.Context, ids []int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := r.db.WithTx(tx)
	if err = q.ReorderProcessesMark(ctx, ids); err != nil {
		return err
	}
	if err = q.ReorderProcessesApply(ctx, ids); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *ProcessRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: row.ProjectOwner, ProcessOwner: row.ProcessOwner}, nil
}
