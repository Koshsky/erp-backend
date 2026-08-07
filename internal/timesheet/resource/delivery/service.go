package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
)

type ResourceService interface {
	ListResources(ctx context.Context, userID int64, role string) ([]dto.ResourceResponse, error)
	FindResource(ctx context.Context, id int64, userID int64, role string) (*dto.ResourceResponse, error)
	CreateResource(
		ctx context.Context,
		resource dto.CreateResourceRequest,
		userID int64,
	) (*dto.ResourceResponse, error)
	DeleteResource(ctx context.Context, id int64, userID int64, role string) error
	UpdateResource(
		ctx context.Context,
		id int64,
		resource dto.UpdateResourceRequest,
		userID int64,
		role string,
	) (*dto.ResourceResponse, error)
}
