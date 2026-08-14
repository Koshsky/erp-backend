-- Initial schema for enterprise resource planning

-- Люди (единая таблица). Рабочий — это пользователь с ролью worker; профиль
-- рабочего (должность, даты, руководитель) хранится прямо на users.
CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	role TEXT NOT NULL,
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
	CHECK (end_date >= start_date)
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
	deleted_at TIMESTAMPTZ DEFAULT NULL
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
	deleted_at TIMESTAMPTZ DEFAULT NULL
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
	deleted_at TIMESTAMPTZ DEFAULT NULL
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
CREATE EXTENSION IF NOT EXISTS btree_gist;

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

CREATE INDEX idx_user_states_dates ON user_states (start_date, end_date);

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
	deleted_at TIMESTAMPTZ DEFAULT NULL
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

