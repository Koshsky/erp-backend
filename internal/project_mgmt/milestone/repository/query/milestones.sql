-- name: CreateMilestone :one
INSERT INTO milestones (process_id, title, content, date)
VALUES (@process_id, @title, @content, @date)
RETURNING *;

-- name: ListMilestones :many
SELECT m.*
FROM milestones m
JOIN processes p ON p.id = m.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE m.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint))
  )
  AND (@owner_id::bigint = 0 OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint)
ORDER BY m.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountMilestones :one
SELECT COUNT(*)
FROM milestones m
JOIN processes p ON p.id = m.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE m.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint))
  )
  AND (@owner_id::bigint = 0 OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint);

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
