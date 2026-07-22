-- name: CreateTask :one
INSERT INTO tasks (process_id, title, start_date, end_date)
VALUES ($1, $2, $3, $4)
RETURNING id, process_id, title, start_date, end_date;

-- name: GetTask :one
SELECT id, process_id, title, start_date, end_date
FROM tasks
WHERE id = $1
	AND deleted_at IS NULL;

-- name: UpdateTask :one
UPDATE tasks
SET
	title = $1,
	start_date = $2,
	end_date = $3,
	updated_at = NOW()
WHERE id = $4
	AND deleted_at IS NULL
RETURNING id, process_id, title, start_date, end_date;

-- name: DeleteTask :exec
UPDATE tasks
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
	AND deleted_at IS NULL;


-- name: ListTasksByProcessID :many
SELECT id, process_id, title, start_date, end_date
FROM tasks
WHERE process_id = $1
	AND deleted_at IS NULL
ORDER BY start_date ASC, end_date ASC, id ASC;