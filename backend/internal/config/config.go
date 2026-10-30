// Package config loads application configuration from the environment.
//
// Values are read from OS environment variables. A local .env file
// (git-ignored) is applied first for any variable not already present in the
// process environment, so the OS environment always takes precedence
// (12-factor style) and real secrets stay out of the repository and shell
// history. Error messages name variables, never values.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime settings for the server.
type Config struct {
	Host      string // PGHOST
	Port      string // PGPORT
	Database  string // PGDATABASE
	User      string // PGUSER
	Password  string // PGPASSWORD (never logged or included in errors)
	SSLMode   string // PGSSLMODE
	AppPort   string // APP_PORT (HTTP listen port)
	JWTSecret string // JWT_SECRET (HS256 signing key; never logged)
}

var validSSLModes = map[string]bool{
	"disable": true, "allow": true, "prefer": true,
	"require": true, "verify-ca": true, "verify-full": true,
}

// Load applies the given optional .env files (missing files are skipped) and
// then reads and validates the configuration from the environment.
func Load(envPaths ...string) (Config, error) {
	for _, p := range envPaths {
		applyDotEnv(p)
	}

	cfg := Config{
		Host:      os.Getenv("PGHOST"),
		Port:      os.Getenv("PGPORT"),
		Database:  os.Getenv("PGDATABASE"),
		User:      os.Getenv("PGUSER"),
		Password:  os.Getenv("PGPASSWORD"),
		SSLMode:   os.Getenv("PGSSLMODE"),
		AppPort:   os.Getenv("APP_PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
	cfg.setDefaults()
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) setDefaults() {
	if c.SSLMode == "" {
		c.SSLMode = "disable" // local development default (docs/architecture.md §10)
	}
	if c.AppPort == "" {
		c.AppPort = "8080" // docs/architecture.md §10 dev port
	}
}

func (c Config) validate() error {
	required := []struct{ name, val string }{
		{"PGHOST", c.Host},
		{"PGPORT", c.Port},
		{"PGDATABASE", c.Database},
		{"PGUSER", c.User},
		{"PGPASSWORD", c.Password},
		{"JWT_SECRET", c.JWTSecret},
	}
	for _, r := range required {
		if strings.TrimSpace(r.val) == "" {
			return fmt.Errorf("config: %s is required", r.name)
		}
	}
	// A short HS256 secret is brute-forceable; 32 bytes is the working minimum.
	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("config: JWT_SECRET must be at least 32 characters")
	}
	if n, err := strconv.Atoi(c.Port); err != nil || n < 1 || n > 65535 {
		return fmt.Errorf("config: PGPORT must be an integer between 1 and 65535")
	}
	if n, err := strconv.Atoi(c.AppPort); err != nil || n < 1 || n > 65535 {
		return fmt.Errorf("config: APP_PORT must be an integer between 1 and 65535")
	}
	if !validSSLModes[c.SSLMode] {
		return fmt.Errorf("config: PGSSLMODE %q is not a valid libpq sslmode", c.SSLMode)
	}
	return nil
}

// applyDotEnv applies KEY=VALUE lines from path for variables not already set
// in the process environment. A missing file is not an error; malformed lines
// are skipped. Values are never logged.
func applyDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}
