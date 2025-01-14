FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy entire project for module resolution
COPY . .

# Download dependencies
RUN go mod download

# Build with proper flags
RUN CGO_ENABLED=0 GOOS=linux go build -o proxy ./cmd/proxy

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy binary and config
COPY --from=builder /app/proxy .
COPY --from=builder /app/configs/config.yaml /app/configs/

EXPOSE 60408

# Add health check
HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./proxy"]