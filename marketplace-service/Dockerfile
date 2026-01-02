# Stage 1: Build the Go Binary
# OLD: FROM golang:1.23-alpine AS builder
# NEW:
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
# -o main : Output filename
# ./cmd/api : Entry point
RUN go build -o main ./cmd/api

# Stage 2: Run the Binary (Small Image)
FROM alpine:latest
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run
CMD ["./main"]