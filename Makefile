# Makefile para el proyecto Hospital API

.PHONY: help build run test clean docker-build docker-run docker-stop deps seed seed-clean db-reset fresh-start dev-with-data

# Variables
BINARY_NAME=hospital-api
DOCKER_IMAGE=hospital-api:latest

# Help
help: ## Muestra esta ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Dependencias
deps: ## Instala y actualiza dependencias
	go mod download
	go mod tidy

# Build
build: deps ## Compila la aplicación
	go build -o bin/$(BINARY_NAME) ./cmd/server

# Run
run: ## Ejecuta la aplicación en modo desarrollo
	go run ./cmd/server/main.go

# Test
test: ## Ejecuta las pruebas
	go test -v ./...

# Clean
clean: ## Limpia archivos de compilación
	rm -rf bin/
	go clean

# Docker
docker-build: ## Construye la imagen Docker
	docker-compose build

docker-run: ## Ejecuta la aplicación con Docker Compose
	docker-compose up -d

docker-stop: ## Detiene los contenedores
	docker-compose down

docker-logs: ## Muestra logs de los contenedores
	docker-compose logs -f

docker-rebuild: ## Reconstruye y ejecuta los contenedores
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Desarrollo
dev: ## Inicia el entorno de desarrollo completo
	docker-compose up -d db
	sleep 5
	go run ./cmd/server/main.go

# Base de datos
db-reset: ## Reinicia la base de datos
	docker-compose down db
	docker volume rm api-go_postgres_data
	docker-compose up -d db

# Seeding
seed: ## Inserta datos de prueba en la base de datos
	go run cmd/seed/main.go

seed-clean: ## Limpia la base de datos e inserta datos frescos
	go run cmd/seed/main.go -clean

# Comandos combinados
fresh-start: db-reset seed ## Reinicia la BD e inserta datos de prueba
	@echo "🎉 Base de datos reiniciada y datos de prueba insertados"

dev-with-data: docker-run ## Inicia la aplicación con datos de prueba automáticos
	@echo "🚀 Aplicación iniciada con seeding automático habilitado"

dev-no-seed: ## Inicia la aplicación SIN datos de prueba automáticos
	AUTO_SEED=false docker-compose up -d
	@echo "🚀 Aplicación iniciada SIN seeding automático"

# Logs
logs: ## Muestra logs de la aplicación
	docker-compose logs -f app

# Producción
prod: ## Ejecuta en modo producción
	GIN_MODE=release go run ./cmd/server/main.go
