package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/timesheet/employee/service"
)

// ProvideEmployeeHandler builds the EmployeeHandler handler.
func ProvideEmployeeHandler(logger *slog.Logger, svc *service.EmployeeService, mw *rbac.Middleware) *EmployeeHandler {
	return &EmployeeHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
