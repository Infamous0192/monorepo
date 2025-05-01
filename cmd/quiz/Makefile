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