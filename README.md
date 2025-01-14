# HTTP Reverse Proxy

A Go-based HTTP reverse proxy with load balancing, rate limiting, and CORS support.

## Quick Start

```bash
# Clone and run
git clone https://github.com/yourusername/http-reverse-proxy.git
cd http-reverse-proxy
docker-compose up --build

# Test the proxy
curl localhost:8080/test
```

## Features

- Load Balancing (Round Robin)
- Rate Limiting
- CORS Support
- Health Checks
- Multiple Backend Support
- Docker Integration

# Development

## Prerequisites

- Docker & Docker Compose
- Go 1.20+

## Project Structure

```bash
reverse-proxy/
├── cmd/
│   └── proxy/
│       └── main.go
│       └── Dockerfile
│   └── backenda/
│       └── main.go
│       └── Dockerfile
│   └── backendb/
│       └── main.go
│       └── Dockerfile
├── internal/
│   ├── proxy/
│   │   ├── handler.go
│   │   └── proxy.go
│   │   └── router.go
│   ├── loadbalancer/
│   │   └── loadbalancer.go
│   └── middleware/
│       ├── logging.go
│       ├── middleware.go
│       ├── rate_limiting.go
│       └── cors.go
├── pkg/
│   ├── utils/
│   │   ├── headers.go
│   │   └── config.go
│   ├── models/
│   │   ├── config.go
│   ├── server/
│   │   ├── server.go
│   └── healthcheck/
│       └── healthcheck.go
├── configs/
│   ├── config.yaml
│   └── logging.yaml
│   └── backenda.yaml
│   └── backendb.yaml
├── scripts/
│   └── wait-for.sh
│   └── run-tests.sh
├── tests/
│   ├── helpers/
│   │   ├── logger.go
│   │   └── mock.go
│   │   ├── proxy.go
│   │   └── utils.go
│   ├── utils/
│   │   ├── cors_test.go
│   │   └── health_test.go
│   │   ├── loadbalancer_test.go
│   │   └── proxy_test.go
│   │   ├── rate_limiting_test.go
│   │   └── status_test.go
│   ├── Dockerfile
├── docker-compose.yaml
├── go.mod
├── go.sum
└── README.md
```

## Testing

```bash
# Run all integration tests (no parallel)
go test -v ./tests/integration/... -p 1


# Run specific test
go test ./tests/integration -run TestHealthChecks

```
