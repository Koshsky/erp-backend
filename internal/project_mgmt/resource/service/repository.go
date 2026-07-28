package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
)

type RepositoryInterface interface {
	CreateResource(ctx context.Context, Resource domain.Resource) (*domain.Resource, error)
	GetResource(ctx context.Context, id int64) (*domain.Resource, error)
	UpdateResource(ctx context.Context, new domain.Resource) (*domain.Resource, error)
	DeleteResource(ctx context.Context, id int64) error
	ListResources(ctx context.Context) ([]domain.Resource, error)
}
