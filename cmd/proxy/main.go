package main

import (
	"context"
	"http-reverse-proxy/internal/loadbalancer"
	"http-reverse-proxy/internal/proxy"
	"http-reverse-proxy/pkg/logger"
	"http-reverse-proxy/pkg/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	lb, err := loadbalancer.NewRoundRobin(config.Backends)
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

	httpServer := &http.Server{
		Addr:         config.Server.Address,
		Handler:      router,
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

	gracefulShutdown(httpServer, zapLogger)
}

// gracefulShutdown handles server shutdown upon receiving termination signals
func gracefulShutdown(server *http.Server, logger *zap.Logger) {
	// Create a channel to listen for OS signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-signalChannel
	logger.Info("Shutdown signal received")

	// Create a context with timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting gracefully")
}
