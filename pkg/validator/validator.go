package validator

import (
	"strings"
	"time"

	"github.com/Koshsky/erp-backend/pkg/errors"
)

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateDateRange(start, end time.Time, entity string) error {
	if end.Before(start) {
		return errors.NewFieldError("end_date", codeDateRange, msgDateRange(entity))
	}
	return nil
}

func (v *Validator) ValidateRequiredText(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return errors.NewFieldError(field, codeRequired, msgRequired(field))
	}
	return nil
}

func (v *Validator) ValidatePositiveID(id int64, field string) error {
	if id <= 0 {
		return errors.NewFieldError(field, codeMinValue, msgGreaterThan(field, 0))
	}
	return nil
}

func (v *Validator) ValidateRequiredDate(value time.Time, field string) error {
	if value.IsZero() {
		return errors.NewFieldError(field, codeRequired, msgRequired(field))
	}
	return nil
}
