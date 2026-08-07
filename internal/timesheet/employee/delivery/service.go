package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
)

type EmployeeService interface {
	ListEmployeesByResource(ctx context.Context, resourceID int64) ([]dto.EmployeeResponse, error)
	ListEmployees(ctx context.Context, managerID *int64, userID int64, role string) ([]dto.EmployeeResponse, error)
	FindEmployee(ctx context.Context, id int64, userID int64, role string) (*dto.EmployeeResponse, error)
	CreateEmployee(
		ctx context.Context,
		resourceID int64,
		req dto.CreateEmployeeRequest,
		userID int64,
		role string,
	) (*dto.EmployeeResponse, error)
	UpdateEmployee(
		ctx context.Context,
		id int64,
		req dto.UpdateEmployeeRequest,
		userID int64,
		role string,
	) (*dto.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, id int64, userID int64, role string) error

	ListStates(
		ctx context.Context,
		employeeID int64,
		start, end date.Date,
		userID int64,
		role string,
	) ([]dto.EmployeeStateResponse, error)
	SetDays(ctx context.Context, employeeID int64, req dto.SetDaysRequest, userID int64, role string) error
	DeleteDays(
		ctx context.Context,
		employeeID int64,
		start, end date.Date,
		stateID *int64,
		userID int64,
		role string,
	) error
}
