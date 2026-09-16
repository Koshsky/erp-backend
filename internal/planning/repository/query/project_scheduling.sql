-- name: ListProjects :many
SELECT * FROM projects
WHERE (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'own' AND owner_id = @user_id::bigint)
)
ORDER BY priority ASC;

-- name: ListProcesses :many
SELECT sqlc.embed(p), pr.code AS project_code
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND pr.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@scope_view::text = 'own' AND p.owner_id = @user_id::bigint)
);

-- name: ListProcessesByTaskScope :many
-- Processes for the TASK planning aggregate: the scope comes from the task
-- matrix, where 'parent' means "in my processes" (p.owner_id), not "in my
-- projects" (pr.owner_id) — the process list must follow the task semantics,
-- otherwise a process owner (vp) sees an empty task diagram although the
-- process view lists their processes.
SELECT sqlc.embed(p), pr.code AS project_code
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@scope_view::text = 'own' AND EXISTS (
        SELECT 1 FROM tasks t
        WHERE t.process_id = p.id
          AND t.owner_id = @user_id::bigint
    ))
);

-- name: ListResources :many
SELECT * FROM resources;


-- name: ListProjectsByIDs :many
SELECT * FROM projects
WHERE id = ANY(@ids::bigint[]);


-- name: ListProcessesByProjectIDs :many
SELECT * FROM processes
WHERE project_id = ANY(@project_ids::bigint[])
ORDER BY sort_order ASC, id ASC;

-- name: ListMilestonesByProcessIDs :many
SELECT * FROM milestones
WHERE process_id = ANY(@process_ids::bigint[])
ORDER BY id ASC;

-- name: ListTasksByProcessIDs :many
SELECT * FROM tasks
WHERE process_id = ANY(@process_ids::bigint[])
-- Top-level tasks (parent group 0) first in display order, then each
-- parent's subtasks in their own display order.
ORDER BY COALESCE(parent_id, 0), sort_order ASC, id ASC;

-- name: ListAssignmentsByTaskIDs :many
SELECT * FROM assignments
WHERE task_id = ANY(@task_ids::bigint[])
ORDER BY id ASC;

-- name: ListTaskCommentCountsByTaskIDs :many
SELECT task_id, COUNT(*)::bigint AS comments_count
FROM task_comments
WHERE task_id = ANY(@task_ids::bigint[])
GROUP BY task_id;