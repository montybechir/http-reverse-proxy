package middleware

import (
	"http-reverse-proxy/pkg/models"
	"log"
	"net/http"
)

func CORSMiddleware(cfg *models.CORSConfig) Middleware {

	// Convert slice to map for O(1) lookup
	originsMap := make(map[string]bool)
	for _, origin := range cfg.AllowedOrigins {
		originsMap[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Allow requests with no origin (like curl)
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if origin is allowed
			if !originsMap[origin] {
				if cfg.Debug {
					log.Printf("CORS: Rejected origin: %s", origin)
				}
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}

			// Set secure CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			// w.Header().Set("Vary", "Origin")

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
