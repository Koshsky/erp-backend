-- Initial schema for enterprise resource planning

CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	role TEXT NOT NULL,
	username TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
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

-- Справочник категорий ресурсов (специализаций).
CREATE TABLE resources (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
    code TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Справочник состояний сотрудника: доступность задаётся флагом is_available.
CREATE TABLE states (
	id BIGSERIAL PRIMARY KEY,
	code TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	is_available BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Конкретный сотрудник (уникальный ресурс) категории resources.
CREATE TABLE employees (
	id BIGSERIAL PRIMARY KEY,
	resource_id BIGINT NOT NULL REFERENCES resources(id) ON DELETE RESTRICT,
	name TEXT NOT NULL,
	-- Должность — свободный текст (не тип ресурса).
	position TEXT NOT NULL DEFAULT '',
	-- Руководитель сотрудника (аккаунт пользователя); NULL — подчинённости нет.
	manager_id BIGINT REFERENCES users(id),
	-- NULL означает "в штате с начала времён"; до hire_date ресурс не учитывается.
	hire_date DATE DEFAULT NULL,
	-- NULL означает "работает поныне"; после termination_date ресурс не учитывается.
	termination_date DATE DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL,
	CHECK (termination_date IS NULL OR hire_date IS NULL OR termination_date >= hire_date)
);

-- Диапазоны состояний сотрудника (только не-явка): одна строка = интервал [start_date, end_date].
-- Отсутствие записей = явка (present). Удаляются жёстко (журнал).
-- EXCLUDE запрещает пересечение состояний одного сотрудника (одно состояние в день).
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE employee_states (
	id BIGSERIAL PRIMARY KEY,
	employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
	state_id BIGINT NOT NULL REFERENCES states(id) ON DELETE RESTRICT,
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CHECK (end_date >= start_date),
	EXCLUDE USING gist (
		employee_id WITH =,
		daterange(start_date, end_date, '[]') WITH &&
	)
);

CREATE INDEX idx_employee_states_dates ON employee_states (start_date, end_date);

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

