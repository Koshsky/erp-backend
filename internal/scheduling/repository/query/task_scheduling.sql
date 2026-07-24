-- name: GetTaskScheduling :many
SELECT sqlc.embed(t) tasks,
    sqlc.embed(p) processes,
    sqlc.embed(pr) projects,
    sqlc.embed(a) assignments
FROM tasks t
JOIN processes p ON t.process_id = p.id
JOIN projects pr ON p.project_id = pr.id
LEFT JOIN assignments a ON t.id = a.task_id
WHERE tasks.deleted_at IS NULL
    AND (
        p.owner_id = @user_id OR
        pr.owner_id = @user_id OR
        @role = 'ДП'
    )
ORDER BY pr.priority ASC;

-- name: GetResources :many
SELECT * FROM resources
WHERE deleted_at IS NULL;