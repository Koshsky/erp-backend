package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateMilestone(ctx context.Context, Milestone Milestone) (*Milestone, error)
	GetMilestone(ctx context.Context, id int64) (*Milestone, error)
	UpdateMilestone(ctx context.Context, new Milestone) (*Milestone, error)
	DeleteMilestone(ctx context.Context, id int64) error
	ListMilestones(ctx context.Context) ([]Milestone, error)
}
