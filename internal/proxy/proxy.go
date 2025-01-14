package proxy

import (
	"encoding/json"
	"http-reverse-proxy/pkg/models"
	"net/http"

	"http-reverse-proxy/internal/loadbalancer"

	"go.uber.org/zap"
)

type ReverseProxy struct {
	LoadBalancer *loadbalancer.RoundRobin
	Logger       *zap.Logger
	Config       *models.Config
}

// NewReverseProxy initializes a new ReverseProxy instance
func NewReverseProxy(lb *loadbalancer.RoundRobin, logger *zap.Logger, config *models.Config) (*ReverseProxy, error) {
	return &ReverseProxy{
		LoadBalancer: lb,
		Logger:       logger,
		Config:       config,
	}, nil
}

// StatusResponse defines the structure of the status response
type StatusResponse struct {
	Status  string `json:"status"`
	Uptime  string `json:"uptime"`
	Version string `json:"version"`
}

// StatusHandler provides the current status and uptime of the proxy
func (rp *ReverseProxy) StatusHandler(w http.ResponseWriter, r *http.Request) {
	// Example static response; in a real-world scenario, dynamically calculate uptime and version
	response := StatusResponse{
		Status:  "running",
		Uptime:  "72h",   // This should be dynamically calculated
		Version: "1.0.0", // Ideally fetched from build variables
	}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		rp.Logger.Error("Failed to encode status response", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log the successful status check
	rp.Logger.Info("Status endpoint accessed",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)
}

// Helper function to copy HTTP headers
func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
