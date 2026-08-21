-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
  -- Для не-admin — только прямые подчинённые (manager_id = текущий пользователь);
  -- admin видит всех. «Сам пользователь» сюда не включается (табель добавляет
  -- себя на клиенте отдельно).
  AND (@is_admin::bool OR manager_id = @user_id::bigint)
  AND (@role_filter::text = '' OR role = @role_filter::text)
  AND (@manager_id::bigint = 0 OR manager_id = @manager_id::bigint)
ORDER BY id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountUsers :one
SELECT COUNT(*)
FROM users
WHERE deleted_at IS NULL
  AND (@is_admin::bool OR manager_id = @user_id::bigint)
  AND (@role_filter::text = '' OR role = @role_filter::text)
  AND (@manager_id::bigint = 0 OR manager_id = @manager_id::bigint);

-- name: ListAllUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: FindUser :one
SELECT *
FROM users
WHERE id = @user_id
	AND deleted_at IS NULL
LIMIT 1;

-- name: FindUserByUsername :one
SELECT *
FROM users
WHERE username = @username
	AND deleted_at IS NULL
LIMIT 1;

-- name: UsernameExists :one
SELECT EXISTS(
	SELECT 1
	FROM users
	WHERE username = @username
		AND deleted_at IS NULL
);

-- name: CreateUser :one
INSERT INTO users (last_name, first_name, middle_name, username, role, password_hash, manager_id, position, hire_date, termination_date)
VALUES (@last_name, @first_name, @middle_name, @username, @role, @password_hash, @manager_id, @position, @hire_date, @termination_date)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
	last_name = @last_name,
	first_name = @first_name,
	middle_name = @middle_name,
	username = @username,
	role = @role,
	password_hash = @password_hash,
	manager_id = @manager_id,
	position = @position,
	hire_date = @hire_date,
	termination_date = @termination_date,
	updated_at = NOW()
WHERE id = @user_id
	AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = NOW()
WHERE id = @user_id
	AND deleted_at IS NULL;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @user_id
	AND deleted_at IS NULL;

-- name: OwnerChain :one
-- Владелец записи: руководитель, а при его отсутствии — сам пользователь
-- (чтобы пользователь мог видеть/менять свой табель).
SELECT COALESCE(manager_id, id)::bigint AS owner_id
FROM users
WHERE id = @id::bigint
	AND deleted_at IS NULL;

-- ================= worker days (user_states) =================

-- name: ListStatesByUserRange :many
SELECT es.id, es.user_id, es.start_date, es.end_date, es.state_id,
	s.code AS state_code, s.name AS state_name, s.is_available
FROM user_states es
JOIN states s ON s.id = es.state_id
WHERE es.user_id = @user_id::bigint
	AND es.end_date >= @start_date::date
	AND es.start_date <= @end_date::date
ORDER BY es.start_date ASC;

-- name: ListOverlappingStates :many
SELECT es.id, es.user_id, es.start_date, es.end_date, es.state_id
FROM user_states es
WHERE es.user_id = @user_id::bigint
	AND es.end_date >= @start_date::date
	AND es.start_date <= @end_date::date
ORDER BY es.start_date ASC
FOR UPDATE;

-- name: ListOverlappingStatesByState :many
SELECT es.id, es.user_id, es.start_date, es.end_date, es.state_id
FROM user_states es
WHERE es.user_id = @user_id::bigint
	AND es.state_id = @state_id::bigint
	AND es.end_date >= @start_date::date
	AND es.start_date <= @end_date::date
ORDER BY es.start_date ASC
FOR UPDATE;

-- name: DeleteOverlapping :exec
DELETE FROM user_states
WHERE user_id = @user_id::bigint
	AND end_date >= @start_date::date
	AND start_date <= @end_date::date;

-- name: DeleteOverlappingByState :exec
DELETE FROM user_states
WHERE user_id = @user_id::bigint
	AND state_id = @state_id::bigint
	AND end_date >= @start_date::date
	AND start_date <= @end_date::date;

-- name: InsertStateRange :one
INSERT INTO user_states (user_id, state_id, start_date, end_date)
VALUES (@user_id, @state_id, @start_date::date, @end_date::date)
RETURNING *;

-- name: NormalizeUserStates :exec
-- Сливает смежные диапазоны одинакового (user_id, state_id) в непрерывные.
SELECT fn_normalize_user_states();
