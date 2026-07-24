package domain

import (
	"github.com/Koshsky/erp-backend/internal/resource/dto"
)

type ResourceMapper struct{}

func NewResourceMapper() *ResourceMapper {
	return &ResourceMapper{}
}

func (m *ResourceMapper) ToDTO(resource *Resource) *dto.ResourceResponse {
	if resource == nil {
		return nil
	}
	return &dto.ResourceResponse{
		ID:       resource.ID,
		Code:     resource.Code,
		Title:    resource.Title,
		Quantity: resource.Quantity,
	}
}

func (m *ResourceMapper) ToDTOs(resources []Resource) []dto.ResourceResponse {
	if resources == nil {
		return []dto.ResourceResponse{}
	}

	responses := make([]dto.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = *m.ToDTO(&resource)
	}
	return responses
}

func (m *ResourceMapper) ToDomainFromCreate(req dto.CreateResourceRequest) Resource {
	return Resource{
		Code:     req.Code,
		Title:    req.Title,
		Quantity: req.Quantity,
	}
}

func (m *ResourceMapper) ApplyUpdateToDomain(resource *Resource, req dto.UpdateResourceRequest) {
	if resource == nil {
		return
	}

	if req.Code != nil {
		resource.Code = *req.Code
	}
	if req.Title != nil {
		resource.Title = *req.Title
	}
	if req.Quantity != nil {
		resource.Quantity = *req.Quantity
	}
}
