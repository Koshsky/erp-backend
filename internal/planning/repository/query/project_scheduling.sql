-- name: ListProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
AND (
    @role::text IN ('admin', 'dp') OR
    owner_id = @user_id::bigint
)
ORDER BY priority ASC;

-- name: ListProcesses :many
SELECT sqlc.embed(p), pr.code AS project_code
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.deleted_at IS NULL
AND (
    @role::text IN ('admin', 'dp') OR
    p.owner_id = @user_id::bigint OR
    pr.owner_id = @user_id::bigint
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
AND deleted_at IS NULL;

-- name: ListMilestonesByProcessIDs :many
SELECT * FROM milestones
WHERE process_id = ANY(@process_ids::bigint[])
AND deleted_at IS NULL;

-- name: ListTasksByProcessIDs :many
SELECT * FROM tasks
WHERE process_id = ANY(@process_ids::bigint[])
AND deleted_at IS NULL;

-- name: ListAssignmentsByTaskIDs :many
SELECT * FROM assignments
WHERE task_id = ANY(@task_ids::bigint[])
AND deleted_at IS NULL;