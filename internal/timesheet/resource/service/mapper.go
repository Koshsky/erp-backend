package service

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type ResourceMapper struct{}

func NewResourceMapper() *ResourceMapper {
	return &ResourceMapper{}
}

func (m *ResourceMapper) ToDTO(resource *domain.Resource) *dto.ResourceResponse {
	if resource == nil {
		return nil
	}
	return &dto.ResourceResponse{
		ID:             resource.ID,
		Code:           resource.Code,
		Title:          resource.Title,
		Color:          resource.Color,
		OwnerID:        resource.OwnerID,
		EmployeesCount: resource.EmployeesCount,
	}
}

func (m *ResourceMapper) ToDTOs(resources []domain.Resource) []dto.ResourceResponse {
	if resources == nil {
		return []dto.ResourceResponse{}
	}

	responses := make([]dto.ResourceResponse, len(resources))
	for i, resource := range resources {
		responses[i] = *m.ToDTO(&resource)
	}
	return responses
}

func (m *ResourceMapper) ToDomainFromCreate(req dto.CreateResourceRequest) domain.Resource {
	return domain.Resource{
		Code:    req.Code,
		Title:   req.Title,
		Color:   req.Color,
		OwnerID: req.OwnerID,
	}
}

func (m *ResourceMapper) ApplyUpdateToDomain(resource *domain.Resource, req dto.UpdateResourceRequest) {
	if resource == nil {
		return
	}

	if req.Code != nil {
		resource.Code = *req.Code
	}
	if req.Title != nil {
		resource.Title = *req.Title
	}
	if req.Color != nil {
		if *req.Color == "" {
			resource.Color = nil
		} else {
			resource.Color = req.Color
		}
	}
	if req.OwnerID != nil {
		resource.OwnerID = req.OwnerID
	}
}

func (m *ResourceMapper) ToMemberDTO(member *domain.ResourceMember) *dto.ResourceMemberResponse {
	if member == nil {
		return nil
	}
	return &dto.ResourceMemberResponse{
		ID:              member.ID,
		Name:            member.Name,
		Preset:          member.Preset,
		Position:        member.Position,
		ManagerID:       member.ManagerID,
		HireDate:        datePtr(member.HireDate),
		TerminationDate: datePtr(member.TerminationDate),
	}
}

func (m *ResourceMapper) ToMemberDTOs(members []domain.ResourceMember) []dto.ResourceMemberResponse {
	if members == nil {
		return []dto.ResourceMemberResponse{}
	}
	responses := make([]dto.ResourceMemberResponse, len(members))
	for i, member := range members {
		responses[i] = *m.ToMemberDTO(&member)
	}
	return responses
}

func (m *ResourceMapper) ToAbsenceDTOs(absences []domain.ResourceAbsence) []dto.ResourceAbsenceResponse {
	if absences == nil {
		return []dto.ResourceAbsenceResponse{}
	}
	responses := make([]dto.ResourceAbsenceResponse, len(absences))
	for i, a := range absences {
		responses[i] = dto.ResourceAbsenceResponse{
			UserID:    a.UserID,
			UserName:  a.UserName,
			StateID:   a.StateID,
			StateCode: a.StateCode,
			StateName: a.StateName,
			StartDate: date.From(a.StartDate),
			EndDate:   date.From(a.EndDate),
		}
	}
	return responses
}

func datePtr(t *time.Time) *date.Date {
	if t == nil {
		return nil
	}
	d := date.From(*t)
	return &d
}
