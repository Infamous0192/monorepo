FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

# Copy source code and build
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o quiz ./cmd/quiz

# Use a small alpine image for the final container
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/quiz .
COPY --from=builder /app/cmd/quiz/config /app/config

# Create a non-root user to run the application
RUN adduser -D -g '' -u 1000 appuser

# Create uploads directory and set permissions
RUN mkdir -p /app/uploads && \
    chown -R appuser:appuser /app && \
    chmod -R 755 /app && \
    chmod 777 /app/uploads

# Switch to non-root user
USER appuser

# Set environment variables
ENV PORT=8080
ENV APP_ENVIRONMENT=production

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./quiz", "--config", "/app/config/config.yml"] 