//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/repository/sqlc"
)

type ResourceRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewResourceRepository(logger *slog.Logger, pool *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *ResourceRepository) CreateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error) {
	row, err := r.db.CreateResource(ctx, sqlc.CreateResourceParams{
		Title: resource.Title,
		Code:  resource.Code,
	})
	if err != nil {
		return nil, err
	}

	return r.withEmployeesCount(ctx, row)
}

func (r *ResourceRepository) FindResource(ctx context.Context, id int64) (*domain.Resource, error) {
	row, err := r.db.FindResource(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Resource{
		ID:             row.ID,
		Title:          row.Title,
		Code:           row.Code,
		EmployeesCount: int(row.EmployeesCount),
	}, nil
}

func (r *ResourceRepository) UpdateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error) {
	row, err := r.db.UpdateResource(ctx, sqlc.UpdateResourceParams{
		ResourceID: resource.ID,
		Title:      resource.Title,
		Code:       resource.Code,
	})
	if err != nil {
		return nil, err
	}

	return r.withEmployeesCount(ctx, row)
}

func (r *ResourceRepository) DeleteResource(ctx context.Context, id int64) error {
	return r.db.DeleteResource(ctx, id)
}

func (r *ResourceRepository) ListResources(ctx context.Context) ([]domain.Resource, error) {
	rows, err := r.db.ListResources(ctx)
	if err != nil {
		return nil, err
	}

	resources := make([]domain.Resource, 0, len(rows))
	for _, row := range rows {
		resources = append(resources, domain.Resource{
			ID:             row.ID,
			Title:          row.Title,
			Code:           row.Code,
			EmployeesCount: int(row.EmployeesCount),
		})
	}
	return resources, nil
}

// withEmployeesCount дополняет модель ресурса количеством активных сотрудников.
func (r *ResourceRepository) withEmployeesCount(ctx context.Context, row sqlc.Resource) (*domain.Resource, error) {
	count, err := r.db.CountEmployeesByResourceID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return &domain.Resource{
		ID:             row.ID,
		Title:          row.Title,
		Code:           row.Code,
		EmployeesCount: int(count),
	}, nil
}
