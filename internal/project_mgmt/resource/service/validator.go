package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/validator"
)

type ResourceValidator struct {
	validator.Validator
}

func (v *ResourceValidator) ValidateResource(code, title string, quantity int) error {
	if err := v.ValidateRequiredText(code, "code"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(title, "title"); err != nil {
		return err
	}
	if quantity < 0 {
		return fmt.Errorf("quantity must be positive")
	}
	return nil
}
