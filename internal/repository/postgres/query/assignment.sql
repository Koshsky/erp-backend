-- name: CreateAssignment :one
INSERT INTO assignments (task_id, resource_id, quantity)
VALUES ($1, $2, $3)
RETURNING id, task_id, resource_id, quantity;

-- name: GetAssignment :one
SELECT id, task_id, resource_id, quantity
FROM assignments
WHERE id = $1
	AND deleted_at IS NULL;

-- name: UpdateAssignment :one
UPDATE assignments
SET
	quantity = $1,
	updated_at = NOW()
WHERE id = $2
	AND deleted_at IS NULL
RETURNING id, task_id, resource_id, quantity;

-- name: DeleteAssignment :exec
UPDATE assignments
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
	AND deleted_at IS NULL;


-- name: ListAssignmentsByTaskID :many
SELECT id, task_id, resource_id, quantity
FROM assignments
WHERE task_id = $1
	AND deleted_at IS NULL
ORDER BY id ASC;