package domain

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/validator"
)

type AssignmentValidator struct {
	validator.Validator
}

func (v *AssignmentValidator) ValidateAssignment(taskID, resourceID int64, quantity int) error {
	if err := v.ValidatePositiveID(taskID, "task_id"); err != nil {
		return err
	}
	if err := v.ValidatePositiveID(resourceID, "resource_id"); err != nil {
		return err
	}
	if quantity < 1 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	return nil
}
