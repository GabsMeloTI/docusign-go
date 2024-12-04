# include
include .env

# Swagger
swag:
	@echo "$(YELLOW) Generating $(CYAN) Swagger $(GREEN). $(NC)"
	swag init --parseDependency --parseInternal -g cmd/api.go --output docs/app

# Local Setup
postgres-setup:
	docker run --name simpplify-postgres -p ${POSTGRES_PORT}:${POSTGRES_PORT}/tcp -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:12-alpine

start-postgres:
	docker start simpplify-postgres

stop-postgres:
	docker stop simpplify-postgres

createdb: 
	docker exec -it simpplify-postgres createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} $(POSTGRES_DB)

# local migrations
migration-up:
	@echo "$(CYAN)Starting local migration... $(GREEN) UP $(NC)"
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose up

migration-down:
	@echo "$(CYAN)Starting local migration... $(RED) DOWN $(NC)"
	migrate -path db/migration -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -verbose down

# run locally
run: 
	@echo "$(GREEN) Running Golang App... $(CYAN)LOCAL $(NC)"
	go run ./main.go local

# internal commands
sqlc:
	@echo "$(YELLOW) Generating $(CYAN) sqlc $(YELLOW)files... $(NC)"
	sqlc generate

.PHONY: run sqlc postgres-setup start-postgres stop-postgres createdb migration-up migration-down swag

run-rabbit:
	go run cmd/rabbit/main.go local

