
-- name: ListProcesses :many
SELECT *
FROM processes p
WHERE p.deleted_at IS NULL
  AND (
    (@role::text = 'ДП') OR
    (p.owner_id = @user_id::bigint) OR
    (p.project_id IN (
        SELECT id FROM projects pr 
        WHERE pr.owner_id = @user_id::bigint
    ))
  )
ORDER BY 
    (SELECT priority FROM projects pr WHERE pr.id = p.project_id) ASC,
    p.start_date ASC, 
    p.end_date ASC, 
    p.id ASC;

-- name: ListTasksWithAssignmentsByProcessIDs :many
SELECT 
    t.id AS task_id,
    t.process_id,
    t.title AS task_title,
    t.start_date AS task_start_date,
    t.end_date AS task_end_date,
    t.created_at AS task_created_at,
    t.updated_at AS task_updated_at,
    t.deleted_at AS task_deleted_at,
    a.id AS assignment_id,
    a.task_id AS assignment_task_id,
    a.resource_id,
    a.quantity AS assignment_quantity,
    a.created_at AS assignment_created_at,
    a.updated_at AS assignment_updated_at,
    a.deleted_at AS assignment_deleted_at
FROM tasks t
LEFT JOIN assignments a ON a.task_id = t.id AND a.deleted_at IS NULL
WHERE t.process_id = ANY(@process_ids::bigint[]) AND t.deleted_at IS NULL
ORDER BY t.id, a.id;

-- name: ListMilestonesByProcessIDs :many
SELECT * FROM milestones 
WHERE process_id = ANY(@process_ids::bigint[]) AND deleted_at IS NULL 
ORDER BY id;

-- name: GetProcess :one
SELECT *
FROM processes
WHERE id = @id::bigint
	AND deleted_at IS NULL;

-- name: GetDetailedProcess :one
SELECT 
    p.id,
    p.project_id,
    p.owner_id,
    p.title,
    p.start_date,
    p.end_date,
    p.created_at,
    p.updated_at,
    p.deleted_at,
    COALESCE(
        (SELECT jsonb_agg(
            jsonb_build_object(
                'id', t.id,
                'process_id', t.process_id,
                'title', t.title,
                'start_date', t.start_date,
                'end_date', t.end_date,
                'created_at', t.created_at,
                'updated_at', t.updated_at,
                'deleted_at', t.deleted_at,
                'assignments', COALESCE(
                    (SELECT jsonb_agg(
                        jsonb_build_object(
                            'id', a.id,
                            'task_id', a.task_id,
                            'resource_id', a.resource_id,
                            'quantity', a.quantity
                        )
                        ORDER BY a.id
                    ) FROM assignments a WHERE a.task_id = t.id AND a.deleted_at IS NULL),
                    '[]'::jsonb
                )
            )
            ORDER BY t.id
        ) FROM tasks t WHERE t.process_id = p.id AND t.deleted_at IS NULL),
        '[]'::jsonb
    ) AS tasks,
    COALESCE(
        (SELECT jsonb_agg(
            jsonb_build_object(
                'id', m.id,
                'process_id', m.process_id,
                'title', m.title,
                'content', m.content,
                'date', m.date,
                'created_at', m.created_at,
                'updated_at', m.updated_at,
                'deleted_at', m.deleted_at
            )
            ORDER BY m.id
        ) FROM milestones m WHERE m.process_id = p.id AND m.deleted_at IS NULL),
        '[]'::jsonb
    ) AS milestones
FROM processes p
WHERE p.id = $1 AND p.deleted_at IS NULL;
    
-- name: CanUserManageProcess :one
SELECT EXISTS (
    SELECT 1 FROM processes p
    JOIN projects pr ON pr.id = p.project_id
    WHERE p.id = @process_id::bigint
      AND p.deleted_at IS NULL
      AND pr.owner_id = @user_id::bigint
) AS can_manage;

-- name: CanUserCreateProcess :one
SELECT EXISTS(
	SELECT 1 FROM projects
	WHERE id = @project_id::bigint
		AND deleted_at IS NULL
		AND owner_id = @user_id::bigint
) AS can_create;

-- name: CreateProcess :one
INSERT INTO processes (project_id, title, start_date, end_date, owner_id)
VALUES (
	@project_id::bigint,
	@title::text,
	@start_date::date,
	@end_date::date,
	@owner_id::bigint
)
RETURNING *;

-- name: UpdateProcess :one
UPDATE processes
SET
	title = @title,
	start_date = @start_date,
	end_date = @end_date,
	owner_id = COALESCE(@owner_id, owner_id),
	updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @process_id::bigint
RETURNING *;

-- name: DeleteProcess :exec
UPDATE processes
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @process_id::bigint;