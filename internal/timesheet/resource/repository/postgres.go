//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type ResourceRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewResourceRepository builds the ResourceRepository repository.
func NewResourceRepository(logger *slog.Logger, pool *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
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

func (r *ResourceRepository) ListResources(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
	limit, offset int,
) ([]domain.Resource, error) {
	rows, err := r.db.ListResources(ctx, sqlc.ListResourcesParams{
		Role:       role,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
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

func (r *ResourceRepository) CountResources(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
) (int64, error) {
	return r.db.CountResources(ctx, sqlc.CountResourcesParams{Role: role, UserID: userID, OwnerID: ownerID})
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

// withEmployeesCount enriches the resource model with the members count.
func (r *ResourceRepository) withEmployeesCount(ctx context.Context, row sqlc.Resource) (*domain.Resource, error) {
	count, err := r.db.CountMembersByResourceID(ctx, row.ID)
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

// ListMembersByResourceID returns the users attached to a resource.
func (r *ResourceRepository) ListMembersByResourceID(
	ctx context.Context,
	resourceID int64,
) ([]domain.ResourceMember, error) {
	rows, err := r.db.ListMembersByResourceID(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	members := make([]domain.ResourceMember, 0, len(rows))
	for _, row := range rows {
		members = append(members, mapMember(row))
	}
	return members, nil
}

func (r *ResourceRepository) AddMember(ctx context.Context, resourceID, userID int64) error {
	return r.db.AddMember(ctx, sqlc.AddMemberParams{ResourceID: resourceID, UserID: userID})
}

func (r *ResourceRepository) RemoveMember(ctx context.Context, resourceID, userID int64) error {
	return r.db.RemoveMember(ctx, sqlc.RemoveMemberParams{ResourceID: resourceID, UserID: userID})
}

// FindUserManager returns the manager id of a user (nil — no manager).
func (r *ResourceRepository) FindUserManager(ctx context.Context, userID int64) (*int64, error) {
	managerID, err := r.db.FindUserManager(ctx, userID)
	if err != nil {
		return nil, err
	}
	return nullable.Int64Ptr(managerID), nil
}

func mapMember(row sqlc.ListMembersByResourceIDRow) domain.ResourceMember {
	return domain.ResourceMember{
		ID:              row.ID,
		Name:            row.Name,
		Role:            row.Role,
		Position:        row.Position,
		ManagerID:       nullable.Int64Ptr(row.ManagerID),
		HireDate:        fromDate(row.HireDate),
		TerminationDate: fromDate(row.TerminationDate),
	}
}

// ListAbsence returns the absence ranges (is_available=false states) of the
// resource members overlapping [start, end].
func (r *ResourceRepository) ListAbsence(
	ctx context.Context,
	resourceID int64,
	start, end time.Time,
) ([]domain.ResourceAbsence, error) {
	rows, err := r.db.ListResourceAbsence(ctx, sqlc.ListResourceAbsenceParams{
		ResourceID: resourceID,
		StartDate:  start,
		EndDate:    end,
	})
	if err != nil {
		return nil, err
	}
	absences := make([]domain.ResourceAbsence, 0, len(rows))
	for _, row := range rows {
		absences = append(absences, domain.ResourceAbsence{
			UserID:    row.UserID,
			UserName:  row.UserName,
			StateID:   row.StateID,
			StateCode: row.StateCode,
			StateName: row.StateName,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		})
	}
	return absences, nil
}

// fromDate unwraps a nullable date (pgtype.Date) into [time.Time].
func fromDate(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
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
