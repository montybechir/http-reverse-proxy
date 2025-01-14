package middleware

import (
	"http-reverse-proxy/pkg/models"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func CORSMiddleware(cfg *models.CORSConfig, log *zap.Logger) Middleware {

	// Convert slice to map for O(1) lookup
	originsMap := make(map[string]bool)
	for _, origin := range cfg.AllowedOrigins {
		originsMap[origin] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Validate origin
			if !isOriginAllowed(origin, originsMap, cfg) {
				if cfg.Debug {
					log.Debug("CORS rejected",
						zap.String("origin", origin),
						zap.Strings("allowed", cfg.AllowedOrigins),
					)
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

func isOriginAllowed(origin string, originsMap map[string]bool, cfg *models.CORSConfig) bool {
	// Check exact match or if all origins are allowed
	if originsMap[origin] || originsMap["*"] {
		return true
	}

	// Check wildcard domains
	originURL, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := originURL.Host
	for allowed := range originsMap {
		if strings.HasPrefix(allowed, "*.") &&
			strings.HasSuffix(host, allowed[1:]) {
			return true
		}
	}

	return false
}
