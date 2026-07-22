package service

import (
	"strings"
	"time"

	"github.com/Koshsky/erp/api/internal/domain"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) validateDateRange(start, end time.Time, entity string) error {
	if end.Before(start) {
		return newFieldError("end_date", codeDateRange, msgDateRange(entity))
	}
	return nil
}

func (v *Validator) validateRequiredText(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return newFieldError(field, codeRequired, msgRequired(field))
	}
	return nil
}

func (v *Validator) validatePositiveID(id int64, field string) error {
	if id <= 0 {
		return newFieldError(field, codeMinValue, msgGreaterThan(field, 0))
	}
	return nil
}

func (v *Validator) validateRequiredDate(value time.Time, field string) error {
	if value.IsZero() {
		return newFieldError(field, codeRequired, msgRequired(field))
	}
	return nil
}

func (v *Validator) ValidateProject(reqCode string, startDate, endDate time.Time, priority int) error {
	if err := v.validateRequiredText(reqCode, "code"); err != nil {
		return err
	}
	if priority < 0 {
		return newFieldError("priority", codeMinValue, msgGreaterThanOrEqual("priority", 0))
	}
	if err := v.validateRequiredDate(startDate, "start_date"); err != nil {
		return err
	}
	if err := v.validateRequiredDate(endDate, "end_date"); err != nil {
		return err
	}
	return v.validateDateRange(startDate, endDate, "project")
}

func (v *Validator) ValidateProcess(projectID int64, title string, startDate, endDate time.Time) error {
	if err := v.validatePositiveID(projectID, "project_id"); err != nil {
		return err
	}
	if err := v.validateRequiredText(title, "title"); err != nil {
		return err
	}
	if err := v.validateRequiredDate(startDate, "start_date"); err != nil {
		return err
	}
	if err := v.validateRequiredDate(endDate, "end_date"); err != nil {
		return err
	}
	return v.validateDateRange(startDate, endDate, "process")
}

func (v *Validator) ValidateTask(processID int64, title string, startDate, endDate time.Time) error {
	if err := v.validatePositiveID(processID, "process_id"); err != nil {
		return err
	}
	if err := v.validateRequiredText(title, "title"); err != nil {
		return err
	}
	if err := v.validateRequiredDate(startDate, "start_date"); err != nil {
		return err
	}
	if err := v.validateRequiredDate(endDate, "end_date"); err != nil {
		return err
	}
	return v.validateDateRange(startDate, endDate, "task")
}

func (v *Validator) ValidateResource(code, title string, quantity int) error {
	if err := v.validateRequiredText(code, "code"); err != nil {
		return err
	}
	if err := v.validateRequiredText(title, "title"); err != nil {
		return err
	}
	if quantity < 0 {
		return newFieldError("quantity", codeMinValue, msgGreaterThanOrEqual("quantity", 0))
	}
	return nil
}

func (v *Validator) ValidateMilestone(processID int64, title, content string, date time.Time) error {
	if err := v.validatePositiveID(processID, "process_id"); err != nil {
		return err
	}
	if err := v.validateRequiredText(title, "title"); err != nil {
		return err
	}
	if err := v.validateRequiredText(content, "content"); err != nil {
		return err
	}
	return v.validateRequiredDate(date, "date")
}

func (v *Validator) ValidateAssignment(taskID, resourceID int64, quantity int) error {
	if err := v.validatePositiveID(taskID, "task_id"); err != nil {
		return err
	}
	if err := v.validatePositiveID(resourceID, "resource_id"); err != nil {
		return err
	}
	if quantity < 1 {
		return ErrInvalidAssignmentQuantity
	}
	return nil
}

func (v *Validator) ValidateUserCreate(name, username string, role domain.UserRole, password string) error {
	if err := v.validateRequiredText(name, "name"); err != nil {
		return err
	}
	if err := v.validateRequiredText(username, "username"); err != nil {
		return err
	}
	if role != domain.UserRoleProjectDirector && role != domain.UserRoleProjectManager && role != domain.UserRoleProcessOwner {
		return newFieldError("role", codeOneOf, msgOneOf("role", "ДП, РП, ВП"))
	}
	if strings.TrimSpace(password) == "" {
		return newFieldError("password", codeRequired, msgRequired("password"))
	}
	return nil
}

func (v *Validator) ValidateUserUpdate(name, username string, role domain.UserRole) error {
	if err := v.validateRequiredText(name, "name"); err != nil {
		return err
	}
	if err := v.validateRequiredText(username, "username"); err != nil {
		return err
	}
	if role != domain.UserRoleProjectDirector && role != domain.UserRoleProjectManager && role != domain.UserRoleProcessOwner {
		return newFieldError("role", codeOneOf, msgOneOf("role", "ДП, РП, ВП"))
	}
	return nil
}
