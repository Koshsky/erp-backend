# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Admin password reset now returns the generated password in the response (shown once).
- States support a custom color (`#RRGGBB`) on create/update.

## [1.0.1] - 2026-09-30

### Fixed

- Seed migrations no longer desync serial sequences (`V1017__sync_sequences.sql` re-syncs them); creating users/resources no longer fails with a primary-key 409.
- User deletion no longer fails when a stored auto-create template has a malformed owner id.

### Changed

- Auto-generated usernames use the `surname.initials` format.
- User deletion returns a 409 conflict listing the referencing records.

## [1.0.0] - 2026-09-29

### Added

- Initial production release: layered Go/Gin service, Casbin RBAC/ABAC, JWT auth, planning/timesheet/project domains, archive-based soft delete, audit events to Loki, OpenTelemetry tracing.

### Changed

- Seed migrations moved to `migrations/seed/` (fixed explicit ids, idempotent).