-- TODO: CanUser...User

-- name: ListUsers :many
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

-- name: CreateUser :one
INSERT INTO users (name, username, role, password_hash, manager_id)
VALUES (@name, @username, @role, @password_hash, @manager_id)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
	name = @name,
	username = @username,
	role = @role,
	password_hash = @password_hash,
	manager_id = @manager_id,
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

-- name: ListSubordinates :many
SELECT *
FROM users
WHERE manager_id = @manager_id
	AND deleted_at IS NULL
ORDER BY id ASC;

