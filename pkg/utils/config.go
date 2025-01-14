// pkg/utils/config.go
package utils

import (
	"errors"
	"fmt"
	"http-reverse-proxy/pkg/models"

	"github.com/spf13/viper"
)

// LoadConfig reads config.yaml and unmarshals it into Config struct
func LoadConfig(path string) (*models.Config, error) {
	v := viper.New() // Create a new Viper instance
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config models.Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	return &config, nil
}

// LoadBackendConfig is used to load backend pool configs for testing
func LoadBackendConfig(path string) (*models.BackendConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config models.BackendConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshalling config: %w", err)
	}

	return &config, nil
}

func ValidateConfig(cfg *models.Config) error {
	if cfg.Server.Address == "" {
		return errors.New("server address is required")
	}

	if len(cfg.Backends) == 0 {
		return errors.New("at least one backend is required")
	}

	if cfg.RateLimit.RequestsPerMinute <= 0 {
		return errors.New("rate_limit.requests_per_minute must be positive")
	}

	if cfg.RateLimit.Burst <= 0 {
		return errors.New("rate_limit.burst must be positive")
	}

	return nil
}
