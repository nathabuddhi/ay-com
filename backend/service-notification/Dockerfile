# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder

# Install git for go mod download if needed
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

# Build the Go binary (static build)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service-notification ./main.go

# --- Final minimal image ---
FROM alpine:latest AS final

# Install CA certificates for HTTPS
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache ca-certificates

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the statically built binary from builder
COPY --from=builder /app/service-notification ./service-notification

# Expose the gRPC port
EXPOSE 5009

USER appuser

ENTRYPOINT ["./service-notification"]
