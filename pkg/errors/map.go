package errors

import (
	"strings"

	stderrors "errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// MapPgConstraint превращает нарушение целостности Postgres в понятную
// HTTP-ошибку: уникальность (23505) → 409, внешний ключ (23503) → 400,
// нарушение CHECK (23514) → 400. Прочие ошибки возвращаются без изменений.
func MapPgConstraint(err error) error {
	var pgErr *pgconn.PgError
	if !stderrors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505":
		return Conflict(uniqueMessage(pgErr.ConstraintName))
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

// uniqueMessage формулирует человеческое сообщение для 23505 по имени
// ограничения (например, частичный уникальный индекс логина пользователя).
func uniqueMessage(constraint string) string {
	if strings.Contains(constraint, "username") {
		return "пользователь с таким логином уже существует"
	}
	return "запись с таким значением уже существует"
}
