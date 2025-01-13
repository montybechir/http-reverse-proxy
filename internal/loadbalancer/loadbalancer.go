package loadbalancer

import (
	"errors"
	"net/url"
	"sync"
	"time"
)

type LoadBalancer interface {
	NextBackend() (*url.URL, error)
	Config() Config
}

type RoundRobin struct {
	backends []*url.URL
	current  int
	mu       sync.Mutex
	config   Config
}

type Config struct {
	Timeout    time.Duration
	Backoff    time.Duration
	MaxRetries int
}

func NewRoundRobin(backendURLs []string) (*RoundRobin, error) {
	backends := make([]*url.URL, len(backendURLs))
	for i, b := range backendURLs {
		parsed, err := url.Parse(b)
		if err != nil {
			return nil, err
		}
		backends[i] = parsed
	}
	return &RoundRobin{
		backends: backends,
		current:  0,
		config: Config{
			Timeout:    10 * time.Second,
			Backoff:    1 * time.Second,
			MaxRetries: 3,
		},
	}, nil
}

func (rr *RoundRobin) NextBackend() (*url.URL, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if len(rr.backends) == 0 {
		return nil, errors.New("no backends available")
	}
	backend := rr.backends[rr.current]
	rr.current = (rr.current + 1) % len(rr.backends)
	return backend, nil
}

func (rr *RoundRobin) Config() Config {
	return rr.config
}
