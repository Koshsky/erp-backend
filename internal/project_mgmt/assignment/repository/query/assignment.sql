-- name: CanUserCreateAssignment :one
SELECT EXISTS (
    SELECT 1 FROM tasks t
	JOIN processes p ON p.id = t.process_id
    WHERE t.id = @task_id::bigint
	  AND t.deleted_at is NULL
      AND p.owner_id = @user_id::bigint
) AS can_create;
-- name: CreateAssignment :one
INSERT INTO assignments (task_id, resource_id, quantity)
VALUES (@task_id, @resource_id, @quantity::bigint)
RETURNING *;

-- TODO: write CanUserViewAssignment
-- name: FindAssignment :one
SELECT *
FROM assignments
WHERE id = @assignment_id
	AND deleted_at IS NULL;
-- name: ListAssigments :many
SELECT *
FROM assignments
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: CanUserUpdateAssignment :one
SELECT EXISTS (
    SELECT 1 FROM assignments a
	JOIN tasks t ON t.id = a.task_id
	JOIN processes p ON p.id = t.process_id
    WHERE a.id = @assignment_id::bigint
	  AND a.deleted_at is NULL
      AND p.owner_id = @user_id::bigint
) AS can_manage;
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

-- name: CanUserDeleteAssignment :one
SELECT EXISTS (
    SELECT 1 FROM assignments a
	JOIN tasks t ON t.id = a.task_id
	JOIN processes p ON p.id = t.process_id
    WHERE a.id = @assignment_id::bigint
	  AND a.deleted_at is NULL
      AND p.owner_id = @user_id::bigint
) AS can_manage;
-- name: DeleteAssignment :exec
UPDATE assignments
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @assignment_id
	AND deleted_at IS NULL;