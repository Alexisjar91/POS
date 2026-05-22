.PHONY: dev build docker-up docker-down db-up db-down test clean runtest

# Variables
APP_NAME := pos-api
CMD_DIR := ./cmd/server
BIN_DIR := ./bin

# === Local Development ===

# Levanta DB en Docker y corre la app localmente con go run
dev: db-up
	@echo "Running app locally..."
	go run $(CMD_DIR)

# Construye el binario localmente
build:
	@echo "Building binary..."
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)

# === Docker ===

# Levanta todo en Docker (producción-like)
docker-up:
	docker compose up -d --build

# Detiene todo
docker-down:
	docker compose down

# Solo la base de datos (para desarrollo local) — espera healthcheck
db-up:
	docker compose up -d --wait postgres

# Detiene solo la DB
db-down:
	docker compose stop postgres

# === Utilidades ===

# Ejecuta tests unitarios
test:
	go test ./... -v -count=1

# Ejecuta test de integración contra BD real
runtest: db-up
	@echo "Running integration tests..."
	@go run ./cmd/runtest/

# Limpia binarios
clean:
	rm -rf $(BIN_DIR)

# Mostrar logs de la app en Docker
logs:
	docker compose logs -f app
