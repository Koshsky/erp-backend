package domain

import (
	"fmt"
	"strings"

	"github.com/Koshsky/erp-backend/internal/validator"
)

type UserValidator struct {
	validator.Validator
}

func (v *UserValidator) ValidateUserCreate(name, username string, role UserRole, password string) error {
	if err := v.ValidateRequiredText(name, "name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(username, "username"); err != nil {
		return err
	}
	if role != UserRoleProjectDirector && role != UserRoleProjectManager && role != UserRoleProcessOwner {
		return fmt.Errorf("unsupported role")
	}
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (v *UserValidator) ValidateUserUpdate(name, username string, role UserRole) error {
	if err := v.ValidateRequiredText(name, "name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(username, "username"); err != nil {
		return err
	}
	if role != UserRoleProjectDirector && role != UserRoleProjectManager && role != UserRoleProcessOwner {
		return fmt.Errorf("unsupported role")
	}
	return nil
}
