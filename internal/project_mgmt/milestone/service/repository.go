package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
)

type MilestoneRepository interface {
	CreateMilestone(ctx context.Context, Milestone domain.Milestone) (*domain.Milestone, error)
	FindMilestone(ctx context.Context, id int64) (*domain.Milestone, error)
	UpdateMilestone(ctx context.Context, new domain.Milestone) (*domain.Milestone, error)
	DeleteMilestone(ctx context.Context, id int64) error
	ListMilestones(ctx context.Context) ([]domain.Milestone, error)
}
