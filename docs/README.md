# Swagger Documentation

This directory contains the auto-generated Swagger documentation for the API.

## Generating Documentation

To generate the Swagger documentation, you need to install the `swag` command-line tool:

```bash
# Install swag CLI
go install github.com/swaggo/swag/cmd/swag@latest
```

Then, from the root directory of the project, run:

```bash
# Generate swagger docs
swag init --parseDependency --parseDepth 2 -g cmd/quiz/main.go
```

## Accessing the Documentation

Once the API is running, you can access the Swagger documentation UI at:

```
http://localhost:8080/swagger/
```

## Security

The API uses two security mechanisms:

1. **JWT Authentication** (`Authorization` header)
   - Format: `Bearer <token>`
   - Used for authenticated user endpoints

2. **API Key** (`X-API-Key` header)
   - Used for administrative endpoints
   - The key is defined in the environment variable `API_KEY` 