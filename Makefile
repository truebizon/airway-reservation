docker-build:
	@docker compose -f docker-compose.yml build

docker-up:
	@docker compose -f docker-compose.yml up --build

docker-local-up:
	@docker compose -f docker-compose.local.yml up --build

docker-down:
	@docker compose -f docker-compose.yml down

docker-local-down:
	@docker compose -f docker-compose.local.yml down

migrate-%:
	go run ./cmd/migration/main.go schema apply -r ./database/migration/${@:migrate-%=%} -p 5432 --dbname postgres

.PHONY: migrate
migrate:
	make migrate-airway_reservation
