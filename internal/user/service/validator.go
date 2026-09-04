package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

// maxDayRange is the maximum calendar range width per request (in days).
// maxPresetLen is the maximum preset name length (RBAC catalog).
const (
	maxDayRange  = 730
	maxPresetLen = 32
	hoursPerDay  = 24
)

// usernamePattern — the strict login rule: starts with a letter or digit, then
// letters/digits/./_, total length 3..20. Logins are always validated as
// usernames (email logins are not supported). Pattern assumes lowercase input.
const usernamePattern = "^[a-z0-9][a-z0-9._]{2,19}$"

// NormalizeUsername trims and lowercases a login (avoids User/user duplicates).
func NormalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

// IsUsernameReserved reports whether the login equals a reserved system word.
func IsUsernameReserved(username string) bool {
	for _, r := range [...]string{"admin", "support", "root", "system", "help"} {
		if username == r {
			return true
		}
	}
	return false
}

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
	// Logins are always validated as usernames (no email variant), normalized
	// to lowercase before storage so the case-insensitive login matches.
	user.Username = NormalizeUsername(user.Username)
	if err := v.ValidateRequiredText(user.Username, "username"); err != nil {
		return err
	}
	if err := v.ValidateUsername(user.Username); err != nil {
		return err
	}
	// Presets are configured at runtime by the rbac_presets catalog: here we
	// only check the form (NULL — no base preset is allowed); preset existence
	// is guaranteed by the FK users_preset_fk — violations surface as 400 via
	// mapUserErr.
	if user.Preset != nil {
		if len(*user.Preset) > maxPresetLen {
			return errors.NewFieldError("preset", "too_long", "preset is too long")
		}
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

// ValidateUsername checks the strict username rule (always as username, no
// email variant). Assumes the value is already trimmed and lowercased.
func (v *UserValidator) ValidateUsername(username string) error {
	ok, err := regexp.MatchString(usernamePattern, username)
	if err != nil || !ok {
		return errors.NewFieldError(
			"username",
			"format",
			"Логин: только латиница (a–z), цифры, точка и подчёркивание. Длина от 3 до 20 символов",
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
