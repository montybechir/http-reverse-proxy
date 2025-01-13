// pkg/utils/config.go
package utils

import (
	"errors"
	"http-reverse-proxy/pkg/models"

	"github.com/spf13/viper"
)

// LoadConfig reads config.yaml and unmarshals it into Config struct
func LoadConfig(path string) (*models.Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Read in environment variables that match
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config models.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
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

	// Add more validation rules as necessary

	return nil
}

func NewCORSConfig() *models.CORSConfig {
	return &models.CORSConfig{
		AllowedOrigins:   []string{},
		AllowCredentials: true,
		MaxAge:           300,
	}
}
