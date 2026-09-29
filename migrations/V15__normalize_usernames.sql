-- Normalize existing logins to lowercase so the case-insensitive login flow
-- (and the strict username rule) matches the stored values. Since Flyway
-- versioned migrations run once, a duplicate-after-lowercase would fail the
-- migration loudly; the service layer now lowercases on every write.
UPDATE users SET username = LOWER(username) WHERE username <> LOWER(username);