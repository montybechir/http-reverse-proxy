package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(logLevel string) (*zap.Logger, error) {
	var cfg zap.Config

	// Choose the config based on the log level
	switch strings.ToLower(logLevel) {
	case "development":
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	case "production":
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		// Default to production config if unknown level
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Parse log level
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(strings.ToLower(logLevel))); err != nil {
		return nil, err
	}
	cfg.Level = level

	// Build the logger
	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return zapLogger, nil
}
