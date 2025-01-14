// tests/helpers/mock_backend.go (advanced)

package helpers

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"go.uber.org/zap"
)

// MockBackend represents a mock backend server with dynamic response capabilities.
type MockBackend struct {
	Server            *httptest.Server
	RequestCh         chan *http.Request
	Status            int
	Response          string
	Headers           map[string]string
	Delay             time.Duration                                          // Introduce artificial delays.
	DynamicResponseFn func(r *http.Request) (int, string, map[string]string) // Customize responses.
	mu                sync.Mutex
	Logger            *zap.Logger
}

// Close shuts down the mock backend server.
func (mb *MockBackend) Close() {
	mb.Server.Close()
}

// SetStaticResponse configures static response parameters.
func (mb *MockBackend) SetStaticResponse(status int, response string, headers map[string]string) {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	mb.Status = status
	mb.Response = response
	mb.Headers = headers
}

// NewMockBackend creates and starts a new mock backend server with optional headers.
func NewMockBackend(status int, response string, headers map[string]string, logger *zap.Logger) *MockBackend {
	backend := &MockBackend{
		RequestCh: make(chan *http.Request, 100),
		Status:    status,
		Response:  response,
		Headers:   headers,
		Logger:    logger,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backend.mu.Lock()
		backend.RequestCh <- r
		backend.mu.Unlock()

		// Set headers if any.
		for key, value := range backend.Headers {
			w.Header().Set(key, value)
		}

		w.WriteHeader(backend.Status)
		w.Write([]byte(backend.Response))
	})

	backend.Server = httptest.NewServer(handler)
	return backend
}

// GetRequests retrieves all received requests.
func (mb *MockBackend) GetRequests() []*http.Request {
	mb.mu.Lock()
	defer mb.mu.Unlock()
	requests := []*http.Request{}
	for len(mb.RequestCh) > 0 {
		req := <-mb.RequestCh
		requests = append(requests, req)
	}
	return requests
}
