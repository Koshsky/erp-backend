package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type MilestoneRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewMilestoneRepository(logger *slog.Logger, db *sqlc.Queries) *MilestoneRepository {
	return &MilestoneRepository{logger: logger, db: db}
}

func (r *MilestoneRepository) CreateMilestone(ctx context.Context, milestone domain.Milestone) (*domain.Milestone, error) {
	row, err := r.db.CreateMilestone(ctx, sqlc.CreateMilestoneParams{
		ProcessID: milestone.ProcessID,
		Title:     milestone.Title,
		Content:   milestone.Content,
		Date:      toDate(milestone.Date),
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
		Date:        toDate(milestone.Date),
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
