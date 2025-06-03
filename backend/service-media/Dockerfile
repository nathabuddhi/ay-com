# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder

# Install git for go mod (if needed)
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

# Build the Go binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service-media

# --- Final minimal image ---
FROM alpine:latest AS final

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/service-media ./service-media

# Use non-root user
USER appuser

# Expose the gRPC port
EXPOSE 5013

# Entrypoint
ENTRYPOINT ["./service-media"]
