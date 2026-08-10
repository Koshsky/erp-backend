-- name: CreateMilestone :one
INSERT INTO milestones (process_id, title, content, date)
VALUES (@process_id, @title, @content, @date)
RETURNING *;

-- name: ListMilestones :many
SELECT *
FROM milestones
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: FindMilestone :one
SELECT *
FROM milestones
WHERE id = @milestone_id
	AND deleted_at IS NULL;
	
-- name: UpdateMilestone :one
UPDATE milestones
SET
	process_id = @process_id,
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

-- name: OwnerChain :one
SELECT COALESCE(pr.owner_id, 0)::bigint AS project_owner,
       COALESCE(p.owner_id, 0)::bigint  AS process_owner
FROM milestones m
JOIN processes p ON p.id = m.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE m.id = @id::bigint
	AND m.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND pr.deleted_at IS NULL;
