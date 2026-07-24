package domain

import (
	"fmt"
	"time"

	"github.com/Koshsky/erp/api/internal/validator"
)

type ProjectValidator struct {
	validator.Validator
}

func (v *ProjectValidator) ValidateProject(reqCode string, startDate, endDate time.Time, priority int) error {
	if err := v.ValidateRequiredText(reqCode, "code"); err != nil {
		return err
	}
	if priority < 0 {
		return fmt.Errorf("priority must be positive")
	}
	if err := v.ValidateRequiredDate(startDate, "start_date"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(endDate, "end_date"); err != nil {
		return err
	}
	return v.ValidateDateRange(startDate, endDate, "project")
}
