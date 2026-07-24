-- TODO: write CanUser...Resource

-- name: ListResources :many
SELECT *
FROM resources
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: CreateResource :one
INSERT INTO resources (title, code, quantity)
VALUES (@title, @code, @quantity)
RETURNING *;

-- name: GetResource :one
SELECT *
FROM resources
WHERE deleted_at IS NULL
	AND id = @resource_id::bigint;

-- name: UpdateResource :one
UPDATE resources
SET
	title = @title,
	code = @code,
	quantity = @quantity,
	updated_at = NOW()
WHERE id = @resource_id
    AND deleted_at IS NULL
RETURNING *;

-- name: DeleteResource :exec
UPDATE resources
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @resource_id
    AND deleted_at IS NULL;