package validator

import (
	"strings"
	"time"
)

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateDateRange(start, end time.Time, entity string) error {
	if end.Before(start) {
		return NewFieldError("end_date", codeDateRange, msgDateRange(entity))
	}
	return nil
}

func (v *Validator) ValidateRequiredText(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return NewFieldError(field, codeRequired, msgRequired(field))
	}
	return nil
}

func (v *Validator) ValidatePositiveID(id int64, field string) error {
	if id <= 0 {
		return NewFieldError(field, codeMinValue, msgGreaterThan(field, 0))
	}
	return nil
}

func (v *Validator) ValidateRequiredDate(value time.Time, field string) error {
	if value.IsZero() {
		return NewFieldError(field, codeRequired, msgRequired(field))
	}
	return nil
}
