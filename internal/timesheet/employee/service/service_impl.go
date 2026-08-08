package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/validator"
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

// ListEmployees returns all employees; visibility by manager is filtered by the
// middleware (vp sees only their own subordinates).
func (s *EmployeeService) ListEmployees(ctx context.Context) ([]dto.EmployeeResponse, error) {
	employees, err := s.repository.ListEmployees(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToEmployeeDTOs(employees), nil
}

func (s *EmployeeService) FindEmployee(ctx context.Context, id int64) (*dto.EmployeeResponse, error) {
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		if validator.IsNotFoundError(err) {
			return nil, validator.ErrEmployeeNotFound
		}
		return nil, err
	}
	if employee == nil {
		return nil, validator.ErrEmployeeNotFound
	}
	return s.mapper.ToEmployeeDTO(employee), nil
}

// CreateEmployee creates an employee. The middleware checked permissions; here
// only normalization: vp creates an employee into their own team.
func (s *EmployeeService) CreateEmployee(
	ctx context.Context,
	resourceID int64,
	req dto.CreateEmployeeRequest,
	userID int64,
	role string,
) (*dto.EmployeeResponse, error) {
	if role == userdomain.ProcessOwner {
		// vp creates an employee into their own team: a foreign manager_id is ignored
		req.ManagerID = &userID
	}

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
		return nil, validator.ErrResourceNotFound
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
		if validator.IsNotFoundError(err) {
			return nil, validator.ErrEmployeeNotFound
		}
		return nil, err
	}
	if employee == nil {
		return nil, validator.ErrEmployeeNotFound
	}

	s.mapper.ApplyUpdateToDomain(employee, req)
	// The new position (resource) must exist and not be deleted.
	if req.ResourceID != nil {
		var active bool
		active, err = s.repository.IsResourceActive(ctx, *req.ResourceID)
		if err != nil {
			return nil, err
		}
		if !active {
			return nil, validator.ErrResourceNotFound
		}
	}
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
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		if validator.IsNotFoundError(err) {
			return validator.ErrEmployeeNotFound
		}
		return err
	}
	if employee == nil {
		return validator.ErrEmployeeNotFound
	}

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

	if err := s.ensureEmployeeExists(ctx, employeeID); err != nil {
		return nil, err
	}

	states, err := s.repository.ListStates(ctx, employeeID, start.Time(), end.Time())
	if err != nil {
		return nil, err
	}
	return s.mapper.ToStateDTOs(states), nil
}

func (s *EmployeeService) SetDays(
	ctx context.Context,
	employeeID int64,
	req dto.SetDaysRequest,
) error {
	if err := s.validator.ValidatePositiveID(employeeID, "employee_id"); err != nil {
		return err
	}
	if err := s.validator.ValidatePositiveID(req.StateID, "state_id"); err != nil {
		return err
	}
	if err := s.validator.ValidateDayRange(req.StartDate, req.EndDate); err != nil {
		return err
	}

	if err := s.ensureEmployeeExists(ctx, employeeID); err != nil {
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

	if err := s.ensureEmployeeExists(ctx, employeeID); err != nil {
		return err
	}

	return s.repository.DeleteStateRange(ctx, employeeID, start.Time(), end.Time(), stateID)
}

// ensureEmployeeExists verifies the employee exists (404 otherwise).
func (s *EmployeeService) ensureEmployeeExists(ctx context.Context, employeeID int64) error {
	employee, err := s.repository.FindEmployee(ctx, employeeID)
	if err != nil {
		if validator.IsNotFoundError(err) {
			return validator.ErrEmployeeNotFound
		}
		return err
	}
	if employee == nil {
		return validator.ErrEmployeeNotFound
	}
	return nil
}
