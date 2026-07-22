-- name: CreateMilestone :one
INSERT INTO milestones (process_id, title, content, date)
VALUES ($1, $2, $3, $4)
RETURNING id, process_id, title, content, date;

-- name: GetMilestone :one
SELECT id, process_id, title, content, date
FROM milestones
WHERE id = $1
	AND deleted_at IS NULL;

-- name: UpdateMilestone :one
UPDATE milestones
SET
	title = $1,
	content = $2,
	date = $3,
	updated_at = NOW()
WHERE id = $4
	AND deleted_at IS NULL
RETURNING id, process_id, title, content, date;

-- name: DeleteMilestone :exec
UPDATE milestones
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
	AND deleted_at IS NULL;


-- name: ListMilestonesByProcessID :many
SELECT id, process_id, title, content, date
FROM milestones
WHERE process_id = $1
	AND deleted_at IS NULL
ORDER BY date ASC, id ASC;