package errors

import (
	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// FromPgInvalidParam превращает ошибку Postgres с кодом 22023
// (invalid_parameter_value — выход дат за границы родителя в V16-триггерах)
// в 400-ошибку с сообщением СУБД; остальные ошибки возвращаются как есть.
func FromPgInvalidParam(err error) error {
	var pgErr *pgconn.PgError
	if stderrors.As(err, &pgErr) && pgErr.Code == "22023" {
		return BadRequest(pgErr.Message)
	}
	return err
}
