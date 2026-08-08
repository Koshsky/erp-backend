-- name: ListResources :many
SELECT r.id, r.code, r.title, r.owner_id,
    COUNT(e.id) FILTER (WHERE e.deleted_at IS NULL)::bigint AS employees_count,
    r.created_at, r.updated_at, r.deleted_at
FROM resources r
LEFT JOIN employees e ON e.resource_id = r.id
WHERE r.deleted_at IS NULL
GROUP BY r.id, r.code, r.title, r.owner_id, r.created_at, r.updated_at, r.deleted_at
ORDER BY r.id ASC;

-- name: ListResourcesByOwnerID :many
SELECT r.id, r.code, r.title, r.owner_id,
    COUNT(e.id) FILTER (WHERE e.deleted_at IS NULL)::bigint AS employees_count,
    r.created_at, r.updated_at, r.deleted_at
FROM resources r
LEFT JOIN employees e ON e.resource_id = r.id
WHERE r.deleted_at IS NULL
	AND r.owner_id = @owner_id::bigint
GROUP BY r.id, r.code, r.title, r.owner_id, r.created_at, r.updated_at, r.deleted_at
ORDER BY r.id ASC;

-- name: FindResource :one
SELECT r.id, r.code, r.title, r.owner_id,
    COUNT(e.id) FILTER (WHERE e.deleted_at IS NULL)::bigint AS employees_count,
    r.created_at, r.updated_at, r.deleted_at
FROM resources r
LEFT JOIN employees e ON e.resource_id = r.id
WHERE r.deleted_at IS NULL
	AND r.id = @resource_id::bigint
GROUP BY r.id, r.code, r.title, r.owner_id, r.created_at, r.updated_at, r.deleted_at;

-- name: CreateResource :one
INSERT INTO resources (title, code, owner_id)
VALUES (@title, @code, @owner_id)
RETURNING *;

-- name: CountEmployeesByResourceID :one
SELECT COUNT(*)::bigint
FROM employees
WHERE resource_id = @resource_id::bigint
	AND deleted_at IS NULL;

-- name: UpdateResource :one
UPDATE resources
SET
	title = @title,
	code = @code,
	owner_id = @owner_id,
	updated_at = NOW()
WHERE id = @resource_id
    AND deleted_at IS NULL
RETURNING *;

-- name: DeleteResource :exec
UPDATE resources
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @resource_id
    AND deleted_at IS NULL;

-- name: OwnerChain :one
SELECT COALESCE(owner_id, 0)::bigint AS owner_id
FROM resources
WHERE id = @id::bigint
	AND deleted_at IS NULL;
