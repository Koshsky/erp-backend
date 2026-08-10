package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/resource/repository"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type ResourceService struct {
	logger     *slog.Logger
	repository ResourceRepository
	mapper     *ResourceMapper
	validator  *ResourceValidator
}

// NewResourceService builds the ResourceService service.
func NewResourceService(logger *slog.Logger, r *repo.ResourceRepository) *ResourceService {
	return &ResourceService{
		logger:     logger,
		repository: r,
		mapper:     NewResourceMapper(),
		validator:  &ResourceValidator{},
	}
}

func (s *ResourceService) ListResources(ctx context.Context, limit, offset int) ([]dto.ResourceResponse, int64, error) {
	rows, err := s.repository.ListResources(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountResources(ctx)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}

// CreateResource creates a resource. The middleware checked permissions; here
// only owner normalization: if not set, the creator becomes the owner
// (owner_id is required).
func (s *ResourceService) CreateResource(
	ctx context.Context,
	req dto.CreateResourceRequest,
	userID int64,
) (*dto.ResourceResponse, error) {
	if req.OwnerID == nil {
		req.OwnerID = &userID
	}

	resource := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateResource(&resource); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateResource(ctx, resource)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *ResourceService) FindResource(ctx context.Context, id int64) (*dto.ResourceResponse, error) {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrResourceNotFound
		}
		return nil, err
	}
	if resource == nil {
		return nil, errors.ErrResourceNotFound
	}
	return s.mapper.ToDTO(resource), nil
}

func (s *ResourceService) UpdateResource(
	ctx context.Context,
	id int64,
	req dto.UpdateResourceRequest,
) (*dto.ResourceResponse, error) {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil || resource == nil {
		return nil, errors.ErrResourceNotFound
	}

	s.mapper.ApplyUpdateToDomain(resource, req)
	if err = s.validator.ValidateResource(resource); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateResource(ctx, *resource)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, id int64) error {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil || resource == nil {
		return errors.ErrResourceNotFound
	}

	return s.repository.DeleteResource(ctx, id)
}
