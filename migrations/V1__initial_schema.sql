-- Initial schema for enterprise resource planning.
-- Полная схема: все таблицы создаются один раз, все CHECK/FK-ограничения
-- объявлены инлайн. Никаких последующих ALTER TABLE не требуется.

-- EXCLUDE-констрейнт user_states опирается на btree_gist.
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- =============================================
-- Роли рантайма (каталог ролей RBAC).
-- Создаётся до users: users.role ссылается на rbac_roles(name).
-- =============================================
CREATE TABLE rbac_roles (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

-- Люди (единая таблица). Рабочий — это пользователь с ролью worker; профиль
-- рабочего (должность, даты, руководитель) хранится прямо на users.
CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	last_name TEXT NOT NULL,
	first_name TEXT NOT NULL,
	-- Отчество необязательное.
	middle_name TEXT DEFAULT NULL,
	-- Роль из каталога rbac_roles (admin/dp/rp/vp/worker + рантайм-роли).
	role TEXT NOT NULL REFERENCES rbac_roles(name),
	username TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	-- Руководитель рабочего (user с ролью vp); для остальных ролей — NULL.
	manager_id BIGINT REFERENCES users(id),
	-- Должность — свободный текст (не тип ресурса).
	position TEXT NOT NULL DEFAULT '',
	-- NULL означает "в штате с начала времён"; до hire_date рабочий не учитывается.
	hire_date DATE DEFAULT NULL,
	-- NULL означает "работает поныне"; после termination_date рабочий не учитывается.
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
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
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
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
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

-- Состояния рабочего (только неприсутствие): одна строка = интервал [start_date, end_date].
-- Отсутствие строк — присутствие. Жёстко удаляются (журнал).
-- EXCLUDE запрещает пересечение состояний одного рабочего (одно состояние в день).
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

-- Перечень пользователей ресурса: у ресурса перечислены его работники.
-- UNIQUE(user_id) — рабочий входит не более чем в один ресурс; членство необязательно.
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
	date DATE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Конфигурация автосоздания процессов/задач при вставке проекта (V8).
CREATE TABLE project_auto_create (
	id BIGSERIAL PRIMARY KEY,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	config JSONB NOT NULL DEFAULT '[]'::jsonb,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Серверные refresh-сессии (ротация/отзыв refresh-токенов).
-- Храним только SHA-256 токена; сам opaque-токен живёт в HttpOnly-куке.
CREATE TABLE refresh_sessions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    -- Цепочка ротации: id сессии, которой эта была заменена (для детекта повторного использования).
    replaced_by BIGINT REFERENCES refresh_sessions(id)
);

-- Idempotency keys: идемпотентные create-эндпоинты.
-- PK (key, user_id, method, path): изоляция по клиенту/маршруту.
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

-- Комментарии к задачам: цепочки обсуждений (ответы на ответы — через parent_id).
CREATE TABLE task_comments (
	id BIGSERIAL PRIMARY KEY,
	task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	-- Отправитель комментария (автор берётся из контекста авторизации).
	author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
	-- Ответ на другой комментарий той же задачи; NULL — корневой комментарий.
	parent_id BIGINT REFERENCES task_comments(id) ON DELETE CASCADE,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Матрица прав (rbac_role_rules) и определения маршрутных проверок
-- (rbac_route_policies) конфигурируются в рантайме: источник истины — БД,
-- движок (интерпретация скоупов, реестр kind'ов) остаётся кодом.
CREATE TABLE rbac_role_rules (
    id          BIGSERIAL PRIMARY KEY,
    role        TEXT NOT NULL REFERENCES rbac_roles(name),
    resource    TEXT NOT NULL, -- project|process|task|milestone|assignment|state|resource|worker|comment|user_catalog|rbac_config
    action      TEXT NOT NULL, -- view|create|update|delete
    scope       TEXT NOT NULL, -- all|own|parent|ancestor
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    updated_by  BIGINT,        -- пользователь, внёсший последнее изменение (из JWT)
    UNIQUE (role, resource, action)
);

CREATE TABLE rbac_route_policies (
    name        TEXT PRIMARY KEY,  -- имя, на которое ссылаются маршруты (mw.Check("project.create"))
    kind        TEXT NOT NULL,     -- list|entity|create|owner_match|author_or
    params      JSONB NOT NULL,    -- параметры kind'а (схема — реестр kind'ов в коде)
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    updated_by  BIGINT
);