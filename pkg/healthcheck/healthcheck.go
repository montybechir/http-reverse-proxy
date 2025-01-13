package healthcheck

import (
	"encoding/json"
	"http-reverse-proxy/pkg/models"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HealthChecker interface {
	Status() error
}

type ProxyHealthChecker struct {
	backends []string
	client   *http.Client
}

func HealthHandler(
	logger *zap.Logger,
	version string,
	startTime time.Time,
	checkers []HealthChecker,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the health check access
		logger.Info("Health check endpoint accessed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
		)

		// Initialize HealthStatus
		health := models.HealthStatus{
			Status:    "healthy",
			Uptime:    time.Since(startTime).Round(time.Second).String(),
			Version:   version,
			Checks:    make(map[string]string),
			Timestamp: time.Now(),
		}

		// Perform checks if any HealthCheckers are provided
		if len(checkers) > 0 {
			for _, checker := range checkers {
				// Assuming each HealthChecker has a unique identifier (type or name)
				// For demonstration, we'll use the type name as the key
				key := "dependency"

				// Perform the health check
				err := checker.Status()
				if err != nil {
					health.Status = "unhealthy"
					health.Checks[key] = "fail: " + err.Error()
				} else {
					health.Checks[key] = "pass"
				}
			}
		}

		// Determine the HTTP status code based on overall health
		var httpStatus int
		if health.Status == "healthy" {
			httpStatus = http.StatusOK
		} else {
			httpStatus = http.StatusServiceUnavailable
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)

		// Encode and send the response
		if err := json.NewEncoder(w).Encode(health); err != nil {
			logger.Error("Failed to encode health status response", zap.Error(err))
			// In case of encoding failure, send a plain text response
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func NewProxyHealthChecker(backends []string) *ProxyHealthChecker {
	return &ProxyHealthChecker{
		backends: backends,
		client: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       100,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		},
	}
}
