
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git zeromq-dev build-base

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates zeromq

# Copy binary from builder
COPY --from=builder /app/main .

# Copy configuration and web files
COPY --from=builder /app/config ./config
COPY --from=builder /app/web/build ./web/build

# Create logs directory
RUN mkdir -p logs

# Expose ports
EXPOSE 8080 5555

# Run the application
CMD ["./main"]