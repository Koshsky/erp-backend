package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

type ProcessValidator struct {
	validator.Validator
}

func (v *ProcessValidator) ValidateProcess(process *domain.Process) error {
	if err := v.ValidatePositiveID(process.ProjectID, "project_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(process.Title, "title"); err != nil {
		return err
	}
	if err := v.ValidateOptionalColor(process.Color, "color"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(process.StartDate, "start_date"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(process.EndDate, "end_date"); err != nil {
		return err
	}
	return v.ValidateDateRange(process.StartDate, process.EndDate, "process")
}
