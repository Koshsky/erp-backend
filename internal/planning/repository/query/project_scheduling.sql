-- name: ListProjects :many
SELECT * FROM projects
WHERE (
    @sc_all::boolean OR
    (@sc_self::boolean AND owner_id = @user_id::bigint) OR
    (@sc_down::boolean AND EXISTS (
        SELECT 1 FROM tasks d
        JOIN processes dp ON dp.id = d.process_id
        WHERE dp.project_id = projects.id AND d.owner_id = @user_id::bigint
    )) OR
    @sc_none::boolean
)
ORDER BY priority ASC;

-- name: ListProcesses :many
SELECT sqlc.embed(p), pr.code AS project_code
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE (
    @sc_all::boolean OR
    (@sc_self::boolean AND p.owner_id = @user_id::bigint) OR
    (@sc_parent::boolean AND pr.owner_id = @user_id::bigint) OR
    (@sc_ancestor::boolean AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@sc_sib::boolean AND EXISTS (
        SELECT 1 FROM processes s
        WHERE s.project_id = p.project_id AND s.owner_id = @user_id::bigint
    )) OR
    (@sc_down::boolean AND EXISTS (
        SELECT 1 FROM tasks d WHERE d.process_id = p.id AND d.owner_id = @user_id::bigint
    )) OR
    @sc_none::boolean
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
    @sc_all::boolean OR
    (@sc_self::boolean AND EXISTS (
        SELECT 1 FROM tasks t WHERE t.process_id = p.id AND t.owner_id = @user_id::bigint
    )) OR
    (@sc_parent::boolean AND p.owner_id = @user_id::bigint) OR
    (@sc_ancestor::boolean AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@sc_sib::boolean AND EXISTS (
        SELECT 1 FROM processes s
        WHERE s.project_id = p.project_id AND s.owner_id = @user_id::bigint
    )) OR
    (@sc_down::boolean AND EXISTS (
        SELECT 1 FROM tasks d
        WHERE d.process_id = p.id AND d.owner_id = @user_id::bigint
    )) OR
    @sc_none::boolean
);

-- name: ListResources :many
-- Resources embedded into task/reference data, scoped by the caller's
-- resource view zone: a task-view holder with a resource zone must not
-- receive resource rows outside it (resource is own-scoped for vp). A caller
-- with no resource view rule (empty zone — "no rule / none") gets every row:
-- reference data must not hide the visible tasks.
SELECT * FROM resources
WHERE (
    @sc_all::boolean OR
    (@sc_self::boolean AND owner_id = @user_id::bigint) OR
    @sc_none::boolean
);


-- name: ListProjectsByIDs :many
-- Projects attached as parent context of the caller's visible processes,
-- scoped by the caller's project view zone: the process aggregate must not
-- disclose full project rows of projects outside that zone. A caller with no
-- project view rule (empty zone — e.g. vp, which sees processes through
-- process.view=all) gets every requested parent project: reference rows with
-- no resolved zone must not hide the visible processes.
SELECT * FROM projects
WHERE id = ANY(@ids::bigint[])
  AND (
    @sc_all::boolean OR
    (@sc_self::boolean AND owner_id = @user_id::bigint) OR
    (@sc_down::boolean AND EXISTS (
        SELECT 1 FROM tasks d
        JOIN processes dp ON dp.id = d.process_id
        WHERE dp.project_id = projects.id AND d.owner_id = @user_id::bigint
    )) OR
    @sc_none::boolean
  );


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

-- name: ListTaskDependenciesByProcessIDs :many
-- Scheduling links between tasks of the given processes (both ends are
-- top-level tasks of the same process — enforced at creation).
SELECT sqlc.embed(td)
FROM task_dependencies td
JOIN tasks t ON t.id = td.task_id
WHERE t.process_id = ANY(@process_ids::bigint[])
ORDER BY td.task_id ASC, td.id ASC;