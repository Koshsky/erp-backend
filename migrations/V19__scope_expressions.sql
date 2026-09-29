-- =============================================
-- SCOPE EXPRESSIONS (OWNERSHIP-TREE VISIBILITY)
-- =============================================
-- The ownership scope of a rule is now an EXPRESSION over the ownership tree
-- (all | self | up1 | up | sib | down | combinations like "self sib");
-- legacy zone codes (own|parent|ancestor) are canonicalized on save. The
-- expression syntax is validated by the RBAC service (engine scopeexpr.go);
-- the DB only guards against empty values.
--
-- Existing deployments drop the legacy-only CHECK; a fresh database already
-- has the new CHECK inline in V1 (this ALTER is a no-op there — guarded).

DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'user_permissions_check'
    ) THEN
        ALTER TABLE user_permissions DROP CONSTRAINT user_permissions_check;
    END IF;
END $$;

ALTER TABLE user_permissions ADD CONSTRAINT user_permissions_check
CHECK (NOT granted OR length(trim(scope)) > 0);