//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository/sqlc"
)

type MilestoneRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewMilestoneRepository(logger *slog.Logger, pool *pgxpool.Pool) *MilestoneRepository {
	return &MilestoneRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *MilestoneRepository) CreateMilestone(ctx context.Context, milestone domain.Milestone) (*domain.Milestone, error) {
	row, err := r.db.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		ProcessID: milestone.ProcessID,
		Title:     milestone.Title,
		Content:   milestone.Content,
		Date:      milestone.Date,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapMilestone(row)
	return &mapped, nil
}

func (r *MilestoneRepository) GetMilestone(ctx context.Context, id int64) (*domain.Milestone, error) {
	row, err := r.db.GetMilestone(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapMilestone(row)
	return &mapped, nil
}

func (r *MilestoneRepository) UpdateMilestone(ctx context.Context, milestone domain.Milestone) (*domain.Milestone, error) {
	row, err := r.db.UpdateMilestone(ctx, sqlc.UpdateMilestoneParams{
		MilestoneID: milestone.ID,
		Date:        milestone.Date,
		Title:       milestone.Title,
		Content:     milestone.Content,
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

func (r *MilestoneRepository) ListMilestones(ctx context.Context) ([]domain.Milestone, error) {
	rows, err := r.db.ListMilestones(ctx)
	if err != nil {
		return nil, err
	}
	milestones := make([]domain.Milestone, 0, len(rows))
	for _, row := range rows {
		milestones = append(milestones, mapMilestone(row))
	}
	return milestones, nil
}

func mapMilestone(row sqlc.Milestone) domain.Milestone {
	return domain.Milestone{
		ID:        row.ID,
		ProcessID: row.ProcessID,
		Title:     row.Title,
		Content:   row.Content,
		Date:      row.Date,
	}
}
