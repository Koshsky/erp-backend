-- =============================================
-- RUNTIME RBAC POLICIES
-- =============================================
-- Таблицы rbac_roles / rbac_role_rules / rbac_route_policies — в V1.
-- Матрица прав (rbac_role_rules) и определения маршрутных проверок
-- (rbac_route_policies) конфигурируются в рантайме: источник истины — БД,
-- движок (интерпретация скоупов, реестр kind'ов) остаётся кодом.
--
-- Скоупы (отсутствие строки = нет доступа):
--   all      — полный доступ (dp на просмотрах; vp на process.view/state.view)
--   own      — владелец самой строки (проект для rp; resource/worker для vp)
--   parent   — владелец непосредственного родителя (управление: процесса —
--              проект; задачи/вехи/назначения — процесс)
--   ancestor — любой владелец в цепочке предков (просмотр задач для rp)
-- admin — жёсткий bypass в коде (ScopeAll), в БД не хранится.

-- =============================================
-- SEED: каталог ролей
-- =============================================
INSERT INTO rbac_roles (name, description) VALUES
    ('admin',  'system administrator — full access (bypass in code)'),
    ('dp',     'project portfolio director'),
    ('rp',     'project manager'),
    ('vp',     'process owner'),
    ('worker', 'worker (no rights yet)');

-- =============================================
-- SEED: матрица прав (текущее поведение в новых скоупах)
-- =============================================
INSERT INTO rbac_role_rules (role, resource, action, scope) VALUES
    -- projects: dp — все, rp — свои (просмотр/редактирование/удаление), rp создаёт в свою собственность
    ('dp', 'project', 'view',   'all'),
    ('rp', 'project', 'view',   'own'),
    ('rp', 'project', 'create', 'own'),
    ('dp', 'project', 'update', 'all'),
    ('rp', 'project', 'update', 'own'),
    ('rp', 'project', 'delete', 'own'),
    -- processes: dp — все, rp — процессы своих проектов (parent), vp — справочно все
    ('dp', 'process', 'view',   'all'),
    ('rp', 'process', 'view',   'parent'),
    ('vp', 'process', 'view',   'all'),
    ('rp', 'process', 'create', 'parent'),
    ('rp', 'process', 'update', 'parent'),
    ('rp', 'process', 'delete', 'parent'),
    -- tasks / milestones / assignments: dp — все, rp — просмотр через проект (ancestor), vp — управление в своём процессе (parent)
    ('dp', 'task',   'view',   'all'),
    ('rp', 'task',   'view',   'ancestor'),
    ('vp', 'task',   'view',   'parent'),
    ('vp', 'task',   'create', 'parent'),
    ('vp', 'task',   'update', 'parent'),
    ('vp', 'task',   'delete', 'parent'),
    ('dp', 'milestone',   'view',   'all'),
    ('rp', 'milestone',   'view',   'ancestor'),
    ('vp', 'milestone',   'view',   'parent'),
    ('vp', 'milestone',   'create', 'parent'),
    ('vp', 'milestone',   'update', 'parent'),
    ('vp', 'milestone',   'delete', 'parent'),
    ('dp', 'assignment',  'view',   'all'),
    ('rp', 'assignment',  'view',   'ancestor'),
    ('vp', 'assignment',  'view',   'parent'),
    ('vp', 'assignment',  'create', 'parent'),
    ('vp', 'assignment',  'update', 'parent'),
    ('vp', 'assignment',  'delete', 'parent'),
    -- states (справочник табеля): vp — справочно, управление — только admin (bypass)
    ('vp', 'state', 'view', 'all'),
    -- timesheet resources: vp — свои (own)
    ('vp', 'resource', 'view',   'own'),
    ('vp', 'resource', 'create', 'own'),
    ('vp', 'resource', 'update', 'own'),
    ('vp', 'resource', 'delete', 'own'),
    -- workers: создание сотрудников — только admin (bypass); vp — свои подчинённые (own)
    ('vp', 'worker', 'view',   'own'),
    ('vp', 'worker', 'update', 'own'),
    ('vp', 'worker', 'delete', 'own'),
    -- comments: прав нет в матрице — права считаются по родительской задаче
    -- user_catalog (виртуальный ресурс: каталог пользователей для пикеров): dp/rp/vp + admin (bypass)
    ('dp', 'user_catalog', 'view', 'all'),
    ('rp', 'user_catalog', 'view', 'all'),
    ('vp', 'user_catalog', 'view', 'all');
    -- rbac_config (виртуальный ресурс: управление автосозданием/RBAC): только admin (bypass), строк нет

-- =============================================
-- SEED: маршрутные проверки (имя → kind + параметры)
-- =============================================
INSERT INTO rbac_route_policies (name, kind, params) VALUES
    -- projects
    ('project.list',    'list',   '{"resource":"project","query_key":"owner_id"}'),
    ('project.view',    'entity', '{"resource":"project","action":"view","owner":"id"}'),
    ('project.create',  'create', '{"resource":"project","owner_key":"owner_id","default_self":true}'),
    ('project.update',  'entity', '{"resource":"project","action":"update","owner":"id"}'),
    ('project.delete',  'entity', '{"resource":"project","action":"delete","owner":"id"}'),
    -- processes
    ('process.list',    'list',   '{"resource":"process","query_key":"owner_id"}'),
    ('process.view',    'entity', '{"resource":"process","action":"view","owner":"id"}'),
    ('process.create',  'create', '{"resource":"process","parent_resource":"project","parent_from":"project_id"}'),
    ('process.update',  'entity', '{"resource":"process","action":"update","owner":"id"}'),
    ('process.delete',  'entity', '{"resource":"process","action":"delete","owner":"id"}'),
    -- tasks
    ('task.list',       'list',   '{"resource":"task","query_key":"owner_id"}'),
    ('task.view',       'entity', '{"resource":"task","action":"view","owner":"id"}'),
    ('task.create',     'create', '{"resource":"task","parent_resource":"process","parent_from":"process_id"}'),
    ('task.update',     'entity', '{"resource":"task","action":"update","owner":"id"}'),
    ('task.delete',     'entity', '{"resource":"task","action":"delete","owner":"id"}'),
    -- milestones
    ('milestone.list',    'list',   '{"resource":"milestone","query_key":"owner_id"}'),
    ('milestone.view',    'entity', '{"resource":"milestone","action":"view","owner":"id"}'),
    ('milestone.create',  'create', '{"resource":"milestone","parent_resource":"process","parent_from":"process_id"}'),
    ('milestone.update',  'entity', '{"resource":"milestone","action":"update","owner":"id"}'),
    ('milestone.delete',  'entity', '{"resource":"milestone","action":"delete","owner":"id"}'),
    -- assignments: create — кросс-сущностная проверка владельцев (kind owner_match)
    ('assignment.list',   'list',   '{"resource":"assignment","query_key":"owner_id"}'),
    ('assignment.view',   'entity', '{"resource":"assignment","action":"view","owner":"id"}'),
    ('assignment.create', 'owner_match', '{"resource":"assignment","action":"create","primary_resource":"task","primary_from":"task_id","compare_resource":"resource","compare_from":"resource_id","exempt_roles":["admin"]}'),
    ('assignment.update', 'entity', '{"resource":"assignment","action":"update","owner":"id"}'),
    ('assignment.delete', 'entity', '{"resource":"assignment","action":"delete","owner":"id"}'),
    -- comments: list/create по task.view; delete — автор или право обновления задачи
    ('task.comment.list',   'entity', '{"resource":"task","action":"view","owner":"id"}'),
    ('task.comment.create', 'entity', '{"resource":"task","action":"view","owner":"id"}'),
    ('task.comment.delete', 'author_or', '{"author_resource":"comment","author_id_param":"comment_id","right_resource":"task","right_action":"update"}'),
    -- workers
    ('worker.list',    'list',   '{"resource":"worker","query_key":"manager_id"}'),
    ('worker.view',    'entity', '{"resource":"worker","action":"view","owner":"id"}'),
    ('worker.create',  'create', '{"resource":"worker","owner_key":"manager_id","default_self":true}'),
    ('worker.update',  'entity', '{"resource":"worker","action":"update","owner":"id"}'),
    ('worker.delete',  'entity', '{"resource":"worker","action":"delete","owner":"id"}'),
    -- user catalog (пикер имён): виртуальный ресурс user_catalog
    ('user.picker',    'entity', '{"resource":"user_catalog","action":"view","owner":"none"}'),
    -- timesheet resources + calendar
    ('resource.list',        'list',   '{"resource":"resource","query_key":"owner_id"}'),
    ('resource.view',        'entity', '{"resource":"resource","action":"view","owner":"id"}'),
    ('resource.create',      'create', '{"resource":"resource","owner_key":"owner_id","default_self":true}'),
    ('resource.update',      'entity', '{"resource":"resource","action":"update","owner":"id"}'),
    ('resource.delete',      'entity', '{"resource":"resource","action":"delete","owner":"id"}'),
    ('resource.member-list',   'entity', '{"resource":"resource","action":"view","owner":"id"}'),
    ('resource.member-add',    'entity', '{"resource":"resource","action":"update","owner":"id"}'),
    ('resource.member-remove', 'entity', '{"resource":"resource","action":"update","owner":"id"}'),
    ('calendar.view',      'list',   '{"resource":"resource","query_key":"owner_id"}'),
    -- states (справочник без владельца)
    ('state.list',   'entity', '{"resource":"state","action":"view","owner":"none"}'),
    ('state.view',   'entity', '{"resource":"state","action":"view","owner":"none"}'),
    ('state.create', 'entity', '{"resource":"state","action":"create","owner":"none"}'),
    ('state.update', 'entity', '{"resource":"state","action":"update","owner":"none"}'),
    ('state.delete', 'entity', '{"resource":"state","action":"delete","owner":"none"}'),
    -- auto-create config + RBAC admin API (виртуальный ресурс rbac_config: только admin)
    ('autocreate.list',   'entity', '{"resource":"rbac_config","action":"view","owner":"none"}'),
    ('autocreate.update', 'entity', '{"resource":"rbac_config","action":"update","owner":"none"}'),
    ('rbac.manage',       'entity', '{"resource":"rbac_config","action":"view","owner":"none"}');

-- =============================================
-- BLOCK hard DELETE (та же защита, что и для остальных таблиц, V5)
-- =============================================
CREATE TRIGGER block_hard_delete_on_rbac_roles
BEFORE DELETE ON rbac_roles
FOR EACH ROW EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_rbac_role_rules
BEFORE DELETE ON rbac_role_rules
FOR EACH ROW EXECUTE FUNCTION block_hard_delete();

CREATE TRIGGER block_hard_delete_on_rbac_route_policies
BEFORE DELETE ON rbac_route_policies
FOR EACH ROW EXECUTE FUNCTION block_hard_delete();