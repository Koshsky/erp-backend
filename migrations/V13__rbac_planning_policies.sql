-- =============================================
-- PLANNING AGGREGATE ROUTE POLICIES
-- =============================================
-- The /planning/* endpoints return scoped aggregates (projects / processes /
-- tasks). They are gated by the view matrix of the underlying domain through
-- the list kind: a role without the view right gets 403 instead of an empty
-- (or leaking) response; row scoping stays in the SQL (ViewScopeCode).
INSERT INTO rbac_route_policies (name, kind, params, active) VALUES
    ('planning.projects', 'list', '{"resource":"project","query_key":"owner_id"}', TRUE),
    ('planning.processes', 'list', '{"resource":"process","query_key":"owner_id"}', TRUE),
    ('planning.tasks', 'list', '{"resource":"task","query_key":"owner_id"}', TRUE)
ON CONFLICT (name) DO NOTHING;