package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type UserValidator struct {
	validator.Validator
}

func (v *UserValidator) ValidateUserCreate(name, username, role string) error {
	if err := v.ValidateRequiredText(name, "name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(username, "username"); err != nil {
		return err
	}
	if role != domain.ProjectDirector && role != domain.ProjectManager && role != domain.ProcessOwner {
		return fmt.Errorf("unsupported role")
	}
	return nil
}

func (v *UserValidator) ValidateUserUpdate(name, username, role string) error {
	if err := v.ValidateRequiredText(name, "name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(username, "username"); err != nil {
		return err
	}
	if role != domain.ProjectDirector && role != domain.ProjectManager && role != domain.ProcessOwner {
		return fmt.Errorf("unsupported role")
	}
	return nil
}
