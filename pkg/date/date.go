// Package date — a calendar date in YYYY-MM-DD format (without time or timezone).
//
// Single API contract for all "calendar" fields (start_date, end_date, date):
// in JSON they encode as a YYYY-MM-DD string (OpenAPI format: date), which
// matches the DATE column type in the database and the calendar grid on the
// frontend. [time.Time] for these fields would use RFC3339 with time and zone —
// a source of format mismatches.
package date

import (
	"fmt"
	"strings"
	"time"
)

// Layout — the only date format used in the API.
const Layout = "2006-01-02"

// Date — a calendar date without time, stored as a YYYY-MM-DD string.
// The string representation lets swag (OpenAPI) treat it as a string.
type Date string //nolint:recvcheck // json.Unmarshaler needs a pointer, the other methods are on value

// Parse parses a YYYY-MM-DD string.
func Parse(s string) (Date, error) {
	if _, err := time.Parse(Layout, s); err != nil {
		return "", fmt.Errorf("invalid date %q: %w", s, err)
	}
	return Date(s), nil
}

// From returns a Date from [time.Time], keeping only the calendar part of the date.
func From(t time.Time) Date {
	return Date(t.Format(Layout))
}

// Time returns the midnight in UTC matching the date.
func (d Date) Time() time.Time {
	t, _ := time.Parse(Layout, string(d))
	return t
}

// String returns the date in YYYY-MM-DD format.
func (d Date) String() string {
	return string(d)
}

// UnmarshalJSON accepts YYYY-MM-DD (the main format) and, for compatibility,
// RFC3339; the result is always normalized to a calendar date.
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		*d = ""
		return nil
	}
	if _, err := time.Parse(Layout, s); err == nil {
		*d = Date(s)
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("invalid date %q: %w", s, err)
	}
	*d = Date(t.Format(Layout))
	return nil
}

// MarshalJSON encodes the date as a YYYY-MM-DD string.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d + `"`), nil
}
