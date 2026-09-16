include .env
export

# Dev stack: infra (db + flyway + jaeger + loki + grafana) in docker, API
# served locally by air. TRACING_ENDPOINT points air at the in-docker Jaeger
# collector via its published port (the docker-service name "jaeger" is not
# resolvable from the host; config.yaml keeps "jaeger:4317" for the
# all-in-docker run). Loki is published on 3100 so air (audit.url in
# config.yaml) reaches it at localhost:3100.
.PHONY: dev infra reset stop lint

dev:
	docker compose up -d db flyway jaeger loki grafana
	DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable \
	TRACING_ENDPOINT=localhost:4317 \
	air

# Lint with the go.mod toolchain pinned: golangci-lint v2 is built with an
# older Go, so against a newer system Go it fails to type-check the standard
# library (math/rand/v2 "method must have no type parameters" / panic).
lint:
	GOTOOLCHAIN=go1.25.14 golangci-lint run

reset:
	docker compose down -v
	docker compose up -d db flyway jaeger loki grafana

stop:
	docker compose down
