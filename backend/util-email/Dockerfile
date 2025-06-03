# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Install git for go get if needed, and ca-certificates for HTTPS
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache git ca-certificates

# Copy go.mod and go.sum first for dependency caching
COPY --link go.mod go.sum ./

# Download dependencies (use cache mounts for Go modules and build cache)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy the rest of the source code
COPY --link . .

# Build the Go binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o util-email

# --- Final minimal image ---
FROM alpine:latest AS final

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install ca-certificates for TLS
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/util-email ./util-email

# Use non-root user
USER appuser

# Expose port if needed (not strictly necessary for this service, as it listens to RabbitMQ)
# EXPOSE 8080

ENTRYPOINT ["./util-email"]
