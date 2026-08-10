package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
)

type ResourceRepository interface {
	CreateResource(ctx context.Context, Resource domain.Resource) (*domain.Resource, error)
	FindResource(ctx context.Context, id int64) (*domain.Resource, error)
	UpdateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error)
	DeleteResource(ctx context.Context, id int64) error
	ListResources(ctx context.Context) ([]domain.Resource, error)
	ListResourcesByOwnerID(ctx context.Context, ownerID int64) ([]domain.Resource, error)
}
