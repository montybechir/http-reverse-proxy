package main

import (
	"http-reverse-proxy/internal/loadbalancer"
	"http-reverse-proxy/internal/middleware"
	"http-reverse-proxy/internal/proxy"
	"http-reverse-proxy/pkg/logger"
	"http-reverse-proxy/pkg/server"
	"http-reverse-proxy/pkg/utils"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {

	// Load configuration
	// ctx := context.Background()

	config, err := utils.LoadConfig("configs/config.yaml")

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize logger based on config
	zapLogger, err := logger.NewZapLogger(config.Logging.Level)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	defer zapLogger.Sync() // Flushes buffer, if any

	// Validate configuration
	if err := utils.ValidateConfig(config); err != nil {
		zapLogger.Fatal("Invalid configuration", zap.Error(err))
	}

	// Initialize load balancer
	lb, err := loadbalancer.NewRoundRobin(config, zapLogger)
	if err != nil {
		zapLogger.Fatal("Failed to initialize load balancer", zap.Error(err))
	}

	// Initialize the reverse proxy handler
	proxyHandler, err := proxy.NewReverseProxy(lb, zapLogger, config)
	if err != nil {
		zapLogger.Fatal("Failed to initialize proxy", zap.Error(err))
	}

	// setup routes with handlers and middleware
	router := proxyHandler.SetupRoutes()

	loggingMiddleware := middleware.LoggingMiddleware(zapLogger)
	corsMiddleware := middleware.CORSMiddleware(&config.CORS, zapLogger)
	rateLimiterMiddleware := middleware.NewRateLimiter(&config.RateLimit, zapLogger).Middleware()

	chainedHandler := middleware.Chain(router, loggingMiddleware, corsMiddleware, rateLimiterMiddleware)

	httpServer := &http.Server{
		Addr:         config.Server.Address,
		Handler:      chainedHandler,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}

	// start server in a goroutine
	go func() {
		zapLogger.Info("Starting server", zap.String("address", config.Server.Address))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("ListenAndServe error:", zap.Error(err))
		}
	}()

	server.GracefulShutdown(httpServer, zapLogger)
}
