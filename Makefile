# ============================================
# VARIABLES
# ============================================
COMPOSE_FILE := docker-compose.yaml
ENV_FILE := .env

ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

# Путь к миграциям (относительно корня проекта)
MIGRATIONS_PATH := /migrations
USER_MIGRATIONS_HOST_PATH := ./user-service/internal/repository/db/migrations

# Подключение к БД (из .env)
DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable

# Сеть Docker (из docker-compose)
NETWORK_NAME := pic-by-text_microservices


.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=proto proto/user.proto
	@echo "Proto files generated"



.PHONY: run-user-service
run-user-service:
	cd user-service && go run cmd/app/main.go

.PHONY: run-api-gateway
run-api-gateway:
	cd api-gateway && go run cmd/app/main.go


.PHONY: up
up:
	docker compose -f $(COMPOSE_FILE) up -d

.PHONY: down
down:
	docker compose -f $(COMPOSE_FILE) down

.PHONY: down-volume
down-volume:
	docker compose -f $(COMPOSE_FILE) down -v

.PHONY: stop
stop:
	docker compose -f $(COMPOSE_FILE) stop

.PHONY: logs
logs:
	docker compose -f $(COMPOSE_FILE) logs -f

.PHONY: build
build:
	docker compose -f $(COMPOSE_FILE) build

.PHONY: rebuild
rebuild: down-volume build up
	@echo "Rebuilt completely"


.PHONY: migrate-create
migrate-create:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
	@echo "Created migration: $(name)"

.PHONY: migrate-up
migrate-up:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up
	@echo "Migrations applied"

.PHONY: migrate-down
migrate-down:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1
	@echo "Last migration rolled back"

.PHONY: migrate-down-all
migrate-down-all:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down -all
	@echo "All migrations rolled back"

.PHONY: migrate-force
migrate-force:
	@if "$(version)"=="" echo "Ошибка: укажите версию: make migrate-force version=20250101120000" && exit 1
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)
	@echo "Forced version to $(version)"

.PHONY: migrate-version
migrate-version:
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

.PHONY: full-reset
full-reset: down-volume up migrate-up
