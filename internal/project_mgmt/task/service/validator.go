package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

type TaskValidator struct {
	validator.Validator
}

func (v *TaskValidator) ValidateTask(task *domain.Task) error {
	if err := v.ValidatePositiveID(task.ProcessID, "process_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(task.Title, "title"); err != nil {
		return err
	}
	if err := v.ValidateOptionalColor(task.Color, "color"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(task.StartDate, "start_date"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(task.EndDate, "end_date"); err != nil {
		return err
	}
	return v.ValidateDateRange(task.StartDate, task.EndDate, "task")
}
