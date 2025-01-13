package models

import "time"

type Config struct {
	Server    ServerConfig `mapstructure:"server"`
	Backends  []string     `mapstructure:"backends"`
	RateLimit struct {
		RequestsPerMinute int `mapstructure:"requests_per_minute"`
		Burst             int `mapstructure:"burst"`
	} `mapstructure:"rate_limit"`
	CORS    CORSConfig    `mapstructure:"cors"`
	Logging LoggingConfig `mapstructure:"logging"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
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

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
	Services  map[string]bool   `json:"services"`
}
