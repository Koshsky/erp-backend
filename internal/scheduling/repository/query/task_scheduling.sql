-- name: GetDescribedTasks :many
SELECT sqlc.embed(t), sqlc.embed(p), sqlc.embed(pr)
FROM tasks t
JOIN processes p ON t.process_id = p.id
JOIN projects pr ON p.project_id = pr.id
WHERE t.deleted_at IS NULL
  AND (
    @role = 'ДП' OR
    p.owner_id = @user_id OR
    pr.owner_id = @user_id
  )
ORDER BY pr.priority ASC, p.id ASC;

-- name: ListAssignmentsByTaskIDs :many
SELECT *
FROM assignments a
WHERE a.task_id = ANY(@task_ids::bigint[])
  AND a.deleted_at IS NULL;

-- name: ListMilestonesByProcessIDs :many
SELECT *
FROM milestones m
WHERE m.process_id = ANY(@process_ids::bigint[])
  AND m.deleted_at IS NULL;

-- name: GetResources :many
SELECT * FROM resources
WHERE deleted_at IS NULL;