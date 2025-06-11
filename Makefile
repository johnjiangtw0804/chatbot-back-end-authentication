.EXPORT_ALL_VARIABLES:
### PostgreSQL ENV FOR LOCAL (https://hub.docker.com/_/postgres)###
DB_USER ?= root
DB_PASSWORD ?= root
DB_NAME ?= user
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_LOG_MODE ?= True
PG_CONTAINER_NAME ?= local-postgres

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
