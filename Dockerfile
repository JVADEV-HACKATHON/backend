# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Build the seeder application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

WORKDIR /root/

# Copy the binaries from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/seed .

# Copy .env file if exists
COPY --from=builder /app/.env* ./

# Copy entrypoint script
COPY scripts/entrypoint.sh ./entrypoint.sh

# Fix line endings and set permissions
RUN dos2unix ./entrypoint.sh || true
RUN chmod +x ./entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]