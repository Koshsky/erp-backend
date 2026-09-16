package errors

import (
	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// FromPgInvalidParam maps a Postgres error with code 22023
// (invalid_parameter_value — dates exceeding a parent's bounds, raised by the
// V16 triggers) into a 400 error carrying the DB message; all other errors are
// returned unchanged.
func FromPgInvalidParam(err error) error {
	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) && pgErr.Code == "22023" {
		return BadRequest(pgErr.Message)
	}
	return err
}
