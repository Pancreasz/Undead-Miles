# 1. Build Stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o notification-service main.go

# 2. Run Stage (Small image)
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/notification-service .

# Expose the port
EXPOSE 8082

# Run the binary
CMD ["./notification-service"]