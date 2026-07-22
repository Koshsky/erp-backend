-- =============================================
-- CONSTRAINTS with soft delete support
-- =============================================

-- ---------- users ----------
CREATE UNIQUE INDEX users_username_unique_active ON users (username) WHERE deleted_at IS NULL;
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('ДП', 'РП', 'ВП'));

-- ---------- projects ----------
CREATE UNIQUE INDEX projects_code_unique_active ON projects (code) WHERE deleted_at IS NULL;
ALTER TABLE projects ADD CONSTRAINT projects_priority_check CHECK (priority >= 0);
ALTER TABLE projects ADD CONSTRAINT projects_dates_check CHECK (end_date >= start_date);

-- ---------- processes ----------
ALTER TABLE processes ADD CONSTRAINT processes_dates_check CHECK (end_date >= start_date);

-- ---------- tasks ----------
ALTER TABLE tasks ADD CONSTRAINT tasks_dates_check CHECK (end_date >= start_date);

-- ---------- resources ----------
CREATE UNIQUE INDEX resources_title_unique_active ON resources (title) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX resources_code_unique_active ON resources (code) WHERE deleted_at IS NULL;
ALTER TABLE resources ADD CONSTRAINT resources_quantity_check CHECK (quantity >= 0);
ALTER TABLE resources ADD CONSTRAINT resources_code_check CHECK (length(trim(code)) > 0);

-- ---------- assignments ----------
CREATE UNIQUE INDEX unique_task_resource_active ON assignments (task_id, resource_id) WHERE deleted_at IS NULL;
ALTER TABLE assignments ADD CONSTRAINT assignments_quantity_check CHECK (quantity > 0);
