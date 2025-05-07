FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o quiz ./cmd/quiz

# Development stage
FROM golang:1.24-alpine AS development
WORKDIR /app

# Install air for hot-reloading in development
RUN go install github.com/air-verse/air@latest

# Copy the entire project for development
COPY . .

# Command to run air for hot reloading
CMD ["air", "-c", ".air/quiz.toml"]

# Production stage
FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/quiz .

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./quiz"] 