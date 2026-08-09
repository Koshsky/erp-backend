package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

type AssignmentValidator struct {
	validator.Validator
}

func (v *AssignmentValidator) ValidateAssignment(assignment *domain.Assignment) error {
	if err := v.ValidatePositiveID(assignment.TaskID, "task_id"); err != nil {
		return err
	}
	if err := v.ValidatePositiveID(assignment.ResourceID, "resource_id"); err != nil {
		return err
	}
	if assignment.Quantity < 1 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	return nil
}
