package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/domain"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

type MilestoneValidator struct {
	validator.Validator
}

func (v *MilestoneValidator) ValidateMilestone(milestone *domain.Milestone) error {
	if err := v.ValidatePositiveID(milestone.ProcessID, "process_id"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(milestone.Title, "title"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(milestone.Content, "content"); err != nil {
		return err
	}
	return v.ValidateRequiredDate(milestone.Date, "date")
}
