-- name: ListProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'own' AND owner_id = @user_id::bigint)
)
ORDER BY priority ASC;

-- name: ListProcesses :many
SELECT sqlc.embed(p), pr.code AS project_code
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.deleted_at IS NULL
AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND pr.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@scope_view::text = 'own' AND p.owner_id = @user_id::bigint)
);

-- name: ListResources :many
SELECT * FROM resources
WHERE deleted_at IS NULL;


-- name: ListProjectsByIDs :many
SELECT * FROM projects
WHERE id = ANY(@ids::bigint[])
AND deleted_at IS NULL;


-- name: ListProcessesByProjectIDs :many
SELECT * FROM processes
WHERE project_id = ANY(@project_ids::bigint[])
AND deleted_at IS NULL
ORDER BY sort_order ASC, id ASC;

-- name: ListMilestonesByProcessIDs :many
SELECT * FROM milestones
WHERE process_id = ANY(@process_ids::bigint[])
AND deleted_at IS NULL
ORDER BY id ASC;

-- name: ListTasksByProcessIDs :many
SELECT * FROM tasks
WHERE process_id = ANY(@process_ids::bigint[])
AND deleted_at IS NULL
ORDER BY sort_order ASC, id ASC;

-- name: ListAssignmentsByTaskIDs :many
SELECT * FROM assignments
WHERE task_id = ANY(@task_ids::bigint[])
AND deleted_at IS NULL
ORDER BY id ASC;

-- name: ListTaskCommentCountsByTaskIDs :many
SELECT task_id, COUNT(*)::bigint AS comments_count
FROM task_comments
WHERE task_id = ANY(@task_ids::bigint[])
AND deleted_at IS NULL
GROUP BY task_id;