package service

import (
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/domain"
	"github.com/Koshsky/erp-backend/pkg/validator"
)

type ResourceValidator struct {
	validator.Validator
}

func (v *ResourceValidator) ValidateResource(resource *domain.Resource) error {
	if err := v.ValidateRequiredText(resource.Code, "code"); err != nil {
		return err
	}
	return v.ValidateRequiredText(resource.Title, "title")
}
