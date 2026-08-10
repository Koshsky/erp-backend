package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/timesheet/employee/domain"
)

type EmployeeRepository interface {
	IsResourceActive(ctx context.Context, resourceID int64) (bool, error)
	ListEmployeesByResourceID(ctx context.Context, resourceID int64) ([]domain.Employee, error)
	ListEmployees(ctx context.Context, limit, offset int) ([]domain.Employee, error)
	CountEmployees(ctx context.Context) (int64, error)
	ListEmployeesByManagerID(ctx context.Context, managerID int64) ([]domain.Employee, error)
	FindEmployee(ctx context.Context, id int64) (*domain.Employee, error)
	CreateEmployee(ctx context.Context, employee domain.Employee) (*domain.Employee, error)
	UpdateEmployee(ctx context.Context, employee domain.Employee) (*domain.Employee, error)
	DeleteEmployee(ctx context.Context, id int64) error

	ListStates(ctx context.Context, employeeID int64, start, end time.Time) ([]domain.EmployeeState, error)
	SetStateRange(ctx context.Context, employeeID, stateID int64, start, end time.Time) error
	DeleteStateRange(ctx context.Context, employeeID int64, start, end time.Time, stateID *int64) error
}
