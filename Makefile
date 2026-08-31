include .env
export

# Dev stack: infra (db + flyway + jaeger) in docker, API served locally by air.
# TRACING_ENDPOINT points air at the in-docker Jaeger collector via its
# published port (the docker-service name "jaeger" is not resolvable from the
# host; config.yaml keeps "jaeger:4317" for the all-in-docker run).
.PHONY: dev infra reset stop

dev:
	docker compose up -d db flyway jaeger
	DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable \
	TRACING_ENDPOINT=localhost:4317 \
	air

reset:
	docker compose down -v
	docker compose up -d db flyway jaeger

stop:
	docker compose down
