-- =============================================
-- CASCADE SOFT DELETE FUNCTIONS
-- =============================================
-- for projects
CREATE OR REPLACE FUNCTION cascade_soft_delete_processes()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE processes
    SET deleted_at = NOW()
    WHERE project_id = OLD.id
      AND deleted_at IS NULL;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- for processes
CREATE OR REPLACE FUNCTION cascade_soft_delete_tasks_and_milestones()
RETURNS TRIGGER AS $$
BEGIN
    -- soft delete all tasks associated with the process
    UPDATE tasks
    SET deleted_at = NOW()
    WHERE process_id = OLD.id
      AND deleted_at IS NULL;

    -- soft delete all milestones
    UPDATE milestones
    SET deleted_at = NOW()
    WHERE process_id = OLD.id
      AND deleted_at IS NULL;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- for tasks
CREATE OR REPLACE FUNCTION cascade_soft_delete_assignments()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE assignments
    SET deleted_at = NOW()
    WHERE task_id = OLD.id
      AND deleted_at IS NULL;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- for resources: мягкое удаление категории мягко удаляет её сотрудников
CREATE OR REPLACE FUNCTION cascade_soft_delete_employees()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE employees
    SET deleted_at = NOW()
    WHERE resource_id = OLD.id
      AND deleted_at IS NULL;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

-- for employees: мягкое удаление сотрудника жёстко удаляет его состояния (интервалы)
CREATE OR REPLACE FUNCTION cascade_delete_employee_states()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM employee_states
    WHERE employee_id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;


-- =============================================
-- TRIGGER 1: Projects → Processes
-- =============================================
CREATE TRIGGER trigger_cascade_soft_delete_processes
AFTER UPDATE OF deleted_at ON projects
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION cascade_soft_delete_processes();

-- =============================================
-- TRIGGER 2: Processes → Tasks & Milestones
-- =============================================
CREATE TRIGGER trigger_cascade_soft_delete_tasks_and_milestones
AFTER UPDATE OF deleted_at ON processes
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION cascade_soft_delete_tasks_and_milestones();

-- =============================================
-- TRIGGER 3: Tasks → Assignments
-- =============================================
CREATE TRIGGER trigger_cascade_soft_delete_assignments
AFTER UPDATE OF deleted_at ON tasks
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION cascade_soft_delete_assignments();

-- =============================================
-- TRIGGER 4: Resources → Employees
-- =============================================
CREATE TRIGGER trigger_cascade_soft_delete_employees
AFTER UPDATE OF deleted_at ON resources
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION cascade_soft_delete_employees();

-- =============================================
-- TRIGGER 5: Employees → Employee States
-- =============================================
CREATE TRIGGER trigger_cascade_delete_employee_states
AFTER UPDATE OF deleted_at ON employees
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION cascade_delete_employee_states();

-- =============================================
-- BLOCK hards DELETE
-- =============================================
CREATE OR REPLACE FUNCTION block_hard_delete()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Hard DELETE is disabled. Use soft delete (UPDATE deleted_at = NOW()) instead.';
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER block_hard_delete_on_users
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_projects
BEFORE DELETE ON projects
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_processes
BEFORE DELETE ON processes
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_tasks
BEFORE DELETE ON tasks
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_resources
BEFORE DELETE ON resources
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_assignments
BEFORE DELETE ON assignments
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_milestones
BEFORE DELETE ON milestones
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_employees
BEFORE DELETE ON employees
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();