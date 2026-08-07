-- Должность сотрудника — свободный текст, не тип ресурса.
-- Для уже существующих БД (на свежих колонка есть из V1).
ALTER TABLE employees
	ADD COLUMN IF NOT EXISTS position TEXT NOT NULL DEFAULT '';
