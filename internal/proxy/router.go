package proxy

import (
	"net/http"
	"time"
)

var (
	startTime = time.Now()
	version   = "1.0.0" // Should come from build info
)

// SetupRoutes registers all necessary routes and returns the mux
func (rp *ReverseProxy) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// checkers := []healthcheck.HealthChecker{
	// 	healthcheck.NewProxyHealthChecker(rp.Config.Backends),
	// }
	// //Register health and status handlers
	// mux.HandleFunc("/health", healthcheck.HealthHandler(rp.Logger))
	mux.HandleFunc("/status", rp.StatusHandler)

	// Proxy all other routes
	mux.Handle("/", http.HandlerFunc(rp.ProxyHandler))

	return mux
}
