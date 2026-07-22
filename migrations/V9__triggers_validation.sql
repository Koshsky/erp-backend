-- =============================================
-- Block child date range violations relative to parent boundaries.
-- =============================================
CREATE OR REPLACE FUNCTION fn_processes_validate_within_project_dates()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
	parent_start DATE;
	parent_end DATE;
BEGIN
	SELECT p.start_date, p.end_date
	INTO parent_start, parent_end
	FROM projects p
	WHERE p.id = NEW.project_id
	  AND p.deleted_at IS NULL;

	IF parent_start IS NULL OR parent_end IS NULL THEN
		RAISE EXCEPTION
			'Parent project % not found or deleted',
			NEW.project_id;
	END IF;

	IF NEW.start_date < parent_start OR NEW.end_date > parent_end THEN
		RAISE EXCEPTION
			'Process dates [% - %] must be within project dates [% - %]',
			NEW.start_date, NEW.end_date, parent_start, parent_end;
	END IF;

	RETURN NEW;
END;
$$;

CREATE TRIGGER trg_processes_validate_within_project_dates
BEFORE INSERT OR UPDATE OF start_date, end_date, project_id ON processes
FOR EACH ROW
WHEN (NEW.deleted_at IS NULL)
EXECUTE FUNCTION fn_processes_validate_within_project_dates();

-- =============================================
-- Tasks must stay within Process dates
-- =============================================
CREATE OR REPLACE FUNCTION fn_tasks_validate_within_process_dates()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
	parent_start DATE;
	parent_end DATE;
BEGIN
	SELECT p.start_date, p.end_date
	INTO parent_start, parent_end
	FROM processes p
	WHERE p.id = NEW.process_id
	  AND p.deleted_at IS NULL;

	IF parent_start IS NULL OR parent_end IS NULL THEN
		RAISE EXCEPTION
			'Parent process % not found or deleted',
			NEW.process_id;
	END IF;

	IF NEW.start_date < parent_start OR NEW.end_date > parent_end THEN
		RAISE EXCEPTION
			'Task dates [% - %] must be within process dates [% - %]',
			NEW.start_date, NEW.end_date, parent_start, parent_end;
	END IF;

	RETURN NEW;
END;
$$;

CREATE TRIGGER trg_tasks_validate_within_process_dates
BEFORE INSERT OR UPDATE OF start_date, end_date, process_id ON tasks
FOR EACH ROW
WHEN (NEW.deleted_at IS NULL)
EXECUTE FUNCTION fn_tasks_validate_within_process_dates();
