// Package date — календарная дата в формате YYYY-MM-DD (без времени и часового пояса).
//
// Единый контракт API для всех «календарных» полей (start_date, end_date, date):
// в JSON они кодируются строкой YYYY-MM-DD (OpenAPI format: date), что совпадает
// с типом DATE в БД и дневной сеткой календаря на фронтенде. [time.Time] в этих
// полях использовал бы RFC3339 с временем и зоной — источник расхождений формата.
package date

import (
	"fmt"
	"strings"
	"time"
)

// Layout — единственный формат дат API.
const Layout = "2006-01-02"

// Date — календарная дата без времени, хранится как строка YYYY-MM-DD.
// Строковое представление позволяет swag (OpenAPI) видеть его как string.
type Date string //nolint:recvcheck // json.Unmarshaler требует указатель, остальные методы — на значении

// Parse разбирает строку YYYY-MM-DD.
func Parse(s string) (Date, error) {
	if _, err := time.Parse(Layout, s); err != nil {
		return "", fmt.Errorf("invalid date %q: %w", s, err)
	}
	return Date(s), nil
}

// From возвращает Date по [time.Time], оставляя только календарную часть даты.
func From(t time.Time) Date {
	return Date(t.Format(Layout))
}

// Time возвращает полночь в UTC, соответствующую дате.
func (d Date) Time() time.Time {
	t, _ := time.Parse(Layout, string(d))
	return t
}

// String возвращает дату в формате YYYY-MM-DD.
func (d Date) String() string {
	return string(d)
}

// UnmarshalJSON принимает YYYY-MM-DD (основной формат) и, для совместимости,
// RFC3339; на выходе всегда нормализуется к календарной дате.
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

// MarshalJSON кодирует дату строкой YYYY-MM-DD.
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d + `"`), nil
}
