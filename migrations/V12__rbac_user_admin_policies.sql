-- =============================================
-- USER-ADMIN ROUTE POLICIES
-- =============================================
-- An employee IS a system user (single users table), so profile mutations
-- (create/update/delete user, set manager, reset password) are gated by the
-- user_admin virtual resource instead of worker.*: only a holder of the
-- user-edit right may edit a user, and that happens through the admin users
-- section. The matrix has no user_admin rows by default (admin bypass), the
-- right is grantable via the RBAC matrix (scope "all" is the only applicable
-- one — the backend validates this).
INSERT INTO rbac_route_policies (name, kind, params, active) VALUES
    ('user_admin.create', 'entity', '{"resource":"user_admin","action":"create","owner":"none"}', TRUE),
    ('user_admin.update', 'entity', '{"resource":"user_admin","action":"update","owner":"none"}', TRUE),
    ('user_admin.delete', 'entity', '{"resource":"user_admin","action":"delete","owner":"none"}', TRUE)
ON CONFLICT (name) DO NOTHING;