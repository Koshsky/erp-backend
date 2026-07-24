package domain

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/validator"
)

type TaskValidator struct {
	validator.Validator
}

func (v *TaskValidator) ValidateTask(processID int64, title string, startDate, endDate time.Time) error {
	if err := v.ValidatePositiveID(processID, "process_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(title, "title"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(startDate, "start_date"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(endDate, "end_date"); err != nil {
		return err
	}
	return v.ValidateDateRange(startDate, endDate, "task")
}
