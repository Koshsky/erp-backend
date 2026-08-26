-- =============================================
-- РОЛИ РАНТАЙМА: users.role -> rbac_roles
-- =============================================
-- CHECK (V3) ограничивал роль пятью константами; каталог ролей с V15
-- конфигурируется в рантайме (GET/PUT /api/v1/rbac/roles не CRUD — каталог
-- прав, но сам каталог хранится в БД). Заменяем CHECK на FK, чтобы новые
-- роли можно было заводить в rbac_roles и назначать пользователям
-- (INSERT rbac_roles + UPDATE users.role); нарушение FK (роль не в каталоге)
-- отдаётся API как 400.

ALTER TABLE users DROP CONSTRAINT users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_fk FOREIGN KEY (role) REFERENCES rbac_roles(name);