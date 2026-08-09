package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/timesheet/employee/domain"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

// maxDayRange is the maximum calendar range width per request (in days).
const (
	maxDayRange = 730
	hoursPerDay = 24
)

type EmployeeValidator struct {
	validator.Validator
}

func (v *EmployeeValidator) ValidateEmployee(employee *domain.Employee) error {
	if err := v.ValidatePositiveID(employee.ResourceID, "resource_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(employee.Name, "name"); err != nil {
		return err
	}
	if employee.ManagerID != nil {
		if err := v.ValidatePositiveID(*employee.ManagerID, "manager_id"); err != nil {
			return err
		}
	}
	if employee.HireDate != nil && employee.TerminationDate != nil &&
		employee.TerminationDate.Before(*employee.HireDate) {
		return errors.NewFieldError(
			"termination_date",
			"date_range",
			"termination_date must be greater than or equal to hire_date",
		)
	}
	return nil
}

func (v *EmployeeValidator) ValidateDayRange(start, end date.Date) error {
	if start == "" {
		return errors.NewFieldError("start_date", "required", "start_date is required")
	}
	if end == "" {
		return errors.NewFieldError("end_date", "required", "end_date is required")
	}

	s, e := start.Time(), end.Time()
	if e.Before(s) {
		return errors.NewFieldError(
			"end_date",
			"date_range",
			"end_date must be greater than or equal to start_date",
		)
	}
	if int(e.Sub(s).Hours()/hoursPerDay) > maxDayRange-1 {
		return errors.NewValidationError(
			fmt.Sprintf("date range must not exceed %d days", maxDayRange),
		)
	}
	return nil
}
