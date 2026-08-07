package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/dto"
)

type EmployeeService interface {
	ListEmployeesByResource(ctx context.Context, resourceID int64) ([]dto.EmployeeResponse, error)
	FindEmployee(ctx context.Context, id int64) (*dto.EmployeeResponse, error)
	CreateEmployee(ctx context.Context, resourceID int64, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error)
	UpdateEmployee(ctx context.Context, id int64, req dto.UpdateEmployeeRequest) (*dto.EmployeeResponse, error)
	DeleteEmployee(ctx context.Context, id int64) error

	ListStates(ctx context.Context, employeeID int64, start, end date.Date) ([]dto.EmployeeStateResponse, error)
	SetDays(ctx context.Context, employeeID int64, req dto.SetDaysRequest) error
	DeleteDays(ctx context.Context, employeeID int64, start, end date.Date, stateID *int64) error
}
