package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
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

// canManageAllEmployees возвращает true, если роль управляет всеми сотрудниками.
func canManageAllEmployees(role string) bool {
	return role == userdomain.Admin
}

// managesEmployee проверяет, что пользователь управляет сотрудником:
// admin — любой, остальные — только подчинённые (manager_id = user.id).
func managesEmployee(role string, userID int64, employee *domain.Employee) bool {
	if canManageAllEmployees(role) {
		return true
	}
	return employee.ManagerID != nil && *employee.ManagerID == userID
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

// ListEmployees возвращает всех сотрудников (admin) или только подчинённых
// текущего пользователя. Для не-admin переданный manager_id игнорируется,
// чтобы нельзя было увидеть чужую команду.
func (s *EmployeeService) ListEmployees(
	ctx context.Context,
	managerID *int64,
	userID int64,
	role string,
) ([]dto.EmployeeResponse, error) {
	if !canManageAllEmployees(role) {
		managerID = &userID
	}

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

func (s *EmployeeService) FindEmployee(
	ctx context.Context,
	id int64,
	userID int64,
	role string,
) (*dto.EmployeeResponse, error) {
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrNotFound
	}
	if !managesEmployee(role, userID, employee) {
		return nil, ErrForbidden
	}
	return s.mapper.ToEmployeeDTO(employee), nil
}

func (s *EmployeeService) CreateEmployee(
	ctx context.Context,
	resourceID int64,
	req dto.CreateEmployeeRequest,
	userID int64,
	role string,
) (*dto.EmployeeResponse, error) {
	switch {
	case canManageAllEmployees(role):
		// admin указывает любого руководителя (или оставляет без руководителя)
	case role == userdomain.ProcessOwner:
		// vp создаёт сотрудника себе в подчинение: чужой manager_id из запроса игнорируется
		req.ManagerID = &userID
	default:
		return nil, ErrForbidden
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
		return nil, ErrResourceNotFound
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
	userID int64,
	role string,
) (*dto.EmployeeResponse, error) {
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrNotFound
	}
	if !managesEmployee(role, userID, employee) {
		return nil, ErrForbidden
	}

	s.mapper.ApplyUpdateToDomain(employee, req)
	// Не-admin не может переподчинить своего сотрудника другому менеджеру.
	if !canManageAllEmployees(role) {
		employee.ManagerID = &userID
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

func (s *EmployeeService) DeleteEmployee(
	ctx context.Context,
	id int64,
	userID int64,
	role string,
) error {
	employee, err := s.repository.FindEmployee(ctx, id)
	if err != nil {
		return err
	}
	if employee == nil {
		return ErrNotFound
	}
	if !managesEmployee(role, userID, employee) {
		return ErrForbidden
	}

	return s.repository.DeleteEmployee(ctx, id)
}

func (s *EmployeeService) ListStates(
	ctx context.Context,
	employeeID int64,
	start, end date.Date,
	userID int64,
	role string,
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
		return nil, ErrNotFound
	}
	if !managesEmployee(role, userID, employee) {
		return nil, ErrForbidden
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
	userID int64,
	role string,
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

	if !s.managesEmployeeID(ctx, employeeID, userID, role) {
		return ErrForbidden
	}

	return s.repository.SetStateRange(ctx, employeeID, req.StateID, req.StartDate.Time(), req.EndDate.Time())
}

func (s *EmployeeService) DeleteDays(
	ctx context.Context,
	employeeID int64,
	start, end date.Date,
	stateID *int64,
	userID int64,
	role string,
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

	if !s.managesEmployeeID(ctx, employeeID, userID, role) {
		return ErrForbidden
	}

	return s.repository.DeleteStateRange(ctx, employeeID, start.Time(), end.Time(), stateID)
}

// managesEmployeeID загружает сотрудника и проверяет право управления им.
func (s *EmployeeService) managesEmployeeID(
	ctx context.Context,
	employeeID int64,
	userID int64,
	role string,
) bool {
	if canManageAllEmployees(role) {
		return true
	}
	employee, err := s.repository.FindEmployee(ctx, employeeID)
	if err != nil || employee == nil {
		return false
	}
	return employee.ManagerID != nil && *employee.ManagerID == userID
}
