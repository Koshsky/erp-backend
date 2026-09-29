# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-09-29

### Added

- Initial production release of the MVS ERP backend:
  - Go 1.27 / Gin service with layered per-domain architecture
    (`internal/<domain>/{delivery,service,repository,domain,dto}`).
  - Casbin-based RBAC/ABAC access control (`internal/authz`): presets, per-user
    grants, ownership-tree scope expressions, record-level view scoping.
  - Auth with JWT access tokens + opaque refresh sessions, idempotent
    mutations, per-IP and per-user rate limiting (Redis-backed).
  - Planning/timesheet/project management domains (projects, processes, tasks,
    dependencies, milestones, resources, worker states, auto-creation
    templates), archive-based soft delete, task comments.
  - Audit events pushed to Grafana Loki; OpenTelemetry tracing (Jaeger).
- Versioning: SemVer `v1.0.0` tag on `main`; the binary is stamped at build
  time via `-ldflags "-X main.version=…"` (see `Dockerfile`, `scripts/deploy.sh`)
  and the version is logged at startup.

### Changed

- Migrations: data seeds moved from `migrations/plugins/` to
  `migrations/seed/`; seeds use fixed explicit ids (never usernames) and are
  idempotent; flyway runs with `-ignoreMigrationPatterns=*:missing`.
- User deletion returns a 409 conflict listing the referencing records
  (managees, resources, projects, processes, tasks, comments, auto-create
  templates) instead of a generic constraint error.