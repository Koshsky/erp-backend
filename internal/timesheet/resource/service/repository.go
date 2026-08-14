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
	ListResources(
		ctx context.Context,
		userID int64,
		role string,
		ownerID int64,
		limit, offset int,
	) ([]domain.Resource, error)
	CountResources(ctx context.Context, userID int64, role string, ownerID int64) (int64, error)
	ListResourcesByOwnerID(ctx context.Context, ownerID int64) ([]domain.Resource, error)
	ListMembersByResourceID(ctx context.Context, resourceID int64) ([]domain.ResourceMember, error)
	AddMember(ctx context.Context, resourceID, userID int64) error
	RemoveMember(ctx context.Context, resourceID, userID int64) error
	FindUserManager(ctx context.Context, userID int64) (*int64, error)
}
