.PHONY: all build test clean docker-up docker-down docker-up-prod docker-down-prod lint generate hot-reload-chat hot-reload-telegram build-user-service build-product-service build-dev-user-service build-dev-product-service swagger

all: build

build:
	go build -v ./...

test:
	go test -v ./...

clean:
	go clean

# Development environment
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Production environment
docker-up-prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

docker-down-prod:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# Hot reload specific services
hot-reload-chat:
	docker-compose up -d --build chat-service

hot-reload-telegram:
	docker-compose up -d --build telegram-service

# Build specific service
build-user-service:
	docker build -f services/user-service/Dockerfile --target production -t your-org/user-service:latest .

build-product-service:
	docker build -f services/product-service/Dockerfile --target production -t your-org/product-service:latest .

# Development builds
build-dev-user-service:
	docker build -f services/user-service/Dockerfile --target development -t your-org/user-service:dev .

build-dev-product-service:
	docker build -f services/product-service/Dockerfile --target development -t your-org/product-service:dev .

lint:
	golangci-lint run

generate:
	go generate ./...

# Generate swagger documentation
swagger:
	swag init --parseDependency --parseDepth 2 -g cmd/quiz/main.go

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o quiz_linux -v ./cmd/quiz
