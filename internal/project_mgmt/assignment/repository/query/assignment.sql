-- name: CreateAssignment :one
INSERT INTO assignments (task_id, resource_id, quantity)
VALUES (@task_id, @resource_id, @quantity::bigint)
RETURNING *;

-- name: FindAssignment :one
SELECT *
FROM assignments
WHERE id = @assignment_id
	AND deleted_at IS NULL;
-- name: ListAssigments :many
SELECT *
FROM assignments
WHERE deleted_at IS NULL
ORDER BY id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountAssignments :one
SELECT COUNT(*)
FROM assignments
WHERE deleted_at IS NULL;

-- name: UpdateAssignment :one
UPDATE assignments
SET
	task_id = @task_id,
	resource_id = @resource_id,
	quantity = @quantity::bigint,
	updated_at = NOW()
WHERE id = @assignment_id
	AND deleted_at IS NULL
RETURNING *;

-- name: DeleteAssignment :exec
UPDATE assignments
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @assignment_id
	AND deleted_at IS NULL;

-- name: OwnerChain :one
SELECT COALESCE(pr.owner_id, 0)::bigint AS project_owner,
       COALESCE(p.owner_id, 0)::bigint  AS process_owner
FROM assignments a
JOIN tasks t ON t.id = a.task_id
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE a.id = @id::bigint
	AND a.deleted_at IS NULL
	AND t.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND pr.deleted_at IS NULL;
