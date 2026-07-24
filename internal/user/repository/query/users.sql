-- TODO: CanUser...User

-- name: ListUsers :many
SELECT *
FROM users
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = @user_id
	AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, username, role, password_hash)
VALUES (@name, @username, @role, @password_hash)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
	name = @name,
	username = @username,
	role = @role,
	password_hash = @password_hash,
	updated_at = NOW()
WHERE id = @user_id
	AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @user_id
	AND deleted_at IS NULL;
