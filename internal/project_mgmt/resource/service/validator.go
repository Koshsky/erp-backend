package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/resource/domain"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type ResourceValidator struct {
	validator.Validator
}

func (v *ResourceValidator) ValidateResource(resource *domain.Resource) error {
	if err := v.ValidateRequiredText(resource.Code, "code"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(resource.Title, "title"); err != nil {
		return err
	}
	if resource.Quantity < 0 {
		return fmt.Errorf("quantity must be positive")
	}
	return nil
}
