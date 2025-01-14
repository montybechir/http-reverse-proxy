// pkg/server/server.go
package server

import (
	"context"
	"fmt"
	"http-reverse-proxy/internal/middleware"
	"http-reverse-proxy/pkg/logger"
	"http-reverse-proxy/pkg/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// gracefulShutdown handles server shutdown upon receiving termination signals
func GracefulShutdown(server *http.Server, logger *zap.Logger) error {
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
	return nil
}

// StartServer initializes and starts the HTTP server based on the provided configuration.
func StartServer(configPath string) error {
	// Load configuration
	config, err := utils.LoadBackendConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger based on config
	zapLogger, err := logger.NewZapLogger(config.Logging.Level)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		if err := zapLogger.Sync(); err != nil {
			zapLogger.Error("Error syncing logger", zap.Error(err))
		}
	}()

	// Initialize ServeMux and handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintf(w, config.Server.Response)
	})

	// Initialize Logging Middleware
	loggingMiddleware := middleware.LoggingMiddleware(zapLogger)

	// Chain middleware with the mux
	chainedHandler := middleware.Chain(mux, loggingMiddleware)

	// Create HTTP server
	srv := &http.Server{
		Addr:    config.Server.Address,
		Handler: chainedHandler,
	}

	// Start server in a goroutine
	fmt.Println("----Config", config, "ConfigPath:", configPath)
	go func() {
		zapLogger.Info("Server is starting", zap.String("address", config.Server.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	return GracefulShutdown(srv, zapLogger)
}
