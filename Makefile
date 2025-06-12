### This Makefile is used for db migration and other local development task
### For deployment, we could simply run db migration with the correct config

include config.env
export

.EXPORT_ALL_VARIABLES:
# 是針對 Makefile「執行中」的指令才有效，而不是讓 shell 全域變數生效
PG_CONTAINER_NAME ?= local-postgres

### GOOSE DB MIGRATION
### Always use := when defining a variable that depends on other variables.
GOOSE_DRIVER ?= postgres
GOOSE_MIGRATION_DIR ?= ./migration
GOOSE_DBSTRING := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)

### @ means Don't echo the command line
### something notice here is that when we have both env and we used the DB string in the argument, this would
### cause an issue
### https://github.com/pressly/goose/issues/239
migrate-up:
	@echo $(GOOSE_DBSTRING)
	@goose -dir $(GOOSE_MIGRATION_DIR) up

run:
	@go run main.go

# ========================
# 🐘 PostgreSQL (Docker)
# ========================
stop-pg:
	@echo "Stopping PostgreSQL container..."
	@docker stop $(PG_CONTAINER_NAME) || true
	@docker rm $(PG_CONTAINER_NAME) || true

start-pg:
	@echo "Starting PostgreSQL container..."
	@docker run -d \
		--name $(PG_CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		postgres:15

print-env-test:
	@echo "ENV: $(APP_ENV), DB: $(DB_NAME)@$(DB_HOST):$(DB_PORT) $(PG_CONTAINER_NAME)"
