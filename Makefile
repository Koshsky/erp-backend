include .env
export

# Dev stack: infra (db + flyway) in docker, API served locally by air.
.PHONY: dev infra reset stop

dev:
	docker compose up -d db flyway
	DATABASE_URL=postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable \
	air

infra:
	docker compose up -d db flyway

reset:
	docker compose down -v
	docker compose up -d db flyway

stop:
	docker compose down
