package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository/sqlc"
)

type MilestoneRepository interface {
	CreateMilestone(
		ctx context.Context,
		processID int64,
		title, content string,
		color *string,
		date time.Time,
	) (*sqlc.Milestone, error)
	FindMilestone(ctx context.Context, id int64) (*sqlc.Milestone, error)
	UpdateMilestone(ctx context.Context, milestone sqlc.Milestone) (*sqlc.Milestone, error)
	DeleteMilestone(ctx context.Context, id int64) error
	ListMilestones(
		ctx context.Context,
		userID int64,
		viewScope string,
		ownerID int64,
		limit, offset int,
	) ([]sqlc.Milestone, error)
	CountMilestones(ctx context.Context, userID int64, viewScope string, ownerID int64) (int64, error)
}
