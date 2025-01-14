package proxy

import (
	"net/http"
)

// SetupRoutes registers all necessary routes and returns the mux
func (rp *ReverseProxy) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", rp.StatusHandler)

	// Proxy all other routes
	mux.Handle("/", http.HandlerFunc(rp.ProxyHandler))

	return mux
}
