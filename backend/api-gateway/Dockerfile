# syntax=docker/dockerfile:1

# ----------- Build Stage -----------
ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder

# Install git for go get if needed (for private modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY --link go.mod go.sum ./

# Download dependencies (use cache mounts for speed)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy the rest of the source code
COPY --link . .

# Build the Go binary (static build, strip debug info)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o api-gateway ./main.go

# ----------- Final Stage -----------
FROM alpine:latest AS final

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/api-gateway /app/api-gateway

# Copy static assets (e.g., docs/swagger.json, docs/swagger.yaml) if needed for http-swagger
COPY --from=builder /app/docs /app/docs

# Expose the service port
EXPOSE 5000

# Use non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["/app/api-gateway"]

# Note: Pass environment variables (e.g., REDIS_URL, JWT_SECRET_KEY, etc.) at runtime, not via .env files.
