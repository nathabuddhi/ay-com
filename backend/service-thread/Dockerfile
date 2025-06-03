# syntax=docker/dockerfile:1

# ----------- Build Stage -----------
ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory inside the builder
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
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service-thread ./main.go

# ----------- Final Stage -----------
FROM alpine:latest AS final

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install ca-certificates for HTTPS
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/service-thread ./service-thread

# Copy any required static assets (if any, e.g. migrations, etc.)
# (none in this service)

# Set permissions
RUN chown -R appuser:appgroup /app

USER appuser

# Expose the gRPC port
EXPOSE 5002

# Set environment variables (if any defaults, but secrets should be passed at runtime)
# ENV EXAMPLE_ENV=example

# Start the service
ENTRYPOINT ["./service-thread"]
