package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

type ProjectValidator struct {
	validator.Validator
}

func (v *ProjectValidator) ValidateProject(project *domain.Project) error {
	if err := v.ValidateRequiredText(project.Code, "code"); err != nil {
		return err
	}
	if err := v.ValidateOptionalColor(project.Color, "color"); err != nil {
		return err
	}
	if project.Priority < 0 {
		return errors.NewValidationError("priority must be positive")
	}
	if err := v.ValidateRequiredDate(project.StartDate, "start_date"); err != nil {
		return err
	}
	if err := v.ValidateRequiredDate(project.EndDate, "end_date"); err != nil {
		return err
	}
	return v.ValidateDateRange(project.StartDate, project.EndDate, "project")
}
