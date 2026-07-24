package domain

import (
	"time"

	"github.com/Koshsky/erp/api/internal/validator"
)

type MilestoneValidator struct {
	validator.Validator
}

func (v *MilestoneValidator) ValidateMilestone(processID int64, title, content string, date time.Time) error {
	if err := v.ValidatePositiveID(processID, "process_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(title, "title"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(content, "content"); err != nil {
		return err
	}
	return v.ValidateRequiredDate(date, "date")
}
