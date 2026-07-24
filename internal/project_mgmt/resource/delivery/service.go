package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/resource/dto"
)

type ResourceService interface {
	ListResources(ctx context.Context) ([]dto.ResourceResponse, error)
	GetResource(ctx context.Context, id int64) (*dto.ResourceResponse, error)
	CreateResource(ctx context.Context, resource dto.CreateResourceRequest) (*dto.ResourceResponse, error)
	DeleteResource(ctx context.Context, id int64) error
	UpdateResource(ctx context.Context, id int64, resource dto.UpdateResourceRequest) (*dto.ResourceResponse, error)
}
