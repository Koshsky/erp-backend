package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

// maxDayRange is the maximum calendar range width per request (in days).
// maxRoleLen is the maximum role name length (RBAC catalog).
const (
	maxDayRange = 730
	maxRoleLen  = 32
	hoursPerDay = 24
)

type UserValidator struct {
	validator.Validator
}

func (v *UserValidator) ValidateUser(user *domain.User) error {
	if err := v.ValidateRequiredText(user.LastName, "last_name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(user.FirstName, "first_name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(user.Username, "username"); err != nil {
		return err
	}
	// Roles are configured at runtime by the rbac_roles catalog (V15): here we
	// only check the form; role existence is guaranteed by the FK
	// users_role_fk (V17) — violations surface as 400 via mapUserErr.
	if user.Role == "" {
		return errors.NewFieldError("role", "required", "role is required")
	}
	if len(user.Role) > maxRoleLen {
		return errors.NewFieldError("role", "too_long", "role is too long")
	}
	if user.ManagerID != nil {
		if err := v.ValidatePositiveID(*user.ManagerID, "manager_id"); err != nil {
			return err
		}
	}
	if user.HireDate != nil && user.TerminationDate != nil &&
		user.TerminationDate.Before(*user.HireDate) {
		return errors.NewFieldError(
			"termination_date",
			"date_range",
			"termination_date must be greater than or equal to hire_date",
		)
	}
	return nil
}

func (v *UserValidator) ValidateDayRange(start, end date.Date) error {
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
