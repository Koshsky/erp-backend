-- name: ListProjects :many
SELECT * FROM projects
WHERE deleted_at IS NULL
AND (
    @role = 'ДП' OR
    owner_id = @user_id
)
ORDER BY priority ASC;

-- name: ListProcessesByProjectID :many
SELECT * FROM processes
WHERE project_id = ANY(@project_ids::bigint[])
AND deleted_at IS NULL;