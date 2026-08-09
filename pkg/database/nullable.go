// Package nullable provides conversions between nullable database values
// (pgtype.Int8) and idiomatic Go pointer types.
package nullable

import "github.com/jackc/pgx/v5/pgtype"

// Int64Ptr converts a pgtype.Int8 into an *int64, returning nil when the value is not set.
func Int64Ptr(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	i := v.Int64
	return &i
}

// ToInt8 converts an *int64 into a pgtype.Int8, leaving it unset when v is nil.
func ToInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}
