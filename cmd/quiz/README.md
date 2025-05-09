# Quiz Service

This is a microservice for the quiz application, providing API endpoints to manage quizzes, questions, answers, and user submissions.

## Features

- RESTful API for managing quizzes and related resources
- Authentication and authorization
- PostgreSQL database for data storage
- Swagger documentation
- Docker support for easy deployment

## Requirements

- Go 1.21 or higher
- PostgreSQL 15 or higher
- Docker and Docker Compose (for containerized deployment)

## Configuration

Configuration is managed through a YAML file at `cmd/quiz/config/config.yml` or through environment variables.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_NAME | Application name | quiz-service |
| APP_ENVIRONMENT | Application environment | development |
| APP_API_KEY | API key for protected endpoints | default-api-key |
| SERVER_PORT | Port for the server to listen on | 8080 |
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL user | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | PostgreSQL database name | quiz_db |

## Running the Service

### Using Go directly

```bash
# Navigate to the project root
cd /path/to/project

# Build and run
go run cmd/quiz/main.go
```

### Using the Makefile

```bash
# Run with Make
make -f cmd/quiz/Makefile run
```

### Using Docker

```bash
# Build the Docker image
make -f cmd/quiz/Makefile docker-build

# Run the container
make -f cmd/quiz/Makefile docker-run
```

### Using Docker Compose

```bash
# Start all services
make -f cmd/quiz/Makefile docker-compose-up

# View logs
make -f cmd/quiz/Makefile docker-compose-logs

# Stop all services
make -f cmd/quiz/Makefile docker-compose-down
```

## API Documentation

Swagger documentation is available at `/docs/` when the service is running. You can also generate the Swagger docs with:

```bash
make -f cmd/quiz/Makefile swagger
```

## Project Structure

```
cmd/quiz/                    # Quiz service main directory
├── config/                  # Configuration files
│   └── config.yml           # Default configuration
├── docs/                    # Swagger documentation
├── Dockerfile               # Dockerfile for containerization
├── docker-compose.yml       # Docker Compose for local development
├── Makefile                 # Make commands for development
├── main.go                  # Application entry point
├── database.go              # Quiz database setup
└── article_database.go      # Article database setup

pkg/quiz/                    # Quiz service packages
├── config/                  # Configuration package
├── domain/                  # Domain models and interfaces
│   ├── entity/              # Data entities
│   └── repository/          # Repository interfaces
├── handlers/                # HTTP handlers
├── middleware/              # HTTP middleware
├── repository/              # Repository implementations
│   └── gorm/                # GORM implementation
└── services/                # Business logic services
```

## License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details. 