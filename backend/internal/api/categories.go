package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// maxCategoryNameLen caps category names. The assignment does not specify a
// length; 100 is a documented working limit (docs/requirements.md §9 style).
const maxCategoryNameLen = 100

type category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func scanCategory(row pgx.Row) (category, error) {
	var c category
	err := row.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

// validateCategoryName returns the normalized name or a user-facing error
// message (AC-16: backend validation independent of the frontend).
func validateCategoryName(name string) (string, string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", "category name is required"
	}
	if len(name) > maxCategoryNameLen {
		return "", "category name must be at most 100 characters"
	}
	return name, ""
}

// categoryID extracts and validates the {id} path segment.
func categoryID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func decodeBody(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	return true
}

// handleListCategories returns only the authenticated user's categories
// (FR-2.4, AC-02). The user ID always comes from the verified JWT context.
func (s *Server) handleListCategories(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	rows, err := s.pool.Query(r.Context(),
		`SELECT id, name, created_at, updated_at FROM categories WHERE user_id = $1 ORDER BY name`,
		userID)
	if err != nil {
		log.Printf("categories: list: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer rows.Close()

	cats := []category{}
	for rows.Next() {
		c, err := scanCategory(rows)
		if err != nil {
			log.Printf("categories: list scan: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		cats = append(cats, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("categories: list rows: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": cats})
}

// handleCreateCategory creates a category for the authenticated user.
// Per-user name uniqueness is enforced by the categories_user_name_unique
// index; a violation maps to 409 (FR-2.2, TC-14).
func (s *Server) handleCreateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	name, errMsg := validateCategoryName(req.Name)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	var c category
	err := s.pool.QueryRow(r.Context(),
		`INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id, name, created_at, updated_at`,
		userID, name).Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			writeError(w, http.StatusConflict, "a category with this name already exists")
			return
		}
		log.Printf("categories: create: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, c)
}

// handleGetCategory returns one own category. Foreign or missing IDs return
// 404 (not 403) so the endpoint does not leak the existence of other users'
// records (docs/architecture.md §4, TC-10).
func (s *Server) handleGetCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := categoryID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	c, err := scanCategory(s.pool.QueryRow(r.Context(),
		`SELECT id, name, created_at, updated_at FROM categories WHERE id = $1 AND user_id = $2`,
		id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	if err != nil {
		log.Printf("categories: get: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// handleUpdateCategory renames one own category. Renaming onto the user's own
// other category name maps to 409 via the unique index.
func (s *Server) handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := categoryID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if !decodeBody(w, r, &req) {
		return
	}
	name, errMsg := validateCategoryName(req.Name)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	c, err := scanCategory(s.pool.QueryRow(r.Context(),
		`UPDATE categories SET name = $1, updated_at = now()
		 WHERE id = $2 AND user_id = $3
		 RETURNING id, name, created_at, updated_at`,
		name, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "a category with this name already exists")
			return
		}
		log.Printf("categories: update: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// handleDeleteCategory deletes one own category, blocked while any product
// references it (AC-08 plus the stricter Q3 rule: inactive products also
// block deletion so their category reference can never dangle). The guard is
// part of the DELETE's WHERE clause, making the check-and-delete atomic.
func (s *Server) handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := categoryID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	tag, err := s.pool.Exec(r.Context(),
		`DELETE FROM categories
		 WHERE id = $1 AND user_id = $2
		   AND NOT EXISTS (SELECT 1 FROM products WHERE products.category_id = categories.id)`,
		id, userID)
	if err != nil {
		log.Printf("categories: delete: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if tag.RowsAffected() == 1 {
		writeJSON(w, http.StatusOK, map[string]string{"message": "category deleted"})
		return
	}

	// Zero rows: either foreign/missing (404) or product-referenced (409).
	var exists bool
	if err := s.pool.QueryRow(r.Context(),
		`SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1 AND user_id = $2)`,
		id, userID).Scan(&exists); err != nil {
		log.Printf("categories: delete exists-check: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !exists {
		writeError(w, http.StatusNotFound, "category not found")
		return
	}
	writeError(w, http.StatusConflict, "category has products and cannot be deleted")
}
