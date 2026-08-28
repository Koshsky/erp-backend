-- =============================================
-- 1. UNIQUE INDICES (soft delete)
-- =============================================
CREATE UNIQUE INDEX idx_users_username_active ON users(username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_projects_code_active ON projects(code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_resources_title_active ON resources(title) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_resources_code_active ON resources(code) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_assignments_unique_active
ON assignments(task_id, resource_id) WHERE deleted_at IS NULL;
-- Order uniqueness within the parent group (soft delete: freed on deletion).
CREATE UNIQUE INDEX idx_processes_project_sort_order
ON processes(project_id, sort_order) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_tasks_process_sort_order
ON tasks(process_id, sort_order) WHERE deleted_at IS NULL;

-- =============================================
-- 2. INDICES FOR Foreign Keys
-- =============================================
CREATE INDEX idx_processes_project_id ON processes(project_id);
CREATE INDEX idx_tasks_process_id ON tasks(process_id);
CREATE INDEX idx_assignments_task_id ON assignments(task_id);
CREATE INDEX idx_assignments_resource_id ON assignments(resource_id);
CREATE INDEX idx_milestones_process_id ON milestones(process_id);
CREATE INDEX idx_users_manager_id ON users(manager_id);
CREATE INDEX idx_user_states_user_id ON user_states(user_id);

-- =============================================
-- 3. INDICES FOR soft delete (search for active)
-- =============================================
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_processes_deleted_at ON processes(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_deleted_at ON tasks(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_resources_deleted_at ON resources(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_milestones_deleted_at ON milestones(deleted_at) WHERE deleted_at IS NULL;

-- =============================================
-- 4. COMPOSITE INDICES (JOIN + soft delete)
-- =============================================
CREATE INDEX idx_processes_project_deleted ON processes(project_id, deleted_at)
WHERE deleted_at IS NULL;
CREATE INDEX idx_tasks_process_deleted ON tasks(process_id, deleted_at)
WHERE deleted_at IS NULL;
CREATE INDEX idx_assignments_task_deleted ON assignments(task_id, deleted_at)
WHERE deleted_at IS NULL;
CREATE INDEX idx_assignments_resource_deleted ON assignments(resource_id, deleted_at)
WHERE deleted_at IS NULL;
CREATE INDEX idx_milestones_process_deleted ON milestones(process_id, deleted_at)
WHERE deleted_at IS NULL;

-- =============================================
-- 5. INDICES FOR dates
-- =============================================
CREATE INDEX idx_user_states_dates ON user_states (start_date, end_date);
CREATE INDEX idx_tasks_dates ON tasks(start_date, end_date);
CREATE INDEX idx_milestones_date ON milestones(date);

-- =============================================
-- 6. INDICES FOR auth / idempotency / comments
-- =============================================
CREATE INDEX refresh_sessions_user_idx    ON refresh_sessions (user_id);
CREATE INDEX refresh_sessions_expires_idx ON refresh_sessions (expires_at);
CREATE INDEX idx_idempotency_keys_expires ON idempotency_keys (expires_at);
CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);
CREATE INDEX idx_task_comments_task_created ON task_comments(task_id, created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_task_comments_parent_id ON task_comments(parent_id) WHERE deleted_at IS NULL;

-- -- =============================================
-- -- 7. OPTIONAL: full-text search
-- -- =============================================
-- If you want to enable full-text search on certain text fields, you can create GIN indices on those fields. For example:
-- CREATE INDEX idx_projects_code_trgm ON projects USING GIN (code gin_trgm_ops);
-- CREATE INDEX idx_tasks_title_trgm ON tasks USING GIN (title gin_trgm_ops);
-- CREATE INDEX idx_resources_title_trgm ON resources USING GIN (title gin_trgm_ops);