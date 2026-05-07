.PHONY: infra-up infra-down run-api run-worker migrate

infra-up:
	cd deploy && docker compose up -d

infra-down:
	cd deploy && docker compose down

run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

migrate:
	docker exec -i newsletter_postgres psql -U newsletter -d newsletter_db < migrations/001_create_subscribers.sql
	docker exec -i newsletter_postgres psql -U newsletter -d newsletter_db < migrations/002_create_newsletter_sends.sql