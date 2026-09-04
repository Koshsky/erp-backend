-- =============================================
-- ENTITY COLORS (#RRGGBB, nullable; NULL — standard color on the frontend)
-- =============================================
-- The columns are also declared in V1 (the sqlc schema source), so a fresh
-- database already has them and the ALTERs below are no-ops there; existing
-- deployments get the columns from this migration.
ALTER TABLE projects    ADD COLUMN IF NOT EXISTS color TEXT;
ALTER TABLE processes   ADD COLUMN IF NOT EXISTS color TEXT;
ALTER TABLE tasks       ADD COLUMN IF NOT EXISTS color TEXT;
ALTER TABLE milestones  ADD COLUMN IF NOT EXISTS color TEXT;
ALTER TABLE resources   ADD COLUMN IF NOT EXISTS color TEXT;

-- =============================================
-- AUTO-CREATE TRIGGER: carry template colors onto created rows
-- =============================================
-- Replaces the V8 function: the JSONB template now may carry an optional
-- "color" (#RRGGBB, hex string) on processes and tasks, which is written
-- into the created rows. Missing/empty color → NULL (standard color).
CREATE OR REPLACE FUNCTION fn_project_auto_create() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
	cfg      RECORD;
	p        RECORD;
	tk       RECORD;
	pid      BIGINT;
	tid      BIGINT;
BEGIN
	SELECT * INTO cfg FROM project_auto_create ORDER BY id LIMIT 1;
	IF NOT FOUND OR NOT cfg.enabled THEN
		RETURN NEW;
	END IF;

	FOR p IN
		SELECT value->>'title' AS title,
		       (value->>'owner_id')::bigint AS owner_id,
		       value->>'color' AS color,
		       value->'tasks' AS tasks,
		       ord
		FROM jsonb_array_elements(cfg.config) WITH ORDINALITY AS pp(value, ord)
	LOOP
		INSERT INTO processes (project_id, title, start_date, end_date, owner_id, color, sort_order)
		VALUES (NEW.id, p.title, NEW.start_date, NEW.end_date, p.owner_id, NULLIF(p.color, ''), p.ord)
		RETURNING id INTO pid;

		FOR tk IN
			SELECT value, ord
			FROM jsonb_array_elements(p.tasks) WITH ORDINALITY AS tt(value, ord)
		LOOP
			INSERT INTO tasks (process_id, title, start_date, end_date, color, sort_order)
			VALUES (pid, tk.value->>'title', NEW.start_date, NEW.end_date, NULLIF(tk.value->>'color', ''), tk.ord)
			RETURNING id INTO tid;

			INSERT INTO assignments (task_id, resource_id, quantity)
			SELECT tid, (res.value->>'resource_id')::bigint, (res.value->>'quantity')::int
			FROM jsonb_array_elements(tk.value->'resources') res
			WHERE (res.value->>'resource_id')::bigint IS NOT NULL;
		END LOOP;
	END LOOP;

	RETURN NEW;
END;
$$;