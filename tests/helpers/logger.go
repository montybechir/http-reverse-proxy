package helpers

import (
	"http-reverse-proxy/pkg/logger"

	"go.uber.org/zap"
)

func NewTestLogger() (*zap.Logger, error) {
	return logger.NewZapLogger("debug")
}
