package service

import (
	"fmt"

	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type UserValidator struct {
	validator.Validator
}

func (v *UserValidator) ValidateUser(user *domain.User) error {
	if err := v.ValidateRequiredText(user.Name, "name"); err != nil {
		return err
	}
	if err := v.ValidateRequiredText(user.Username, "username"); err != nil {
		return err
	}
	switch user.Role {
	case domain.ProjectDirector, domain.ProjectManager, domain.ProcessOwner:
	default:
		return fmt.Errorf("unsupported role: %s", user.Role)
	}
	if user.ManagerID != nil {
		if err := v.ValidatePositiveID(*user.ManagerID, "manager_id"); err != nil {
			return err
		}
	}
	return nil
}
