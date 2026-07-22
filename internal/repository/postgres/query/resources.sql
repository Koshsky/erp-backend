-- name: CreateResource :one
INSERT INTO resources (title, code, quantity)
VALUES ($1, $2, $3)
RETURNING id, title, code, quantity;

-- name: GetResource :one
SELECT id, title, code, quantity
FROM resources
WHERE id = $1
    AND deleted_at IS NULL;

-- name: UpdateResource :one
UPDATE resources
SET
	title = $1,
	code = $2,
	quantity = $3,
	updated_at = NOW()
WHERE id = $4
    AND deleted_at IS NULL
RETURNING id, title, code, quantity;

-- name: DeleteResource :exec
UPDATE resources
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1
    AND deleted_at IS NULL;


-- name: ListResources :many
SELECT id, title, code, quantity
FROM resources
WHERE deleted_at IS NULL
ORDER BY title ASC, id ASC;

-- name: GetResourceUsage :many
SELECT
    r.id,
    r.title,
    r.code,
    r.quantity AS total_quantity,
    COALESCE(SUM(a.quantity), 0)::BIGINT AS used_quantity,
    (r.quantity - COALESCE(SUM(a.quantity), 0))::BIGINT AS available_quantity
FROM resources r
LEFT JOIN assignments a ON a.resource_id = r.id
    AND a.deleted_at IS NULL
LEFT JOIN tasks t ON a.task_id = t.id
    AND t.start_date <= sqlc.arg('target_date')::date
    AND t.end_date > sqlc.arg('target_date')::date
    AND t.deleted_at IS NULL
WHERE r.deleted_at IS NULL
GROUP BY r.id, r.title, r.quantity
ORDER BY r.title ASC;