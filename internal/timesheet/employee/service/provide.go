package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/employee/repository"
)

// ProvideEmployeeService builds the EmployeeService service.
func ProvideEmployeeService(logger *slog.Logger, r *repo.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		logger:     logger,
		repository: r,
		mapper:     NewEmployeeMapper(),
		validator:  &EmployeeValidator{},
	}
}
