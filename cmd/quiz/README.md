# Quiz Service

This service handles the quiz functionality, including quiz creation, question management, user submissions, and scoring.

## Prerequisites

- Docker and Docker Compose installed on your system
- Go 1.21 or higher (for local development without Docker)

## Environment Variables

Before running the service, create a `.env` file in this directory with the following variables:

```
# API Configuration
API_KEY=dev-api-key
APP_API_KEY=dev-app-api-key
PORT=8080

# PostgreSQL Configuration
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=quiz_db
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=quiz_db
```

## Running with Docker Compose

### Development

To run the service in development mode with hot-reloading:

```powershell
docker-compose up
```

or

```powershell
docker-compose up -d
```

to run in detached mode.

### Production

To run the service in production mode:

```powershell
docker-compose -f docker-compose.prod.yml up -d
```

## Manual Setup

### Running locally without Docker

1. Set up a PostgreSQL database
2. Set the required environment variables
3. Run the service:

```powershell
go run .
```

## Database Migrations

Database migrations are handled automatically when the service starts. The schema is defined in the `database.go` and `article_database.go` files. 