-- name: CreateUser :one
INSERT INTO users (name, username, role, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id, name, username, role, password_hash;

-- name: GetUser :one
SELECT id, name, username, role, password_hash
FROM users
WHERE id = $1
	AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
	name = $1,
	username = $2,
	role = $3,
	password_hash = $4,
	updated_at = NOW()
WHERE id = $5
	AND deleted_at IS NULL
RETURNING id, name, username, role, password_hash;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
	AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id, name, username, role, password_hash
FROM users
WHERE deleted_at IS NULL
ORDER BY username ASC, id ASC;
