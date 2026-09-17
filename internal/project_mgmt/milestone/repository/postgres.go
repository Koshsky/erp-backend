//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type MilestoneRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewMilestoneRepository builds the MilestoneRepository repository.
func NewMilestoneRepository(logger *slog.Logger, pool *pgxpool.Pool) *MilestoneRepository {
	return &MilestoneRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *MilestoneRepository) CreateMilestone(
	ctx context.Context,
	processID int64,
	title, content string,
	color *string,
	date time.Time,
) (*sqlc.Milestone, error) {
	row, err := r.db.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		ProcessID: processID,
		Title:     title,
		Content:   content,
		Color:     nullable.ToString(color),
		Date:      date,
	})
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *MilestoneRepository) FindMilestone(ctx context.Context, id int64) (*sqlc.Milestone, error) {
	row, err := r.db.FindMilestone(ctx, id)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *MilestoneRepository) UpdateMilestone(
	ctx context.Context,
	milestone sqlc.Milestone,
) (*sqlc.Milestone, error) {
	row, err := r.db.UpdateMilestone(ctx, sqlc.UpdateMilestoneParams{
		MilestoneID: milestone.ID,
		ProcessID:   milestone.ProcessID,
		Date:        milestone.Date,
		Title:       milestone.Title,
		Content:     milestone.Content,
		Color:       milestone.Color,
	})
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *MilestoneRepository) DeleteMilestone(ctx context.Context, id int64) error {
	return r.db.DeleteMilestone(ctx, id)
}

func (r *MilestoneRepository) ListMilestones(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
	limit, offset int,
) ([]sqlc.Milestone, error) {
	return r.db.ListMilestones(ctx, sqlc.ListMilestonesParams{
		ScopeView:  viewScope,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
}

func (r *MilestoneRepository) CountMilestones(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
) (int64, error) {
	return r.db.CountMilestones(
		ctx,
		sqlc.CountMilestonesParams{
			ScopeView: viewScope,
			UserID:    userID,
			OwnerID:   ownerID,
		},
	)
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *MilestoneRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: row.ProjectOwner, ProcessOwner: row.ProcessOwner}, nil
}
