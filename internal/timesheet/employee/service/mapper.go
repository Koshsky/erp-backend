package service

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
)

type EmployeeMapper struct{}

func NewEmployeeMapper() *EmployeeMapper {
	return &EmployeeMapper{}
}

func (m *EmployeeMapper) ToEmployeeDTO(employee *domain.Employee) *dto.EmployeeResponse {
	if employee == nil {
		return nil
	}
	return &dto.EmployeeResponse{
		ID:              employee.ID,
		ResourceID:      employee.ResourceID,
		ResourceTitle:   employee.ResourceTitle,
		Name:            employee.Name,
		Position:        employee.Position,
		ManagerID:       employee.ManagerID,
		HireDate:        datePtr(employee.HireDate),
		TerminationDate: datePtr(employee.TerminationDate),
	}
}

func (m *EmployeeMapper) ToEmployeeDTOs(employees []domain.Employee) []dto.EmployeeResponse {
	if employees == nil {
		return []dto.EmployeeResponse{}
	}

	responses := make([]dto.EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = *m.ToEmployeeDTO(&employee)
	}
	return responses
}

func (m *EmployeeMapper) ToDomainFromCreate(req dto.CreateEmployeeRequest) domain.Employee {
	return domain.Employee{
		Name:            req.Name,
		Position:        req.Position,
		ManagerID:       req.ManagerID,
		HireDate:        timePtr(req.HireDate),
		TerminationDate: timePtr(req.TerminationDate),
	}
}

func (m *EmployeeMapper) ApplyUpdateToDomain(employee *domain.Employee, req dto.UpdateEmployeeRequest) {
	if employee == nil {
		return
	}

	if req.ResourceID != nil {
		employee.ResourceID = *req.ResourceID
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.Name != nil {
		employee.Name = *req.Name
	}
	if req.ManagerID != nil {
		employee.ManagerID = req.ManagerID
	}
	if req.HireDate != nil {
		t := req.HireDate.Time()
		employee.HireDate = &t
	}
	if req.TerminationDate != nil {
		t := req.TerminationDate.Time()
		employee.TerminationDate = &t
	}
}

func (m *EmployeeMapper) ToStateDTOs(states []domain.EmployeeState) []dto.EmployeeStateResponse {
	if states == nil {
		return []dto.EmployeeStateResponse{}
	}

	responses := make([]dto.EmployeeStateResponse, len(states))
	for i, state := range states {
		responses[i] = dto.EmployeeStateResponse{
			ID:          state.ID,
			StateID:     state.StateID,
			StateCode:   state.StateCode,
			StateName:   state.StateName,
			IsAvailable: state.IsAvailable,
			StartDate:   date.From(state.StartDate),
			EndDate:     date.From(state.EndDate),
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

func timePtr(d *date.Date) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time()
	return &t
}
