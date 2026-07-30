package validator

import "fmt"

// Machine-readable codes for FieldError. Centralized here so message wording
// and codes evolve in one place (e.g. for future localization or structured
// API error payloads) instead of being scattered across Validator methods.
const (
	codeRequired  = "required"
	codeMinValue  = "min_value"
	codeDateRange = "date_range"
	codeOneOf     = "one_of"
)

func msgRequired(field string) string {
	return fmt.Sprintf("%s is required", field)
}

func msgGreaterThan(field string, min int) string {
	return fmt.Sprintf("%s must be greater than %d", field, min)
}

func msgGreaterThanOrEqual(field string, min int) string {
	return fmt.Sprintf("%s must be greater than or equal to %d", field, min)
}

func msgAtLeast(field string, min int) string {
	return fmt.Sprintf("%s must be at least %d", field, min)
}

func msgDateRange(entity string) string {
	return fmt.Sprintf("%s end_date must be greater than or equal to start_date", entity)
}
