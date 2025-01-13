package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func RateLimittingMiddleware(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}

}
