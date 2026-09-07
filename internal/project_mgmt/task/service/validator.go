package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
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
	if err := v.validateStatus(task.Status); err != nil {
		return err
	}
	if err := v.ValidateOptionalPositiveID(task.ParentID, "parent_id"); err != nil {
		return err
	}
	// Subtasks inherit the parent's dates, top-level tasks require their own.
	if task.ParentID == nil {
		if err := v.ValidateRequiredDate(task.StartDate, "start_date"); err != nil {
			return err
		}
		if err := v.ValidateRequiredDate(task.EndDate, "end_date"); err != nil {
			return err
		}
	}
	return v.ValidateDateRange(task.StartDate, task.EndDate, "task")
}

// validateStatus rejects values outside the fixed 3-state catalog.
func (v *TaskValidator) validateStatus(status string) error {
	switch status {
	case domain.StatusNotStarted, domain.StatusInProgress, domain.StatusDone:
		return nil
	default:
		return errors.NewFieldError("status", "one_of", "недопустимый статус задачи")
	}
}
