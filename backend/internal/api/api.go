// Package api wires the HTTP handlers of the server to their dependencies.
package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"inventory-backend/internal/auth"
)

// Server holds HTTP handlers and their dependencies.
type Server struct {
	pool   *pgxpool.Pool
	tokens *auth.TokenManager
	start  time.Time
}

// New registers all routes on a new ServeMux and returns it. Auth endpoints
// are public except logout, which requires a valid token. Later tasks attach
// category/product/stock/dashboard routes here behind s.requireAuth.
func New(pool *pgxpool.Pool, tokens *auth.TokenManager) *http.ServeMux {
	s := &Server{pool: pool, tokens: tokens, start: time.Now()}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /api/v1/auth/register", s.handleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", s.handleLogin)
	mux.Handle("POST /api/v1/auth/logout", s.requireAuth(http.HandlerFunc(s.handleLogout)))
	mux.Handle("GET /api/v1/categories", s.requireAuth(http.HandlerFunc(s.handleListCategories)))
	mux.Handle("POST /api/v1/categories", s.requireAuth(http.HandlerFunc(s.handleCreateCategory)))
	mux.Handle("GET /api/v1/categories/{id}", s.requireAuth(http.HandlerFunc(s.handleGetCategory)))
	mux.Handle("PUT /api/v1/categories/{id}", s.requireAuth(http.HandlerFunc(s.handleUpdateCategory)))
	mux.Handle("DELETE /api/v1/categories/{id}", s.requireAuth(http.HandlerFunc(s.handleDeleteCategory)))
	mux.Handle("GET /api/v1/products", s.requireAuth(http.HandlerFunc(s.handleListProducts)))
	mux.Handle("POST /api/v1/products", s.requireAuth(http.HandlerFunc(s.handleCreateProduct)))
	mux.Handle("GET /api/v1/products/{id}", s.requireAuth(http.HandlerFunc(s.handleGetProduct)))
	mux.Handle("PUT /api/v1/products/{id}", s.requireAuth(http.HandlerFunc(s.handleUpdateProduct)))
	mux.Handle("POST /api/v1/products/{id}/deactivate", s.requireAuth(http.HandlerFunc(s.handleDeactivateProduct)))
	mux.Handle("POST /api/v1/transactions", s.requireAuth(http.HandlerFunc(s.handleCreateMovement)))
	mux.Handle("GET /api/v1/transactions", s.requireAuth(http.HandlerFunc(s.handleListMovements)))
	mux.Handle("GET /api/v1/dashboard/summary", s.requireAuth(http.HandlerFunc(s.handleDashboardSummary)))
	return mux
}

type healthResponse struct {
	Status        string `json:"status"`
	Database      string `json:"database"`
	UptimeSeconds int64  `json:"uptime_seconds"`
	Time          string `json:"time"`
}

// handleHealth reports process health and live database connectivity. A slow
// or failed ping degrades the response to HTTP 503 so callers get an
// accurate picture.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	resp := healthResponse{
		Status:        "ok",
		Database:      "up",
		UptimeSeconds: int64(time.Since(s.start).Seconds()),
		Time:          time.Now().UTC().Format(time.RFC3339),
	}
	code := http.StatusOK
	if err := s.pool.Ping(ctx); err != nil {
		log.Printf("health: database ping failed: %v", err)
		resp.Status = "degraded"
		resp.Database = "down"
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}
