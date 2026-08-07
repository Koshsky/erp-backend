package service

import (
	"github.com/Koshsky/erp-backend/internal/timesheet/state/domain"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type StateValidator struct {
	validator.Validator
}

func (v *StateValidator) ValidateState(state *domain.State) error {
	if err := v.ValidateRequiredText(state.Code, "code"); err != nil {
		return err
	}
	return v.ValidateRequiredText(state.Name, "name")
}
