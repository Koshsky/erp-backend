package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type ResourceRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewResourceRepository(logger *slog.Logger, db *sqlc.Queries) *ResourceRepository {
	return &ResourceRepository{logger: logger, db: db}
}

func (R *ResourceRepository) ListResources(ctx context.Context) ([]domain.Resource, error) {
	rows, err := R.db.ListResources(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]domain.Resource, 0, len(rows))
	for _, row := range rows {
		resources = append(resources, mapResource(row))
	}
	return resources, nil
}

func (r *ResourceRepository) CreateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error) {
	row, err := r.db.CreateResource(ctx, sqlc.CreateResourceParams{
		Title:    resource.Title,
		Code:     resource.Code,
		Quantity: int32(resource.Quantity),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapResource(row)
	return &mapped, nil
}

func (r *ResourceRepository) GetResource(ctx context.Context, id int64) (*domain.Resource, error) {
	row, err := r.db.GetResource(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapResource(row)
	return &mapped, nil
}

func (r *ResourceRepository) UpdateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error) {
	row, err := r.db.UpdateResource(ctx, sqlc.UpdateResourceParams{
		ResourceID: resource.ID,
		Title:      resource.Title,
		Code:       resource.Code,
		Quantity:   int32(resource.Quantity),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapResource(row)
	return &mapped, nil
}

func (r *ResourceRepository) DeleteResource(ctx context.Context, id int64) error {
	return r.db.DeleteResource(ctx, id)
}
