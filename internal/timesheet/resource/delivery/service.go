package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
)

type ResourceService interface {
	ListResources(
		ctx context.Context,
		userID int64,
		role string,
		ownerID int64,
		limit, offset int,
	) ([]dto.ResourceResponse, int64, error)
	FindResource(ctx context.Context, id int64) (*dto.ResourceResponse, error)
	CreateResource(
		ctx context.Context,
		resource dto.CreateResourceRequest,
		userID int64,
	) (*dto.ResourceResponse, error)
	DeleteResource(ctx context.Context, id int64) error
	UpdateResource(
		ctx context.Context,
		id int64,
		resource dto.UpdateResourceRequest,
	) (*dto.ResourceResponse, error)
}
