# Build Stage
FROM golang:1.25.1-alpine AS builder

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
# modernc.org/sqlite is pure Go, so CGO_ENABLED=0 is safe and preferred for alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/bot

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
