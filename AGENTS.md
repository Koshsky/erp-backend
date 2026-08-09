# AGENTS.md

Go 1.25 / Gin backend for the MVS ERP monorepo (`module github.com/Koshsky/erp-backend`). Sibling `../frontend` (Vue) consumes the OpenAPI contract in `docs/swagger/swagger.yaml`. Code comments in this repo are often in Russian — keep that style.

## Commands
- `make dev` — dev stack: `docker compose up -d db flyway` + the API served locally by `air` (auto-reload on change; binary in `build/`, ignored). `make infra` only starts db+flyway, `make reset` wipes the DB volume (`docker compose down -v`), `make stop` stops the containers. `Makefile` includes `backend/.env` (symlink to the repo-root `.env`).
- `go run ./cmd/service` — run the API server directly (config comes from env vars, see below).
- `golangci-lint run` — lint (v2 golden config, requires golangci-lint ≥ 2.12). `--fix` also runs the formatters (`goimports` + `golines`, max line 120). Import groups must keep local `github.com/Koshsky/erp-backend` last.
- `go test ./...` — tests (currently minimal).
- `swag fmt` then `swag init -g cmd/service/main.go -o docs/swagger` — regenerate swagger after touching handler annotations or global docs.
- `GOTOOLCHAIN=go1.25.0 go run github.com/google/wire/cmd/wire@latest ./internal/server/...` — regenerate `internal/server/wire_gen.go` after touching any provider or `internal/server/wire.go` (commit `wire_gen.go`; `wire.go` carries the `//go:generate` directive).

## Generated code — never hand-edit
- `docs/swagger/` (`docs.go`, `swagger.json`, `swagger.yaml`): regenerate via `swag init` and commit the result. The frontend regenerates its entire API client from `docs/swagger/swagger.yaml`, so API contract changes must land there.
- `internal/*/repository/sqlc/`: per-package sqlc output. After editing `internal/*/repository/query/*.sql`, run `go generate ./...` (each `repository/postgres.go` carries the `//go:generate` directive). Every package has its own `sqlc.yaml` (v2, pgx/v5) generating into `./sqlc`.

## Architecture
- Layered per domain: `internal/<domain>/{delivery,service,repository,domain,dto}`. Structs are constructed by exported `New*` constructors kept in the same file as the struct (`handler.go`/`service_impl.go`/`postgres.go`); there are no `provide.go` files in the layers. Each domain root ships a single `module.go` with a `wire.ProviderSet` referencing those `New*` constructors and a `Module` implementing the `server.Module` interface (`RegisterPublicRoutes`/`RegisterProtectedRoutes`) — the module owns route registration, prefixes and per-route RBAC checks. The app only collects modules via `ProvideModules` in `internal/server/providers.go`; the graph is assembled in `internal/server/wire.go` (generated into `wire_gen.go`). Auxiliary HTTP servers live in `internal/server/*` (`profiler` pprof, `swagger` registrar). New handlers are added to the module's `ProviderSet`; new modules add a `Module` + one line in `ProvideModules`.
- All API responses use the `{ data, error }` envelope via `internal/common/response` helpers (`OK`, `Created`, `BadRequest`, `InternalError`, ...). Never define a local response struct in a handler.
- Calendar dates (`start_date`, `end_date`, `date`) are `internal/common/date.Date` (string `YYYY-MM-DD`), never `time.Time`.
- `internal/auth` must not touch the `users` table; it depends on the `internal/user` service (single source of truth for user CRUD). Password hashing lives only in `internal/security/hasher`.
- Authenticated user comes from request context via `internal/common/ctx` (`KeyUser`/`KeyRequest`), set by `internal/middleware/auth`.

## DB & migrations
- Postgres 16; schema is Flyway migrations in `migrations/` (V1..V10, applied by the `flyway` docker service). Editing existing migration files is allowed (repo policy). Every `sqlc.yaml` references only `V1__initial_schema.sql`, so column-type changes (nullability, types) must land in `V1` for `go generate` to pick them up — otherwise add the migration to every sqlc.yaml schema list.
- Soft delete, audit, and date-shift are enforced by DB triggers (V5–V9): SQL queries must filter `deleted_at IS NULL`.
- Config: non-secret settings live in the committed `config.yaml` (server timeouts, postgres pool, jwt expiry/issuer, logging, swagger, profiling); only secrets + DB URL come from env vars — `DATABASE_URL`, `JWT_SECRET_KEY`, `JWT_REFRESH_KEY`. `CONFIG_PATH` env overrides the file (default `./config.yaml`). Compose mounts `config.yaml` read-only (not baked into the image), so config edits need no rebuild. No `.env` is auto-loaded in code; `backend/.env` is a symlink to the repo-root `.env` and is injected by docker-compose (it also keeps `DATABASE_NAME/USER/PASSWORD` for the postgres+flyway containers — the app ignores them). Running natively: point `DATABASE_URL` at localhost (the committed `.env` uses host `db` for docker); the swagger `@host` annotation stays `localhost:8080`.
- `profiling.enabled: true` in config.yaml starts a pprof server on `profiling.address` (`:6060`), `/debug/pprof/*`.

## Running the stack
- Full stack (nginx → erp + frontend + postgres + flyway) lives in the repo-root `../../docker-compose.yml`.
- Standalone backend dev: `make dev` — db + flyway in docker (port 5432 exposed), API via `air` locally on :8080 with `DATABASE_URL` pointing at `localhost` (`Makefile` overrides the `db` host from `.env`). Requires `air` v1.65.3 (Go 1.25-compatible): `GOTOOLCHAIN=go1.25.0 go install github.com/air-verse/air@v1.65.3`.
- `docker-compose.yml` in this dir also keeps an `erp` service (build from Dockerfile) for an all-in-docker dev run: `docker compose up erp`.
- `cmd/swagger/main.go` is a dev-only standalone Swagger UI server on :8080 (conflicts with the main service). The global swagger annotations are duplicated in `cmd/service/main.go` and `cmd/swagger/main.go` — keep both in sync.

## Style / process
- Conventional commits (`feat/fix/refactor/chore/docs`).
- The golangci config is strict: no package-level vars, no `init()`, no magic numbers, `log/slog` only (stdlib `log` is denied outside `main`), doc comments end with a period (exempt in `delivery/`), function length/complexity caps. Don't add `//nolint` without an explanation.
