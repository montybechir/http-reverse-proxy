# HTTP Reverse Proxy

A Go-based HTTP reverse proxy featuring load balancing, rate limiting, and CORS support. This reverse proxy is designed without relying on third-party implementations or the net/http/httputil package, ensuring a lightweight and customizable solution.

### Table of Contents

    Features
    Getting Started
        Prerequisites
        Installation
        Configuration
        Running the Proxy
    Usage
    Project Structure
    Testing
    Design Decisions & Limitations
    Scaling
    Security Enhancements
    Resources
    Contact

## Features

Load Balancing (Round Robin):

    - Distributes incoming requests evenly across healthy backend servers using a simple Round Robin algorithm.
    - Easily extendable to other balancing strategies, such as Least Connections or Weighted Round Robin.

Rate Limiting:

    - Prevents abuse by limiting the number of requests a client can make within a specified timeframe.
    - Configurable to adjust thresholds based on application needs.

CORS Support:

    - Configurable to allow connections from any origin or restrict to specific hosts.
    - Customizable headers and methods to control cross-origin requests.

Health Checks with Recovery:

    - Continuously monitors the health of backend servers.
    - Automatically reintegrates recovered backends into the pool without manual intervention.
    - Ensures the proxy only forwards requests to active and healthy backends.

Multiple Backend Support:

    - Easily configurable to support multiple backend servers via YAML configuration files.
    - Ensures high availability by requiring at least one healthy backend to operate.

Docker Integration:

    - Dockerfiles provided for the proxy and backend servers for streamlined containerization.
    - docker-compose setup included for easy orchestration and deployment.

## Prerequisites

- Docker & Docker Compose
- Go 1.22+

## Quick Start

```bash
git clone https://github.com/montybechir/http-reverse-proxy.git

# or

git clone git@github.com:montybechir/http-reverse-proxy.git

```

Configure Docker Resource Sharing:

    Add the project directory (e.g., /Users/montasir/dev/http-reverse-proxy) to Docker's file sharing settings.
    For Mac:
        Navigate to Docker -> Preferences... -> Resources -> File Sharing.
        Add the project directory and apply changes.
        Restart Docker if necessary.
        Docker File Sharing Documentation

Configuration

    Edit Configuration Files:
        Modify /configs/config.yaml to set proxy configurations.
        Add or update backend configurations in /configs/backenda.yaml and /configs/backendb.yaml.
        Ensure at least one backend is active to allow the proxy to start successfully.

### Build and run the container

```bash

docker-compose up --build

```

### Verify Successful Launch:

Monitor the logs to ensure all services (proxy and backends) start without errors.

Example successful log entries:

```bash
 ✔ Service backendb  Built                                                                                              2.0s
 ✔ Service backenda  Built                                                                                              2.1s
 ✔ Service proxy     Built                                                                                              0.5s
Attaching to backenda-1, backendb-1, proxy-1
backendb-1  | {"level":"info","timestamp":"2025-01-14T17:40:59.728Z","caller":"server/server.go:88","msg":"Server is starting","address":":60409"}
backenda-1  | {"level":"info","timestamp":"2025-01-14T17:40:59.728Z","caller":"server/server.go:88","msg":"Server is starting","address":":60408"}
backenda-1  | {"level":"info","timestamp":"2025-01-14T17:40:59.839Z","caller":"middleware/logging.go:19","msg":"Incoming request","method":"GET","path":"/health","remote_addr":"172.18.0.4:58352"}
backenda-1  | {"level":"info","timestamp":"2025-01-14T17:40:59.943Z","caller":"middleware/logging.go:29","msg":"Completed request","status":200,"method":"GET","path":"/health","duration":0.103956}
backendb-1  | {"level":"info","timestamp":"2025-01-14T17:40:59.955Z","caller":"middleware/logging.go:19","msg":"Incoming request","method":"GET","path":"/health","remote_addr":"172.18.0.4:42270"}
backendb-1  | {"level":"info","timestamp":"2025-01-14T17:41:00.060Z","caller":"middleware/logging.go:29","msg":"Completed request","status":200,"method":"GET","path":"/health","duration":0.1035535}
proxy-1     | {"level":"info","timestamp":"2025-01-14T17:41:00.066Z","caller":"loadbalancer/loadbalancer.go:58","msg":"Ensuring at least one backend is healthy"}
proxy-1     | {"level":"info","timestamp":"2025-01-14T17:41:00.070Z","caller":"proxy/main.go:71","msg":"Starting server","address":":8080"}
proxy-1     | {"level":"info","timestamp":"2025-01-14T17:41:09.843Z","caller":"middleware/logging.go:19","msg":"Incoming request","method":"GET","path":"/health","remote_addr":"[::1]:53694"}
proxy-1     | {"level":"info","timestamp":"2025-01-14T17:41:09.948Z","caller":"proxy/handler.go:89","msg":"Request proxied successfully","method":"GET","path":"/health","backend":"http://backenda:60408/health","bytes_written":23,"status":200}

```

## Usage

You can interact with the proxy using curl or API clients like Postman.

```bash

curl localhost:8080/pleasefowardthiscommand -v

```

### expected result

```bash
MB-Machine:~ montasir$ curl localhost:8080/pleasefowardthiscommand -v
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /pleasefowardthiscommand HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
< Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
< Access-Control-Allow-Origin:
< Content-Length: 23
< Content-Type: text/plain; charset=utf-8
< Date: Tue, 14 Jan 2025 17:43:59 GMT
<
* Connection #0 to host localhost left intact
Response from Backend B

```

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

Run integration tests to ensure all components function as expected.

Note: running all tests in VS code works.

```bash
# Run all integration tests (no parallel)
go test -v ./tests/integration/... -p 1


# Run specific test

go test ./tests/integration -run TestHealthChecks

```

### Design & Limitations

Design Decisions

    Language Choice:
        - Go: Selected for its simplicity, high performance, and excellent concurrency support, making it ideal for building scalable network services like reverse proxies.
    Custom Implementation:
        - Avoided third-party packages and the net/http/httputil package to adhere to the assignment constraints and ensure a deeper understanding of proxy mechanics.
    Load Balancing Strategy:
        - Implemented Round Robin for its simplicity and even distribution of requests across backends.
        - Structured to allow easy integration of alternative algorithms in the future.

### Limitations

    Security Features:
    - Does not currently support SSL termination or advanced authentication mechanisms.
    Caching & Compression:
    - Lacks built-in caching and response compression, which could enhance performance.
    Load Balancing Sophistication:
    - Currently limited to Round Robin without considering backend load or response times.
    Error Handling:
    - Basic error handling in place; more granular logging and alerting could be beneficial.

### Scaling

To handle increased traffic and ensure high availability, the system can be scaled as follows:

    Horizontal Scaling:
        - Deploy the load balancer using A records on DNS records using a round-robin setup. However, this would require us to maintain the IP addresses of incoming requests across our proxies to keep track of requests. This would also require us to ensure traffic is only routed to healthy proxies and may introduce cache invalidation headaches.
    Enhanced Load Balancing:
        - Implement more sophisticated load balancing algorithms (e.g., Least Connections, Weighted Round Robin) to optimize request distribution based on backend performance.
    Caching Mechanisms:
<<<<<<< Updated upstream
        Integrate a caching layer (e.g., Redis or in-memory caches) to store frequently accessed data, reducing load on backend servers.
    Microservices Architecture:
        Break down components into microservices to allow independent scaling based on demand.
=======
        - Integrate a caching layer (e.g., Redis or in-memory caches) to store frequently accessed data, reducing load on backend servers.
    Redundancy & Failover:
        - We can achieve high availability via redundancy for failover by having an additional proxy as well. This can be deployed in another region to further reduce the odds of complete loss of access. However, this may introduce data consistency and latency issues.

K8 can be used to manage and auto-scale proxies depending on load.
>>>>>>> Stashed changes

### Security

Enhancing the security posture of the reverse proxy involves several strategies:

    Input Validation:
        Implement rigorous validation of incoming requests to prevent common attacks such as SQL injection, Cross-Site Scripting (XSS), and others listed in the OWASP Top 10.
    TLS/SSL Support:
        Enable TLS to encrypt data in transit, ensuring secure communication between clients and the proxy.
        Support SSL termination at the proxy, or opt for end-to-end encryption based on security requirements.
    Authentication & Authorization:
        Introduce authentication middleware to enforce API key checks or integrate OAuth/JWT-based authentication.
    Environment Variable Management:
        Use environment variables or secure key management services (e.g., Azure Key Vault, AWS Secrets Manager) to manage sensitive information like API keys and certificates securely.
    Rate Limiting & Throttling:
        Enhance rate limiting strategies to include dynamic thresholds based on user roles or IP reputation, etc. If more than one instance of a proxy is used in our design, we'd need to maintain the request counts across different instances for rate-limitting as well.
    Logging & Monitoring:
        Implement comprehensive logging and monitoring to detect and respond to suspicious activities promptly.
        Integrate with monitoring tools like Prometheus or ELK Stack for real-time insights.
    Security:
        To overcome the shortcomings of the current implementation, leveraging Azure Front Door for CDN + Web Application Firewalls would take care of the caching, and security limitations above.

### Resources

- Go Documentation: https://golang.org/doc/
- Docker Documentation: https://docs.docker.com/
- Proxies vs Reverse Proxies: https://blog.bytebytego.com/p/ep25-proxy-vs-reverse-proxy
- System Design Primer: https://github.com/donnemartin/system-design-primer?tab=readme-ov-file#reverse-proxy-web-server
- AI tools: Github Copilot + Azure AI Services
- Microservices + Go: https://www.youtube.com/watch?v=GtSg1H7SU5Y&list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_&index=12&ab_channel=NicJackson
- HTTP Server Best Practices: https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
