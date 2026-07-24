package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateResource(ctx context.Context, Resource Resource) (*Resource, error)
	GetResource(ctx context.Context, id int64) (*Resource, error)
	UpdateResource(ctx context.Context, new Resource) (*Resource, error)
	DeleteResource(ctx context.Context, id int64) error
	ListResources(ctx context.Context) ([]Resource, error)
}
