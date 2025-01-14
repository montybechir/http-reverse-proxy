// tests/helpers/proxy_setup.go

package helpers

import (
	"context"
	"http-reverse-proxy/internal/loadbalancer"
	"http-reverse-proxy/internal/middleware"
	"http-reverse-proxy/internal/proxy"
	"http-reverse-proxy/pkg/logger"
	"http-reverse-proxy/pkg/models"
	"http-reverse-proxy/pkg/utils"
	"net/http"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// SetupProxy initializes the reverse proxy with given backend URLs and configuration overrides.
// It returns the running HTTP server and a teardown function to gracefully shutdown the server.
func SetupProxy(t *testing.T, backendURLs []string, configOverrides map[string]interface{}) (*http.Server, func()) {
	// Load default config.
	// Assuming tests are run from the project root.
	absPath, err := filepath.Abs(filepath.Join("..", "..", "configs", "config.yaml"))
	assert.NoError(t, err, "Failed to load config")

	// Load the configuration.
	config, err := utils.LoadConfig(absPath)
	assert.NoError(t, err, "Failed to load config")

	// Override config if necessary.
	if addr, ok := configOverrides["address"].(string); ok {
		config.Server.Address = addr
	}
	if corsCfg, ok := configOverrides["cors"].(models.CORSConfig); ok {
		config.CORS = corsCfg
	}
	if rateLimitCfg, ok := configOverrides["ratelimit"].(models.RateLimitConfig); ok {
		config.RateLimit = rateLimitCfg
	}
	if healthCheckFreq, ok := configOverrides["healthCheckFreq"].(time.Duration); ok {
		config.HealthCheck.Frequency = healthCheckFreq
	}

	// The default ocnfig don't have the settings we want
	config.Backends = backendURLs

	// Initialize logger.
	zapLogger, err := logger.NewZapLogger(config.Logging.Level)
	assert.NoError(t, err, "Failed to initialize logger")

	// Initialize load balancer.
	lb, err := loadbalancer.NewRoundRobin(config, zapLogger)
	assert.NoError(t, err, "Failed to initialize load balancer")

	// Initialize reverse proxy handler.
	proxyHandler, err := proxy.NewReverseProxy(lb, zapLogger, config)
	assert.NoError(t, err, "Failed to initialize proxy handler")

	// Setup routes.
	router := proxyHandler.SetupRoutes()

	// Initialize middlewares.
	loggingMiddleware := middleware.LoggingMiddleware(zapLogger)
	corsMiddleware := middleware.CORSMiddleware(&config.CORS, zapLogger)
	rateLimiterMiddleware := middleware.NewRateLimiter(&config.RateLimit, zapLogger).Middleware()

	// Chain middlewares.
	chainedHandler := middleware.Chain(router, loggingMiddleware, corsMiddleware, rateLimiterMiddleware)

	// Create HTTP server.
	httpServer := &http.Server{
		Addr:         config.Server.Address,
		Handler:      chainedHandler,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	// Start server in a goroutine.
	go func() {
		wg.Done()
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("ListenAndServe failed", zap.Error(err))
		}
	}()

	// Wait for server to start.
	wg.Wait()
	// Wait a brief moment to ensure server readiness.
	time.Sleep(100 * time.Millisecond)

	// Teardown function.
	teardown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}

	return httpServer, teardown
}
