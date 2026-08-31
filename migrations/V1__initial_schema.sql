-- Initial schema for enterprise resource planning.
-- Complete schema: all tables are created once, all CHECK/FK constraints
-- are declared inline. No subsequent ALTER TABLE is required.

-- The user_states EXCLUDE constraint relies on btree_gist.
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- =============================================
-- Runtime roles (the RBAC role catalog).
-- Created before users: users.role references rbac_roles(name).
-- =============================================
CREATE TABLE rbac_roles (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- People (a single table). A worker is a user with the worker role; the worker's
-- profile (position, dates, manager) is stored directly on users.
CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	last_name TEXT NOT NULL,
	first_name TEXT NOT NULL,
	-- Patronymic is optional.
	middle_name TEXT DEFAULT NULL,
	-- Role from the rbac_roles catalog (admin/dp/rp/vp/worker + runtime roles).
	role TEXT NOT NULL REFERENCES rbac_roles(name),
	username TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	-- The worker's manager (a user with the vp role); NULL for other roles.
	manager_id BIGINT REFERENCES users(id),
	-- Position is free text (not a resource type).
	position TEXT NOT NULL DEFAULT '',
	-- NULL means "on staff since the beginning of time"; before hire_date the worker is not counted.
	hire_date DATE DEFAULT NULL,
	-- NULL means "still employed"; after termination_date the worker is not counted.
	termination_date DATE DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (termination_date IS NULL OR hire_date IS NULL OR termination_date >= hire_date)
);

CREATE TABLE projects (
	id BIGSERIAL PRIMARY KEY,
	owner_id BIGINT REFERENCES users(id),
	code TEXT NOT NULL,
	-- Optional entity color (#RRGGBB, NULL — standard color on the frontend).
	color TEXT,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	priority INTEGER NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (end_date >= start_date),
	CHECK (priority >= 0)
);

CREATE TABLE processes (
	id BIGSERIAL PRIMARY KEY,
	project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	owner_id BIGINT REFERENCES users(id),
	title TEXT NOT NULL,
	-- Optional entity color (#RRGGBB, NULL — standard color on the frontend).
	color TEXT,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	-- Order of the process within its project (unique per project, see V2). 
	-- Ascending display order; a distinct order column (not id) so rows can be
	-- reordered without rewriting their identity.
	sort_order INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (end_date >= start_date)
);

CREATE TABLE tasks (
	id BIGSERIAL PRIMARY KEY,
	process_id BIGINT NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
	owner_id BIGINT REFERENCES users(id),
	title TEXT NOT NULL,
	-- Optional entity color (#RRGGBB, NULL — standard color on the frontend).
	color TEXT,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	-- Order of the task within its process (unique per process, see V2).
	sort_order INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (end_date >= start_date)
);

-- Resource category dictionary (specializations).
CREATE TABLE resources (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	code TEXT NOT NULL,
	-- Optional resource color (#RRGGBB, NULL — standard color on the frontend).
	color TEXT,
	-- Resource owner (user account) is required.
	owner_id BIGINT NOT NULL REFERENCES users(id),
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (length(trim(code)) > 0)
);

-- Employee state dictionary: availability is set by the is_available flag.
CREATE TABLE states (
	id BIGSERIAL PRIMARY KEY,
	code TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	is_available BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Worker states (absence only): one row = interval [start_date, end_date].
-- Missing rows mean present; rows are hard-deleted (audit log).
-- The EXCLUDE constraint forbids overlapping states per worker (one state per day).
CREATE TABLE user_states (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	state_id BIGINT NOT NULL REFERENCES states(id) ON DELETE RESTRICT,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CHECK (end_date >= start_date),
	EXCLUDE USING gist (
		user_id WITH =,
		daterange(start_date, end_date, '[]') WITH &&
	)
);

-- Resource user list: the resource enumerates its workers.
-- UNIQUE(user_id) — a worker belongs to at most one resource; membership is optional.
CREATE TABLE resource_members (
	resource_id BIGINT NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	PRIMARY KEY (resource_id, user_id),
	UNIQUE (user_id)
);

CREATE TABLE assignments (
	id BIGSERIAL PRIMARY KEY,
	task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	resource_id BIGINT NOT NULL REFERENCES resources(id) ON DELETE RESTRICT,
	quantity INTEGER NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (quantity > 0)
);

CREATE TABLE milestones (
	id BIGSERIAL PRIMARY KEY,
	process_id BIGINT NOT NULL REFERENCES processes(id) ON DELETE CASCADE,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	-- Optional entity color (#RRGGBB, NULL — standard color on the frontend).
	color TEXT,
	date DATE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Auto-creation config for processes/tasks on project insert (V8).
CREATE TABLE project_auto_create (
	id BIGSERIAL PRIMARY KEY,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	config JSONB NOT NULL DEFAULT '[]'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Server-side refresh sessions (refresh-token rotation/revocation).
-- Only the SHA-256 of the token is stored; the opaque token lives in an HttpOnly cookie.
CREATE TABLE refresh_sessions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    -- Rotation chain: id of the session this one replaced (to detect reuse).
    replaced_by BIGINT REFERENCES refresh_sessions(id)
);

-- Idempotency keys: idempotent create endpoints.
-- PK (key, user_id, method, path): isolation per client/route.
CREATE TABLE idempotency_keys (
    key             TEXT NOT NULL,
    user_id         BIGINT NOT NULL,
    method          TEXT NOT NULL,
    path            TEXT NOT NULL,
    response_status INTEGER NOT NULL DEFAULT 0,
    response_body   JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (key, user_id, method, path)
);

-- Task comments: discussion threads (replies to replies via parent_id).
CREATE TABLE task_comments (
	id BIGSERIAL PRIMARY KEY,
	task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	-- Comment author (taken from the auth context).
	author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
	-- Reply to another comment on the same task; NULL — a root comment.
	parent_id BIGINT REFERENCES task_comments(id) ON DELETE CASCADE,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- The rights matrix (rbac_role_rules) and route-policy definitions
-- (rbac_route_policies) are configured at runtime: the source of truth is the DB,
-- while the engine (scope interpretation, kind registry) stays in code.
CREATE TABLE rbac_role_rules (
    id          BIGSERIAL PRIMARY KEY,
    role        TEXT NOT NULL REFERENCES rbac_roles(name),
    resource    TEXT NOT NULL, -- project|process|task|milestone|assignment|state|resource|worker|comment|user_catalog|rbac_config
    action      TEXT NOT NULL, -- view|create|update|delete
    scope       TEXT NOT NULL, -- all|own|parent|ancestor
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    updated_by  BIGINT,        -- user who made the last change (from JWT)
    UNIQUE (role, resource, action)
);

CREATE TABLE rbac_route_policies (
    name        TEXT PRIMARY KEY,  -- name routes reference (mw.Check("project.create"))
    kind        TEXT NOT NULL,     -- list|entity|create|owner_match|author_or|parent_action
    params      JSONB NOT NULL,    -- kind params (schema — the kind registry in code)
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    updated_by  BIGINT
);