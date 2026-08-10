package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type EmployeeService interface {
	ListEmployeesByResource(ctx context.Context, resourceID int64) ([]dto.EmployeeResponse, error)
	ListEmployees(
		ctx context.Context,
		userID int64,
		role string,
		ownerID int64,
		limit, offset int,
	) ([]dto.EmployeeResponse, int64, error)
	FindEmployee(ctx context.Context, id int64) (*dto.EmployeeResponse, error)
	CreateEmployee(
		ctx context.Context,
		resourceID int64,
		req dto.CreateEmployeeRequest,
		userID int64,
		role string,
	) (*dto.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, id int64, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, id int64) error

	ListStates(
		ctx context.Context,
		employeeID int64,
		start, end date.Date,
	) ([]dto.EmployeeStateResponse, error)
	SetDays(ctx context.Context, employeeID int64, req dto.SetDaysRequest) error
	DeleteDays(
		ctx context.Context,
		employeeID int64,
		start, end date.Date,
		stateID *int64,
	) error
}
