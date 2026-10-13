# Migrations

Flyway migrations, applied by the `flyway` docker service (`make infra` / the
deploy stack). Ranges: **core V1+**, **plugins V900+** (deployment-specific
add-ons, untracked), **seed V1010+** (base/demo data, untracked — see
`seed/README.md`).

## Core set

| File | Content |
|---|---|
| `V1__initial_schema.sql` | the whole schema: 33 tables (incl. the `*_deleted` archives), all CHECK/FK constraints inline. **The only schema source for sqlc** — every `sqlc.yaml` points at this file, so column changes (types, nullability) must land here |
| `V2__indexes.sql` | all indexes: unique (username, project code, resource title/code, assignment pair, per-group sort order, per-user permission) and lookups (FKs, dates, sessions/idempotency/comments, task dependencies) |
| `V3__triggers_validation.sql` | process/task/milestone dates must stay inside the parent bounds (milestones are clamped, the rest raise 22023 → HTTP 400) |
| `V4__triggers_date_shift.sql` | shifting a parent range inward shifts children, preserving their duration |
| `V5__triggers_soft_delete.sql` | `fn_archive_row` + one archive trigger per soft-deletable table: a DELETE copies the row into `<table>_deleted` and then really removes it |
| `V6__normalize_user_states.sql` | `fn_normalize_user_states()` — merges touching/overlapping timesheet ranges (called from the user repository) |
| `V7__prevent_manager_cycles.sql` | rejects self-management and cycles in `users.manager_id` |
| `V8__project_auto_create.sql` | the project-creation trigger: clones the `project_auto_create` template into processes, tasks, assignments and subtasks (operations) |
| `V9__rbac_policies.sql` | the RBAC catalog: 5 presets, the 41-row rights matrix, 64 route policies (the baseline the authz PolicyStore loads; the Go engine mirrors it for `/rbac/reset` and as fallback) |

## Rules

- **The core set is frozen once a database is deployed.** Editing an applied
  migration fails `flyway validate` everywhere it ran (checksum mismatch), and
  there is no ignore flag for it — Flyway 9 only knows `missing`, `pending`,
  `ignored` and `future` states. Changes therefore land as a **new** file,
  never as an edit of an existing one.
- This is why the set was rewritten wholesale (21 files → 9) while no database
  held data: dead files, no-op upgrades and duplicated DDL/data were dropped.
  From the first deployment on, that shortcut is gone.
- Keep DDL idempotent-friendly (`IF NOT EXISTS` / `ON CONFLICT`) and never
  duplicate an object another file already creates — a fresh database runs
  every file in order.
- Seeds (`seed/`) insert explicit ids and must `setval()` the sequences they
  touch; core migrations must stay seed-independent.
