-- =============================================
-- STATES COLOR (#RRGGBB, nullable; NULL — standard color on the frontend)
-- =============================================
-- The columns are also declared in V1 (the sqlc schema source), so a fresh
-- database already has them and the ALTERs below are no-ops there; existing
-- deployments get the columns from this migration. states_deleted mirrors
-- the live columns (the archive trigger copies the whole row by column name),
-- so it gains the same column.
ALTER TABLE states ADD COLUMN IF NOT EXISTS color TEXT;
ALTER TABLE states_deleted ADD COLUMN IF NOT EXISTS color TEXT;