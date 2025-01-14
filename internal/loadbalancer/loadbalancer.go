package loadbalancer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type LoadBalancer interface {
	NextBackend() (*url.URL, error)
	Config() Config
}

type RoundRobin struct {
	backends        []*url.URL
	current         int
	mu              sync.Mutex
	config          Config
	healthCheckFreq time.Duration
	healthStatus    map[string]bool
	logger          *zap.Logger
}

type Config struct {
	Timeout    time.Duration
	Backoff    time.Duration
	MaxRetries int
}

func NewRoundRobin(backendURLs []string, logger *zap.Logger) (*RoundRobin, error) {
	backends := make([]*url.URL, 0, len(backendURLs))

	// Parse backend URLs
	for _, urlStr := range backendURLs {
		backendURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, fmt.Errorf("invalid backend URL %s: %w", urlStr, err)
		}
		backends = append(backends, backendURL)
	}

	if len(backends) == 0 {
		return nil, errors.New("no backends provided")
	}

	rr := &RoundRobin{
		backends: backends,
		current:  0,
		config: Config{
			Timeout:    10 * time.Second,
			Backoff:    1 * time.Second,
			MaxRetries: 3,
		},
		healthStatus:    make(map[string]bool),
		healthCheckFreq: 30 * time.Second,
		logger:          logger,
	}

	// Initial health check
	for _, backend := range rr.backends {
		rr.healthStatus[backend.Host] = checkBackendHealth(backend.String(), logger)
	}

	// Ensure at least one backend is healthy
	logger.Info("Ensuring at least one backend is healthy")
	healthy := false
	for _, status := range rr.healthStatus {
		if status {
			healthy = true
			break
		}
	}
	if !healthy {
		return nil, errors.New("no healthy backends available on startup")
	}

	// Start periodic health checks
	go rr.healthChecker()

	return rr, nil
}

// Check for next available backend end
func (rr *RoundRobin) NextBackend() (*url.URL, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	numBackends := len(rr.backends)
	if numBackends == 0 {
		return nil, errors.New("no backends available")
	}

	// loop through all backends and find first available one in RR, exist if none is available
	for i := 0; i < numBackends; i++ {
		backend := rr.backends[rr.current]
		rr.current = (rr.current + 1) % numBackends

		if rr.healthStatus[backend.Host] {
			return backend, nil
		}
	}

	return nil, errors.New("no healthy backends available")
}

func (rr *RoundRobin) healthChecker() {
	ticker := time.NewTicker(rr.healthCheckFreq)
	defer ticker.Stop()
	for range ticker.C {
		for _, backend := range rr.backends {
			healthy := checkBackendHealth(backend.Host, rr.logger)
			rr.mu.Lock()
			rr.healthStatus[backend.Host] = healthy
			rr.mu.Unlock()
			if !healthy {
				// observability
				rr.logger.Warn("Backend marked as unhealthy", zap.String("backend", backend.Host))
			}
		}
	}
}

func checkBackendHealth(backend string, logger *zap.Logger) bool {

	if !strings.HasPrefix(backend, "http://") && !strings.HasPrefix(backend, "https://") {
		backend = "http://" + backend
	}

	resp, err := http.Get(backend + "/health")
	if err != nil {
		logger.Error("Health check request failed", zap.String("backend", backend), zap.Error(err))
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
