//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/repository/sqlc"
)

type ResourceRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func (r *ResourceRepository) CreateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error) {
	row, err := r.db.CreateResource(ctx, sqlc.CreateResourceParams{
		Title:   resource.Title,
		Code:    resource.Code,
		OwnerID: ownerIDValue(resource.OwnerID),
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
		OwnerID:        &row.OwnerID,
		EmployeesCount: int(row.EmployeesCount),
	}, nil
}

func (r *ResourceRepository) UpdateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error) {
	row, err := r.db.UpdateResource(ctx, sqlc.UpdateResourceParams{
		ResourceID: resource.ID,
		Title:      resource.Title,
		Code:       resource.Code,
		OwnerID:    ownerIDValue(resource.OwnerID),
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
			OwnerID:        &row.OwnerID,
			EmployeesCount: int(row.EmployeesCount),
		})
	}
	return resources, nil
}

func (r *ResourceRepository) ListResourcesByOwnerID(ctx context.Context, ownerID int64) ([]domain.Resource, error) {
	rows, err := r.db.ListResourcesByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	resources := make([]domain.Resource, 0, len(rows))
	for _, row := range rows {
		resources = append(resources, domain.Resource{
			ID:             row.ID,
			Title:          row.Title,
			Code:           row.Code,
			OwnerID:        &row.OwnerID,
			EmployeesCount: int(row.EmployeesCount),
		})
	}
	return resources, nil
}

// withEmployeesCount enriches the resource model with the active employees count.
func (r *ResourceRepository) withEmployeesCount(ctx context.Context, row sqlc.Resource) (*domain.Resource, error) {
	count, err := r.db.CountEmployeesByResourceID(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	return &domain.Resource{
		ID:             row.ID,
		Title:          row.Title,
		Code:           row.Code,
		OwnerID:        &row.OwnerID,
		EmployeesCount: int(count),
	}, nil
}

// ownerIDValue unwraps a nullable owner into a required value.
func ownerIDValue(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *ResourceRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	owner, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{Owner: owner}, nil
}
