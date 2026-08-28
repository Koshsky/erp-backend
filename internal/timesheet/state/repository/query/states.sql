-- name: ListStates :many
SELECT * FROM states
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: FindState :one
SELECT * FROM states
WHERE deleted_at IS NULL
	AND id = @state_id::bigint;

-- name: CreateState :one
-- Idempotent create by business key code: if the code already exists we
-- insert nothing; the calling code (repository) turns the conflict into 409.
INSERT INTO states (code, name, is_available)
VALUES (@code, @name, @is_available)
ON CONFLICT ON CONSTRAINT states_code_key
DO NOTHING
RETURNING *;

-- name: UpdateState :one
UPDATE states
SET
	code = @code,
	name = @name,
	is_available = @is_available,
	updated_at = NOW()
WHERE id = @state_id
	AND deleted_at IS NULL
RETURNING *;

-- name: DeleteState :exec
UPDATE states
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @state_id
	AND deleted_at IS NULL;
