// Package database builds the application's PostgreSQL connection pool.
package database

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"inventory-backend/internal/config"
)

// NewPool assembles the DSN from cfg, creates a bounded connection pool, and
// verifies connectivity with a short ping. The DSN contains the password and
// is never logged; wrapped errors carry only pgx's own diagnostics (host,
// port, database, user — no password).
func NewPool(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(DSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("database: parse config: %w", err)
	}
	// Deliberately small pool: single-instance local development.
	poolCfg.MaxConns = 5
	poolCfg.MinConns = 1
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("database: create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database: ping: %w", err)
	}
	return pool, nil
}

// DSN builds the connection string from cfg. The result contains the
// password — callers must never log it. Exported for the migrate command.
func DSN(cfg config.Config) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   cfg.Host + ":" + cfg.Port,
		Path:   cfg.Database,
	}
	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}
