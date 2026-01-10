# Base stage - shared by all other stages ==========================================================
FROM golang:1.25-alpine AS base

WORKDIR /app

# Install git and build dependencies for rtmidi (MIDI support)
RUN apk add --no-cache \
    git \
    build-base \
    alsa-lib-dev \
    jack-dev

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Development stage - includes Air and Delve for hot reload and debugging
FROM base AS development

# Install Air for hot reloading (latest version, requires Go 1.25+)
RUN go install github.com/air-verse/air@latest

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Copy source code
COPY . .

# Expose Delve debugger port
EXPOSE 2345

# Run Air (which will manage rebuilding and launching via Delve)
CMD ["air", "-c", ".air.docker.toml"]

# Builder stage - compiles optimized production binary ==========================================================
FROM base AS builder

# Copy source code
COPY . .

# Build optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o reactor \
    .

# Production stage - minimal image with binary only ==========================================================
FROM alpine:latest AS production

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/reactor .

# Run the application
CMD ["./reactor"]
