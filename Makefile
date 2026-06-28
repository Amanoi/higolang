.PHONY: dev dev-server dev-web build up down logs clean

# Development
dev-server:
	cd server && go run main.go

dev-web:
	cd web && npm run dev

# Build
build:
	docker-compose build

# Docker
up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

# Clean
clean:
	docker-compose down -v
	rm -rf data/

# Help
help:
	@echo "Available commands:"
	@echo "  make dev-server  - Run Go backend locally"
	@echo "  make dev-web     - Run Vue frontend dev server"
	@echo "  make build       - Build Docker images"
	@echo "  make up          - Start all services"
	@echo "  make down        - Stop all services"
	@echo "  make logs        - Follow logs"
	@echo "  make restart     - Restart all services"
	@echo "  make clean       - Stop and remove all data"
