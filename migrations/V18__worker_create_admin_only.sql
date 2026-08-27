-- =============================================
-- СОЗДАНИЕ СОТРУДНИКОВ — ТОЛЬКО ДЛЯ АДМИНИСТРАТОРА
-- =============================================
-- Рантайм-матрица RBAC хранится в rbac_role_rules. Для БД, уже
-- промигрированных версией V15 (vp worker/create own), убираем правило
-- мягким удалением (hard-delete заблокирован триггером V7; отсутствие
-- строки = нет доступа в ScopeFor).
-- На чистой БД (V15 уже без этой строки) запись не находится — no-op.
UPDATE rbac_role_rules
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE role = 'vp'
  AND resource = 'worker'
  AND action = 'create'
  AND deleted_at IS NULL;