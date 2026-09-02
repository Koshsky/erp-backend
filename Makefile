include .env
export

# Dev stack: infra (db + flyway + jaeger + loki + grafana) in docker, API
# served locally by air. TRACING_ENDPOINT points air at the in-docker Jaeger
# collector via its published port (the docker-service name "jaeger" is not
# resolvable from the host; config.yaml keeps "jaeger:4317" for the
# all-in-docker run). Loki is published on 3100 so air (audit.url in
# config.yaml) reaches it at localhost:3100.
.PHONY: dev infra reset stop

dev:
	docker compose up -d db flyway jaeger loki grafana
	DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable \
	TRACING_ENDPOINT=localhost:4317 \
	air

reset:
	docker compose down -v
	docker compose up -d db flyway jaeger loki grafana

stop:
	docker compose down
