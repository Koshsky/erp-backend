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
	codeFormat    = "format"
)

func msgRequired(field string) string {
	return fmt.Sprintf("%s is required", field)
}

func msgGreaterThan(field string, minVal int) string {
	return fmt.Sprintf("%s must be greater than %d", field, minVal)
}

func msgDateRange(entity string) string {
	return fmt.Sprintf("%s end_date must be greater than or equal to start_date", entity)
}

func msgFormat(field string) string {
	return fmt.Sprintf("%s must be a #RRGGBB hex color", field)
}
