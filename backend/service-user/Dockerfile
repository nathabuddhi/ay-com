# syntax=docker/dockerfile:1

# ----------- Build Stage -----------
ARG GO_VERSION=1.24.2
FROM golang:${GO_VERSION}-alpine AS builder

# Install git for go get if needed (for private modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY --link go.mod go.sum ./

# Download dependencies (with cache mounts)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

# Copy the rest of the source code
COPY --link . .

# Build the Go binary (static build, strip debug info)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o service-user ./main.go

# ----------- Final Stage -----------
FROM alpine:latest AS final

# Install CA certificates for HTTPS
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache ca-certificates

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/service-user /app/service-user

# Expose the gRPC port
EXPOSE 5001

# Use non-root user
USER appuser

# Entrypoint
ENTRYPOINT ["/app/service-user"]
