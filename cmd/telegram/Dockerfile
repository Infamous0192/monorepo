# Build stage
FROM golang:1.24-alpine AS builder

# Install required packages
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/telegram-service ./cmd/telegram

# Development stage
FROM golang:1.24-alpine AS development

# Install required packages and Air for hot reloading
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create the config directory if it doesn't exist
RUN mkdir -p /app/cmd/telegram/config

# Expose port
EXPOSE 8080
EXPOSE 8443

# Command to run Air for hot reloading
CMD ["air", "-c", ".air/telegram.toml"]

# Final stage
FROM alpine:latest AS alpine

# Install ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/telegram-service .
COPY --from=builder /app/cmd/telegram/config/config.yml ./config/config.yml

# Expose port
EXPOSE 8080
EXPOSE 8443

# Command to run
CMD ["./telegram-service", "--config", "./config/config.yml"]
