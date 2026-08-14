-- name: ListResources :many
SELECT r.id, r.code, r.title, r.owner_id,
    COUNT(rm.user_id)::bigint AS employees_count,
    r.created_at, r.updated_at, r.deleted_at
FROM resources r
LEFT JOIN resource_members rm ON rm.resource_id = r.id
WHERE r.deleted_at IS NULL
  AND (@role::text = 'admin' OR r.owner_id = @user_id::bigint)
  AND (@owner_id::bigint = 0 OR r.owner_id = @owner_id::bigint)
GROUP BY r.id, r.code, r.title, r.owner_id, r.created_at, r.updated_at, r.deleted_at
ORDER BY r.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountResources :one
SELECT COUNT(*)
FROM resources
WHERE deleted_at IS NULL
  AND (@role::text = 'admin' OR owner_id = @user_id::bigint)
  AND (@owner_id::bigint = 0 OR owner_id = @owner_id::bigint);

-- name: ListResourcesByOwnerID :many
SELECT r.id, r.code, r.title, r.owner_id,
    COUNT(rm.user_id)::bigint AS employees_count,
    r.created_at, r.updated_at, r.deleted_at
FROM resources r
LEFT JOIN resource_members rm ON rm.resource_id = r.id
WHERE r.deleted_at IS NULL
	AND r.owner_id = @owner_id::bigint
GROUP BY r.id, r.code, r.title, r.owner_id, r.created_at, r.updated_at, r.deleted_at
ORDER BY r.id ASC;

-- name: FindResource :one
SELECT r.id, r.code, r.title, r.owner_id,
    COUNT(rm.user_id)::bigint AS employees_count,
    r.created_at, r.updated_at, r.deleted_at
FROM resources r
LEFT JOIN resource_members rm ON rm.resource_id = r.id
WHERE r.deleted_at IS NULL
	AND r.id = @resource_id::bigint
GROUP BY r.id, r.code, r.title, r.owner_id, r.created_at, r.updated_at, r.deleted_at;

-- name: CreateResource :one
INSERT INTO resources (title, code, owner_id)
VALUES (@title, @code, @owner_id)
RETURNING *;

-- name: CountMembersByResourceID :one
SELECT COUNT(*)::bigint
FROM resource_members
WHERE resource_id = @resource_id::bigint;

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

-- ================= resource members =================

-- name: ListMembersByResourceID :many
SELECT u.id, u.name, u.role, u.position, u.hire_date, u.termination_date, u.manager_id
FROM resource_members rm
JOIN users u ON u.id = rm.user_id
WHERE rm.resource_id = @resource_id::bigint
  AND u.deleted_at IS NULL
ORDER BY u.id ASC;

-- name: AddMember :exec
INSERT INTO resource_members (resource_id, user_id)
VALUES (@resource_id::bigint, @user_id::bigint)
ON CONFLICT DO NOTHING;

-- name: RemoveMember :exec
DELETE FROM resource_members
WHERE resource_id = @resource_id::bigint
  AND user_id = @user_id::bigint;

-- name: FindUserManager :one
SELECT manager_id
FROM users
WHERE id = @user_id::bigint
  AND deleted_at IS NULL;
