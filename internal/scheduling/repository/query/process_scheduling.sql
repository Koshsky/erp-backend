-- name: GetProcessScheduling :many
SELECT sqlc.embed(p) processes,
    sqlc.embed(pr) projects
FROM processes p
JOIN projects pr ON p.project_id = pr.id
WHERE  deleted_at IS NULL
    AND (
        pr.owner_id = @user_id OR
        @role = 'ДП'
    )
ORDER BY pr.priority ASC;