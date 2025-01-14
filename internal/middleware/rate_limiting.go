package middleware

import (
	"net/http"
	"sync"
	"time"

	"http-reverse-proxy/pkg/models"

	"go.uber.org/zap"
)

type RateLimiter struct {
	requestsPerMinute int
	burst             int
	clients           map[string]*tokenBucket
	mu                sync.Mutex
	logger            *zap.Logger
}

type tokenBucket struct {
	tokens          int
	lastTokenRefill time.Time
}

func NewRateLimiter(cfg *models.RateLimitConfig, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		requestsPerMinute: cfg.RequestsPerMinute,
		burst:             cfg.Burst,
		clients:           make(map[string]*tokenBucket),
		logger:            logger,
	}
	go rl.cleanupExpiredBuckets()
	return rl
}

func (rl *RateLimiter) cleanupExpiredBuckets() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for client, bucket := range rl.clients {
			if time.Since(bucket.lastTokenRefill) > time.Minute {
				delete(rl.clients, client)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr

			rl.mu.Lock()
			bucket, exists := rl.clients[clientIP]
			if !exists {
				bucket = &tokenBucket{
					tokens:          rl.burst,
					lastTokenRefill: time.Now(),
				}
				rl.clients[clientIP] = bucket
			}

			// Refill tokens based on elapsed time
			elapsed := time.Since(bucket.lastTokenRefill)
			refillTokens := int(elapsed.Minutes() * float64(rl.requestsPerMinute))
			if refillTokens > 0 {
				bucket.tokens += refillTokens
				if bucket.tokens > rl.burst {
					bucket.tokens = rl.burst
				}
				bucket.lastTokenRefill = time.Now()
			}

			if bucket.tokens > 0 {
				bucket.tokens--
				rl.mu.Unlock()
				next.ServeHTTP(w, r)
			} else {
				rl.mu.Unlock()
				rl.logger.Warn("Rate limit exceeded", zap.String("client", clientIP))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			}
		})
	}
}
