-- name: CreateProject :one
INSERT INTO projects (code, start_date, end_date, priority)
VALUES ($1, $2, $3, $4)
RETURNING id, code, start_date, end_date, priority;

-- name: GetProject :one
SELECT id, code, start_date, end_date, priority
FROM projects
WHERE id = $1
	AND deleted_at IS NULL;

-- name: UpdateProject :one
UPDATE projects
SET
	code = $1,
	priority = $2,
	start_date = $3,
	end_date = $4,
	updated_at = NOW()
WHERE id = $5
	AND deleted_at IS NULL
RETURNING id, code, start_date, end_date, priority;

-- name: DeleteProject :exec
UPDATE projects
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
	AND deleted_at IS NULL;


-- name: ListProjects :many
SELECT id, code, start_date, end_date, priority
FROM projects
WHERE deleted_at IS NULL
ORDER BY priority ASC, start_date ASC, id ASC;