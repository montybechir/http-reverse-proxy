package models

import "time"

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Backends    []string          `mapstructure:"backends"`
	RateLimit   RateLimitConfig   `mapstructure:"rate_limit"`
	CORS        CORSConfig        `mapstructure:"cors"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	HealthCheck HealthCheckConfig `mapstructure:"health_check"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type HealthCheckConfig struct {
	Frequency        time.Duration `mapstructure:"frequency"`
	Timeout          time.Duration `mapstructure:"timeout"`
	HealthyThreshold int           `mapstructure:"healthy_threshold"`
	Path             string        `mapstructure:"path"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
	Debug            bool     `mapstructure:"debug"`
}

type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	// Optional: Additional server configurations
	// MaxHeaderBytes int        `mapstructure:"max_header_bytes"`
	// TLSConfig      *TLSConfig `mapstructure:"tls,omitempty"`
}

type BackendServerConfig struct {
	Address  string `mapstructure:"address"`
	Response string `mapstructure:"response"`
}

// ServerConfig holds the configuration for the server.
type BackendConfig struct {
	Logging LoggingConfig       `mapstructure:"logging"`
	Server  BackendServerConfig `mapstructure:"server"`
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
	Services  map[string]bool   `json:"services"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}
