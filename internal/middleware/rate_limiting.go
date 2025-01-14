package middleware

import (
	"http-reverse-proxy/pkg/models"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       *sync.RWMutex
	rate     rate.Limit
	burst    int
	logger   *zap.Logger
	cleanup  time.Duration
}

func NewRateLimiter(cfg *models.RateLimitConfig, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		mu:       &sync.RWMutex{},
		rate:     rate.Limit(cfg.RequestsPerMinute) / 60, // Convert to per-second
		burst:    cfg.Burst,
		logger:   logger,
		cleanup:  time.Hour,
	}

	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			limiter := rl.getVisitor(ip)

			if !limiter.Allow() {
				rl.logger.Warn("Rate limit exceeded",
					zap.String("ip", ip),
					zap.Float64("rate", float64(rl.rate)),
					zap.Int("burst", rl.burst),
				)
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
