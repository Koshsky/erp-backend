-- name: CreateAssignment :one
-- Idempotent create: if the (task_id, resource_id) link is already active
-- we insert nothing; an existing row is reported back to the calling
-- code (repository) via FindAssignmentByKey.
INSERT INTO assignments (task_id, resource_id, quantity)
VALUES (@task_id, @resource_id, @quantity::bigint)
ON CONFLICT (task_id, resource_id)
DO NOTHING
RETURNING *;

-- name: FindAssignmentByKey :one
SELECT *
FROM assignments
WHERE task_id = @task_id::bigint
	AND resource_id = @resource_id::bigint
LIMIT 1;

-- name: FindAssignment :one
SELECT *
FROM assignments
WHERE id = @assignment_id;
-- name: ListAssigments :many
SELECT a.*
FROM assignments a
JOIN tasks t ON t.id = a.task_id
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (t.owner_id = @user_id::bigint OR p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint))
  )
  AND (@owner_id::bigint = 0 OR t.owner_id = @owner_id::bigint OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint)
ORDER BY a.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountAssignments :one
SELECT COUNT(*)
FROM assignments a
JOIN tasks t ON t.id = a.task_id
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (t.owner_id = @user_id::bigint OR p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint))
  )
  AND (@owner_id::bigint = 0 OR t.owner_id = @owner_id::bigint OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint);

-- name: UpdateAssignment :one
UPDATE assignments
SET
	task_id = @task_id,
	resource_id = @resource_id,
	quantity = @quantity::bigint,
	updated_at = NOW()
WHERE id = @assignment_id
RETURNING *;

-- name: DeleteAssignment :exec
DELETE FROM assignments
WHERE id = @assignment_id;

-- name: OwnerChain :one
SELECT COALESCE(pr.owner_id, 0)::bigint AS project_owner,
       COALESCE(p.owner_id, 0)::bigint  AS process_owner,
       COALESCE(t.owner_id, 0)::bigint  AS owner_id
FROM assignments a
JOIN tasks t ON t.id = a.task_id
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE a.id = @id::bigint;