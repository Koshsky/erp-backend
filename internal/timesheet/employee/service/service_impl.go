package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
)

type EmployeeService struct {
	logger     *slog.Logger
	repository EmployeeRepository
	mapper     *EmployeeMapper
	validator  *EmployeeValidator
}

func NewEmployeeService(logger *slog.Logger, repository EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		logger:     logger,
		repository: repository,
		mapper:     NewEmployeeMapper(),
		validator:  &EmployeeValidator{},
	}
}

func (s *EmployeeService) ListEmployeesByResource(
	ctx context.Context,
	resourceID int64,
) ([]dto.EmployeeResponse, error) {
	employees, err := s.repository.ListEmployeesByResourceID(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToEmployeeDTOs(employees), nil
}

// ListEmployees возвращает всех сотрудников или только подчинённых менеджера.
func (s *EmployeeService) ListEmployees(
	ctx context.Context,
	managerID *int64,
) ([]dto.EmployeeResponse, error) {
	if managerID == nil {
		employees, err := s.repository.ListEmployees(ctx)
		if err != nil {
			return nil, err
		}
		return s.mapper.ToEmployeeDTOs(employees), nil
	}

	employees, err := s.repository.ListEmployeesByManagerID(ctx, *managerID)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToEmployeeDTOs(employees), nil
}

func (s *EmployeeService) FindEmployee(ctx context.Context, id int64) (*dto.EmployeeResponse, error) {
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}
	return s.mapper.ToEmployeeDTO(employee), nil
}

func (s *EmployeeService) CreateEmployee(
	ctx context.Context,
	resourceID int64,
	req dto.CreateEmployeeRequest,
) (*dto.EmployeeResponse, error) {
	employee := s.mapper.ToDomainFromCreate(req)
	employee.ResourceID = resourceID
	if err := s.validator.ValidateEmployee(&employee); err != nil {
		return nil, err
	}

	active, err := s.repository.IsResourceActive(ctx, resourceID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, fmt.Errorf("resource not found or deleted")
	}

	created, err := s.repository.CreateEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToEmployeeDTO(created), nil
}

func (s *EmployeeService) UpdateEmployee(
	ctx context.Context,
	id int64,
	req dto.UpdateEmployeeRequest,
) (*dto.EmployeeResponse, error) {
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	s.mapper.ApplyUpdateToDomain(employee, req)
	if err = s.validator.ValidateEmployee(employee); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateEmployee(ctx, *employee)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToEmployeeDTO(updated), nil
}

func (s *EmployeeService) DeleteEmployee(ctx context.Context, id int64) error {
	return s.repository.DeleteEmployee(ctx, id)
}

func (s *EmployeeService) ListStates(
	ctx context.Context,
	employeeID int64,
	start, end date.Date,
) ([]dto.EmployeeStateResponse, error) {
	if err := s.validator.ValidatePositiveID(employeeID, "employee_id"); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateDayRange(start, end); err != nil {
		return nil, err
	}

	employee, err := s.repository.FindEmployee(ctx, employeeID)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	states, err := s.repository.ListStates(ctx, employeeID, start.Time(), end.Time())
	if err != nil {
		return nil, err
	}
	return s.mapper.ToStateDTOs(states), nil
}

func (s *EmployeeService) SetDays(ctx context.Context, employeeID int64, req dto.SetDaysRequest) error {
	if err := s.validator.ValidatePositiveID(employeeID, "employee_id"); err != nil {
		return err
	}
	if err := s.validator.ValidatePositiveID(req.StateID, "state_id"); err != nil {
		return err
	}
	if err := s.validator.ValidateDayRange(req.StartDate, req.EndDate); err != nil {
		return err
	}

	return s.repository.SetStateRange(ctx, employeeID, req.StateID, req.StartDate.Time(), req.EndDate.Time())
}

func (s *EmployeeService) DeleteDays(
	ctx context.Context,
	employeeID int64,
	start, end date.Date,
	stateID *int64,
) error {
	if err := s.validator.ValidatePositiveID(employeeID, "employee_id"); err != nil {
		return err
	}
	if stateID != nil {
		if err := s.validator.ValidatePositiveID(*stateID, "state_id"); err != nil {
			return err
		}
	}
	if err := s.validator.ValidateDayRange(start, end); err != nil {
		return err
	}

	return s.repository.DeleteStateRange(ctx, employeeID, start.Time(), end.Time(), stateID)
}
