package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type ResourceRepository interface {
	CreateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error)
	GetResource(ctx context.Context, id int64) (*domain.Resource, error)
	UpdateResource(ctx context.Context, resource domain.Resource) (*domain.Resource, error)
	DeleteResource(ctx context.Context, id int64) error
	GetResourceUsage(ctx context.Context, targetDate time.Time) ([]domain.ResourceUsage, error)
	ListResources(ctx context.Context) ([]domain.Resource, error)
}

type ResourceService struct {
	logger     *slog.Logger
	repository ResourceRepository
	mapper     *mapper.ResourceMapper
	validator  *Validator
}

func NewResourceService(logger *slog.Logger, repository ResourceRepository, validator *Validator) *ResourceService {
	return &ResourceService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewResourceMapper(),
		validator:  validator,
	}
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

func (s *ResourceService) GetResource(ctx context.Context, id int64) (*dto.ResourceResponse, error) {
	resource, err := s.repository.GetResource(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, ErrResourceNotFound
	}
	return s.mapper.ToDTO(resource), nil
}

func (s *ResourceService) UpdateResource(ctx context.Context, id int64, req dto.UpdateResourceRequest) (*dto.ResourceResponse, error) {
	resource, err := s.repository.GetResource(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, ErrResourceNotFound
	}
	s.mapper.ApplyUpdateToDomain(resource, req)
	if err := s.validator.ValidateResource(resource.Code, resource.Title, resource.Quantity); err != nil {
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

func (s *ResourceService) GetResourceUsage(ctx context.Context, targetDate time.Time) ([]dto.ResourceUsageResponse, error) {
	domainUsages, err := s.repository.GetResourceUsage(ctx, targetDate)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToUsageDTOs(domainUsages), nil
}

func (s *ResourceService) ListResources(ctx context.Context) ([]dto.ResourceResponse, error) {
	domainResources, err := s.repository.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(domainResources), nil
}
