-- name: CreateProcess :one
INSERT INTO processes (project_id, title, start_date, end_date)
VALUES ($1, $2, $3, $4)
RETURNING id, project_id, title, start_date, end_date;

-- name: GetProcess :one
SELECT id, project_id, title, start_date, end_date
FROM processes
WHERE id = $1
	AND deleted_at IS NULL;

-- name: UpdateProcess :one
UPDATE processes
SET
	title = $1,
	start_date = $2,
	end_date = $3,
	updated_at = NOW()
WHERE id = $4
	AND deleted_at IS NULL
RETURNING id, project_id, title, start_date, end_date;

-- name: DeleteProcess :exec
UPDATE processes
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
	AND deleted_at IS NULL;


-- name: ListProcessesByProjectID :many
SELECT id, project_id, title, start_date, end_date
FROM processes
WHERE project_id = $1
	AND deleted_at IS NULL
ORDER BY start_date ASC, end_date ASC, id ASC;