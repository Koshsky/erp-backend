-- name: GetProcessScheduling :many
SELECT sqlc.embed(p),
    sqlc.embed(pr)
FROM processes p
JOIN projects pr ON p.project_id = pr.id
WHERE  p.deleted_at IS NULL
    AND (
        pr.owner_id = @user_id OR
        @role = 'ДП'
    )
ORDER BY pr.priority ASC;