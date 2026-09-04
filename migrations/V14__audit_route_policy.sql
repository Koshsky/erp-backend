-- =============================================
-- AUDIT-LOG ROUTE POLICY
-- =============================================
-- The /audit/events endpoint (admin journal UI) is gated by the audit virtual
-- resource. The matrix has no audit rows by default (admin bypass in code),
-- mirroring rbac_config / user_admin / state_admin.
INSERT INTO rbac_route_policies (name, kind, params, active) VALUES
    ('audit.view', 'entity', '{"resource":"audit","action":"view","owner":"none"}', TRUE)
ON CONFLICT (name) DO NOTHING;