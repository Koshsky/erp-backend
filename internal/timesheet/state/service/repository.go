package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/state/domain"
)

type StateRepository interface {
	CreateState(ctx context.Context, state domain.State) (*domain.State, error)
	FindState(ctx context.Context, id int64) (*domain.State, error)
	UpdateState(ctx context.Context, state domain.State) (*domain.State, error)
	DeleteState(ctx context.Context, id int64) error
	ListStates(ctx context.Context) ([]domain.State, error)
}
