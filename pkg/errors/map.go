package errors

import (
	"strings"

	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// MapPgConstraint maps a Postgres integrity violation into a domain error:
// unique (23505) → 409, foreign key (23503) → 400, CHECK (23514) → 400,
// exclusion (23P01, overlapping ranges such as user_states) → 409.
// All other errors are returned unchanged.
func MapPgConstraint(err error) error {
	var pgErr *pgconn.PgError
	if !stderrors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return Conflict(uniqueMessage(pgErr.ConstraintName))
	case "23P01":
		return Conflict("значение пересекается с существующим диапазоном")
	case "23503":
		if pgErr.ConstraintName == "users_role_fk" {
			return BadRequest("неизвестная роль: её нет в каталоге ролей")
		}
		return BadRequest("запись ссылается на несуществующий объект: " + pgErr.ConstraintName)
	case "23514":
		return BadRequest("значение не прошло проверку базы данных: " + pgErr.ConstraintName)
	default:
		return err
	}
}

// uniqueMessage builds the human-readable message for 23505 from the
// constraint name (e.g. the partial unique index on the user login).
func uniqueMessage(constraint string) string {
	if strings.Contains(constraint, "username") {
		return "пользователь с таким логином уже существует"
	}
	return "запись с таким значением уже существует"
}
