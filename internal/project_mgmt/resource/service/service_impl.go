package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/resource/dto"
)

type ResourceService struct {
	logger     *slog.Logger
	repository ResourceRepository
	mapper     *ResourceMapper
	validator  *ResourceValidator
}

func NewResourceService(logger *slog.Logger, repository ResourceRepository) *ResourceService {
	return &ResourceService{
		logger:     logger,
		repository: repository,
		mapper:     NewResourceMapper(),
		validator:  &ResourceValidator{},
	}
}

func (s *ResourceService) ListResources(ctx context.Context) ([]dto.ResourceResponse, error) {
	resources, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(resources), nil
}

func (s *ResourceService) CreateResource(ctx context.Context, resource dto.CreateResourceRequest) (*dto.ResourceResponse, error) {
	if err := s.validator.ValidateResource(resource.Code, resource.Title, resource.Quantity); err != nil {
		return nil, err
	}
	createdResource, err := s.repository.CreateResource(ctx, s.mapper.ToDomainFromCreate(resource))
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(createdResource), nil
}

func (s *ResourceService) FindResource(ctx context.Context, id int64) (*dto.ResourceResponse, error) {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, fmt.Errorf("resource not found")
	}
	return s.mapper.ToDTO(resource), nil
}

func (s *ResourceService) UpdateResource(ctx context.Context, id int64, req dto.UpdateResourceRequest) (*dto.ResourceResponse, error) {
	resource, err := s.repository.FindResource(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, fmt.Errorf("resource not found")
	}
	s.mapper.ApplyUpdateToDomain(resource, req)
	if err = s.validator.ValidateResource(resource.Code, resource.Title, resource.Quantity); err != nil {
		return nil, err
	}
	updatedResource, err := s.repository.UpdateResource(ctx, *resource)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updatedResource), nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, id int64) error {
	return s.repository.DeleteResource(ctx, id)
}
