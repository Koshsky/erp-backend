-- name: GetMilestone :one
SELECT *
FROM milestones
WHERE id = @milestone_id
	AND deleted_at IS NULL;
	
-- name: CanUserManageMilestone :one
SELECT EXISTS(
	SELECT 1 FROM milestones m
	JOIN processes p ON p.id = m.process_id
	WHERE m.id = @milestone_id::bigint
		AND m.deleted_at IS NULL
		AND p.owner_id = @user_id::bigint
) AS can_manage;

-- name: CanUserCreateMilestone :one
SELECT EXISTS(
	SELECT 1 FROM processes p
	WHERE p.id = @process_id::bigint
		AND p.deleted_at IS NULL
		AND p.owner_id = @user_id::bigint
) AS can_create;

-- name: CreateMilestone :one
INSERT INTO milestones (process_id, title, content, date)
VALUES (@process_id, @title, @content, @date)
RETURNING *;

-- name: UpdateMilestone :one
UPDATE milestones
SET
	title = @title,
	content = @content,
	date = @date,
	updated_at = NOW()
WHERE id = @milestone_id
	AND deleted_at IS NULL
RETURNING *;

-- name: DeleteMilestone :exec
UPDATE milestones
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @milestone_id
	AND deleted_at IS NULL;