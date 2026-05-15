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
PICTURE_MIGRATIONS_HOST_PATH := ./picture-service/internal/repository/db/migrations

# Подключение к БД (из .env)
DATABASE_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable

# Сеть Docker (из docker-compose)
NETWORK_NAME := pic-by-text_microservices


.PHONY: proto
proto:
	protoc --go_out=proto --go-grpc_out=proto proto/user.proto
	protoc --go_out=proto --go-grpc_out=proto proto/picture.proto
	@echo "Proto files generated"



.PHONY: run-user-service
run-user-service:
	cd user-service && go run cmd/app/main.go

.PHONY: run-picture-service
run-picture-service:
	cd picture-service && go run cmd/app/main.go

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


# ============================================
# USER SERVICE MIGRATIONS
# ============================================

.PHONY: user-migrate-create
user-migrate-create:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
	@echo "Created user-service migration: $(name)"

.PHONY: user-migrate-up
user-migrate-up:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up
	@echo "User-service migrations applied"

.PHONY: user-migrate-down
user-migrate-down:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1
	@echo "Last user-service migration rolled back"

.PHONY: user-migrate-down-all
user-migrate-down-all:
	@if not exist "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(USER_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down -all
	@echo "All user-service migrations rolled back"

.PHONY: user-migrate-force
user-migrate-force:
	@if "$(version)"=="" echo "Ошибка: укажите версию: make user-migrate-force version=20250101120000" && exit 1
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)
	@echo "Forced user-service version to $(version)"

.PHONY: user-migrate-version
user-migrate-version:
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(USER_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version


# ============================================
# PICTURE SERVICE MIGRATIONS
# ============================================

.PHONY: picture-migrate-create
picture-migrate-create:
	@if not exist "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))"
	docker run --rm -v "$(subst \,/,$(PICTURE_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
	@echo "Created picture-service migration: $(name)"

.PHONY: picture-migrate-up
picture-migrate-up:
	@if not exist "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(PICTURE_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up
	@echo "Picture-service migrations applied"

.PHONY: picture-migrate-down
picture-migrate-down:
	@if not exist "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(PICTURE_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1
	@echo "Last picture-service migration rolled back"

.PHONY: picture-migrate-down-all
picture-migrate-down-all:
	@if not exist "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))" mkdir "$(subst /,\,$(PICTURE_MIGRATIONS_HOST_PATH))"
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(PICTURE_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down -all
	@echo "All picture-service migrations rolled back"

.PHONY: picture-migrate-force
picture-migrate-force:
	@if "$(version)"=="" echo "Ошибка: укажите версию: make picture-migrate-force version=20250101120000" && exit 1
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(PICTURE_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(version)
	@echo "Forced picture-service version to $(version)"

.PHONY: picture-migrate-version
picture-migrate-version:
	docker run --rm --network $(NETWORK_NAME) -v "$(subst \,/,$(PICTURE_MIGRATIONS_HOST_PATH)):$(MIGRATIONS_PATH)" migrate/migrate:v4.18.1 -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version


# ============================================
# ALL MIGRATIONS (BOTH SERVICES)
# ============================================

.PHONY: migrate-create
migrate-create:
	@echo "Creating migrations for both services..."
	@$(MAKE) user-migrate-create name=$(name)
	@$(MAKE) picture-migrate-create name=$(name)

.PHONY: migrate-up
migrate-up:
	@echo "Applying all migrations..."
	@$(MAKE) user-migrate-up
	@$(MAKE) picture-migrate-up

.PHONY: migrate-down
migrate-down:
	@echo "Rolling back last migration for both services..."
	@$(MAKE) user-migrate-down
	@$(MAKE) picture-migrate-down

.PHONY: migrate-down-all
migrate-down-all:
	@echo "Rolling back ALL migrations..."
	@$(MAKE) user-migrate-down-all
	@$(MAKE) picture-migrate-down-all

.PHONY: migrate-version
migrate-version:
	@echo "=== User Service ==="
	@$(MAKE) user-migrate-version
	@echo ""
	@echo "=== Picture Service ==="
	@$(MAKE) picture-migrate-version

.PHONY: full-reset
full-reset: down-volume up migrate-up
	@echo "Full reset completed"

.PHONY: migrate-force
migrate-force:
	@echo "⚠️  Use service-specific commands instead:"
	@echo "  make user-migrate-force version=VERSION"
	@echo "  make picture-migrate-force version=VERSION"