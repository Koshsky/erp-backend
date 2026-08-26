-- =============================================
-- VALIDATION BOUNDARIES: явная ошибка вместо тихого soft-delete
-- =============================================
-- Раньше (V5): процесс/задача с датами вне границ родителя «тихо»
-- помечались удалёнными при вставке (NEW.deleted_at := now(); RETURN NEW):
-- API возвращал 201/500 с фантомной записью, которая сразу удалена, а
-- не-пустой deleted_at ломал сканирование `**time.Time` (500).
-- Теперь: явный отказ с кодом 22023 (invalid_parameter_value) и понятным
-- сообщением; репозитории превращают его в 400.

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

	IF NEW.end_date < parent_start OR NEW.start_date > parent_end THEN
		RAISE EXCEPTION 'Сроки процесса (% - %) выходят за границы проекта (% - %)',
			NEW.start_date, NEW.end_date, parent_start, parent_end
		USING ERRCODE = '22023';
	END IF;

	IF NEW.start_date < parent_start THEN
		NEW.start_date := parent_start;
	END IF;
	IF NEW.end_date > parent_end THEN
		NEW.end_date := parent_end;
	END IF;

	RETURN NEW;
END;
$$;

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

	IF NEW.end_date < parent_start OR NEW.start_date > parent_end THEN
		RAISE EXCEPTION 'Сроки задачи (% - %) выходят за границы процесса (% - %)',
			NEW.start_date, NEW.end_date, parent_start, parent_end
		USING ERRCODE = '22023';
	END IF;

	IF NEW.start_date < parent_start THEN
		NEW.start_date := parent_start;
	END IF;
	IF NEW.end_date > parent_end THEN
		NEW.end_date := parent_end;
	END IF;

	RETURN NEW;
END;
$$;