.PHONY: dev prod build run stop clean

# Development environment
dev:
	docker-compose up

# Run in detached mode
dev-d:
	docker-compose up -d

# Production environment
prod:
	docker-compose -f docker-compose.prod.yml up -d

# Build the service
build:
	docker-compose build

# Stop the service
stop:
	docker-compose down

# Stop production service
stop-prod:
	docker-compose -f docker-compose.prod.yml down

# Clean up volumes
clean:
	docker-compose down -v

# Clean up production volumes
clean-prod:
	docker-compose -f docker-compose.prod.yml down -v

# Open a shell in the development container
shell:
	docker-compose exec quiz-service sh

# View logs
logs:
	docker-compose logs -f

# Run database migrations manually (if needed)
migrate:
	docker-compose exec quiz-service go run cmd/quiz/main.go -migrate

# Create .env file if it doesn't exist
init:
	@if not exist .env (echo Creating .env file && echo # API Configuration > .env && echo API_KEY=dev-api-key >> .env && echo APP_API_KEY=dev-app-api-key >> .env && echo PORT=8080 >> .env && echo. >> .env && echo # PostgreSQL Configuration >> .env && echo POSTGRES_USER=postgres >> .env && echo POSTGRES_PASSWORD=postgres >> .env && echo POSTGRES_DB=quiz_db >> .env && echo DB_HOST=postgres >> .env && echo DB_PORT=5432 >> .env && echo DB_USER=postgres >> .env && echo DB_PASSWORD=postgres >> .env && echo DB_NAME=quiz_db >> .env && echo .env file created) else (echo .env file already exists)

.PHONY: build run docker-build docker-run docker-compose swagger clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
BINARY_NAME=quiz-api
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/quiz

# Run the application
run:
	$(GORUN) ./cmd/quiz/main.go

# Build and run with docker
docker-build:
	docker build -t nosmo/quiz-api -f ./cmd/quiz/Dockerfile .

docker-run:
	docker run -p 8082:8080 nosmo/quiz-api

# Docker Compose commands
docker-compose-up:
	docker-compose -f ./cmd/quiz/docker-compose.yml up -d

docker-compose-down:
	docker-compose -f ./cmd/quiz/docker-compose.yml down

docker-compose-logs:
	docker-compose -f ./cmd/quiz/docker-compose.yml logs -f

# Generate swagger documentation
swagger:
	swag init -g cmd/quiz/main.go -o cmd/quiz/docs

# Clean generated binaries
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Help command
help:
	@echo "Make commands:"
	@echo "build - Build the application"
	@echo "run - Run the application"
	@echo "docker-build - Build the Docker image"
	@echo "docker-run - Run the application in Docker"
	@echo "docker-compose-up - Start all services with Docker Compose"
	@echo "docker-compose-down - Stop all services with Docker Compose"
	@echo "docker-compose-logs - View logs from all services"
	@echo "swagger - Generate swagger documentation"
	@echo "clean - Remove build artifacts" 