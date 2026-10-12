# Build Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build essentials if needed
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# modernc.org/sqlite is pure Go, so CGO_ENABLED=0 is safe and preferred for alpine.
# NOTE: no -a flag on purpose. It force-rebuilds the entire stdlib on every build,
# tripling build time and tripling the Go build cache/disk writes (a common cause of
# ENOSPC failures). -installsuffix cgo is a no-op when CGO_ENABLED=0.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o main ./cmd/bot

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create directory for SQLite
RUN mkdir -p /app/data

# Copy binary from builder
COPY --from=builder /app/main .

# Command to run
CMD ["./main"]
