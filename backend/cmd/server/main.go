// Command server runs the Inventory and Stock Management API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inventory-backend/internal/api"
	"inventory-backend/internal/auth"
	"inventory-backend/internal/config"
	"inventory-backend/internal/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// .env lives in the repository root; support running from either the
	// root or the backend/ directory. OS environment variables always win.
	cfg, err := config.Load(".env", "../.env")
	if err != nil {
		return err
	}
	log.Printf("config: loaded (db=%s user=%s host=%s port=%s sslmode=%s app_port=%s jwt_ttl=24h)",
		cfg.Database, cfg.User, cfg.Host, cfg.Port, cfg.SSLMode, cfg.AppPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()
	log.Print("database: connection pool ready")

	tokens := auth.NewTokenManager(cfg.JWTSecret, 24*time.Hour) // expiry per docs/architecture.md Q6

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           api.New(pool, tokens),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("http: listening on :%s", cfg.AppPort)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Print("shutdown: signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	log.Print("shutdown: complete")
	return nil
}
