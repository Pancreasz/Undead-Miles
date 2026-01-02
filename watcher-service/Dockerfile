# Stage 1: Build
# OLD: FROM golang:1.23-alpine AS builder
# NEW:
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the worker binary
RUN go build -o main ./cmd/worker

# Stage 2: Run
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .

# Expose port 8081 (Watcher API)
EXPOSE 8081

CMD ["./main"]