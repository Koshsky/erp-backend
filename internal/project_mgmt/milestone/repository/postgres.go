//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
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
	milestone domain.Milestone,
) (*domain.Milestone, error) {
	row, err := r.db.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		ProcessID: milestone.ProcessID,
		Title:     milestone.Title,
		Content:   milestone.Content,
		Color:     nullable.ToString(milestone.Color),
		Date:      milestone.Date,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapMilestone(row)
	return &mapped, nil
}

func (r *MilestoneRepository) FindMilestone(ctx context.Context, id int64) (*domain.Milestone, error) {
	row, err := r.db.FindMilestone(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapMilestone(row)
	return &mapped, nil
}

func (r *MilestoneRepository) UpdateMilestone(
	ctx context.Context,
	milestone domain.Milestone,
) (*domain.Milestone, error) {
	row, err := r.db.UpdateMilestone(ctx, sqlc.UpdateMilestoneParams{
		MilestoneID: milestone.ID,
		ProcessID:   milestone.ProcessID,
		Date:        milestone.Date,
		Title:       milestone.Title,
		Content:     milestone.Content,
		Color:       nullable.ToString(milestone.Color),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapMilestone(row)
	return &mapped, nil
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
) ([]domain.Milestone, error) {
	rows, err := r.db.ListMilestones(ctx, sqlc.ListMilestonesParams{
		ScopeView:  viewScope,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	milestones := make([]domain.Milestone, 0, len(rows))
	for _, row := range rows {
		milestones = append(milestones, mapMilestone(row))
	}
	return milestones, nil
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

func mapMilestone(row sqlc.Milestone) domain.Milestone {
	return domain.Milestone{
		ID:        row.ID,
		ProcessID: row.ProcessID,
		Title:     row.Title,
		Content:   row.Content,
		Color:     nullable.StringPtr(row.Color),
		Date:      row.Date,
	}
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *MilestoneRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: row.ProjectOwner, ProcessOwner: row.ProcessOwner}, nil
}
