package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	maxProductNameLen = 200 // documented cap (unspecified by the assignment)
	maxProductSKULen  = 100
	maxPrice          = float64(99999999.99) // NUMERIC(10,2) upper bound
	maxStock          = float64(2147483647)  // INTEGER column upper bound
)

type product struct {
	ID           int64     `json:"id"`
	CategoryID   int64     `json:"category_id"`
	Name         string    `json:"name"`
	SKU          string    `json:"sku"`
	UnitPrice    float64   `json:"unit_price"`
	OpeningStock int       `json:"opening_stock"`
	CurrentStock int       `json:"current_stock"`
	MinStock     int       `json:"min_stock"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const productColumns = `id, category_id, name, sku, unit_price, opening_stock, current_stock, min_stock, status, created_at, updated_at`

func scanProduct(row pgx.Row) (product, error) {
	var p product
	err := row.Scan(&p.ID, &p.CategoryID, &p.Name, &p.SKU, &p.UnitPrice,
		&p.OpeningStock, &p.CurrentStock, &p.MinStock, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// productInput carries client-supplied fields. Pointers distinguish "absent"
// (PUT keeps the stored value) from "present" (validated and applied); JSON
// null behaves like absent for required fields. current_stock is deliberately
// absent — stock changes only through stock movements (AC-12/AC-13).
type productInput struct {
	Name         *string  `json:"name"`
	SKU          *string  `json:"sku"`
	CategoryID   *int64   `json:"category_id"`
	UnitPrice    *float64 `json:"unit_price"`
	OpeningStock *float64 `json:"opening_stock"`
	MinStock     *float64 `json:"min_stock"`
	Status       *string  `json:"status"`
}

// wholeNonNegative validates a stock value: a non-negative whole number that
// fits the INTEGER column. Fractions (2.5), negatives, and non-numbers are
// rejected (AC-05).
func wholeNonNegative(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 && v <= maxStock && v == math.Trunc(v)
}

func validPrice(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 && v <= maxPrice
}

// requireProductFields validates a create request (all fields required) and
// returns the normalized values plus a user-facing error message (AC-16).
func requireProductFields(in productInput) (name, sku string, categoryID int64, price float64, opening, minStock int, errMsg string) {
	if in.Name == nil || strings.TrimSpace(*in.Name) == "" {
		return "", "", 0, 0, 0, 0, "product name is required"
	}
	name = strings.TrimSpace(*in.Name)
	if len(name) > maxProductNameLen {
		return "", "", 0, 0, 0, 0, "product name must be at most 200 characters"
	}
	if in.SKU == nil || strings.TrimSpace(*in.SKU) == "" {
		return "", "", 0, 0, 0, 0, "sku is required"
	}
	sku = strings.TrimSpace(*in.SKU)
	if len(sku) > maxProductSKULen {
		return "", "", 0, 0, 0, 0, "sku must be at most 100 characters"
	}
	if in.CategoryID == nil || *in.CategoryID <= 0 {
		return "", "", 0, 0, 0, 0, "a valid category is required"
	}
	categoryID = *in.CategoryID
	if in.UnitPrice == nil || !validPrice(*in.UnitPrice) {
		return "", "", 0, 0, 0, 0, "unit price must be a non-negative number up to 99999999.99"
	}
	price = *in.UnitPrice
	if in.OpeningStock == nil || !wholeNonNegative(*in.OpeningStock) {
		return "", "", 0, 0, 0, 0, "opening stock must be a non-negative whole number"
	}
	opening = int(*in.OpeningStock)
	if in.MinStock == nil || !wholeNonNegative(*in.MinStock) {
		return "", "", 0, 0, 0, 0, "minimum stock must be a non-negative whole number"
	}
	minStock = int(*in.MinStock)
	return name, sku, categoryID, price, opening, minStock, ""
}

func productID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

// categoryOwnedBy reports whether the category exists and belongs to the user
// (AC-07: a product cannot be created without a valid own category).
func (s *Server) categoryOwnedBy(ctx context.Context, categoryID, userID int64) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM categories WHERE id = $1 AND user_id = $2)`,
		categoryID, userID).Scan(&exists)
	return exists, err
}

// handleListProducts returns only the authenticated user's products with their
// current available stock (FR-3.5, AC-02). Optional filters (T-202):
//   - ?search=…  case-insensitive substring match on product name OR SKU (Q7)
//   - ?category_id=n  restrict to one of the user's categories (FR-3.4)
//
// Filters are appended as bound parameters — never string-interpolated (SQL-
// injection safe).
func (s *Server) handleListProducts(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	q := r.URL.Query()
	search := strings.TrimSpace(q.Get("search"))
	if len(search) > 100 {
		writeError(w, http.StatusBadRequest, "search query is too long")
		return
	}
	var categoryFilter int64
	if raw := q.Get("category_id"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || v <= 0 {
			writeError(w, http.StatusBadRequest, "invalid category_id filter")
			return
		}
		categoryFilter = v
	}

	query := `SELECT ` + productColumns + ` FROM products WHERE user_id = $1`
	args := []any{userID}
	if search != "" {
		args = append(args, "%"+strings.ToLower(search)+"%")
		query += fmt.Sprintf(" AND (lower(name) LIKE $%d OR lower(sku) LIKE $%d)", len(args), len(args))
	}
	if categoryFilter != 0 {
		args = append(args, categoryFilter)
		query += fmt.Sprintf(" AND category_id = $%d", len(args))
	}
	query += " ORDER BY created_at DESC, id DESC"

	rows, err := s.pool.Query(r.Context(), query, args...)
	if err != nil {
		log.Printf("products: list: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer rows.Close()

	list := []product{}
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			log.Printf("products: list scan: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		log.Printf("products: list rows: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"products": list})
}

// handleCreateProduct creates a product for the authenticated user. The
// category must exist and belong to the user (AC-07); the SKU must be unique
// within the user's inventory (AC-06); price and stock fields follow AC-04/05.
// current_stock starts at opening_stock; no stock movement is recorded for the
// opening balance (docs/requirements.md Q13).
func (s *Server) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var in productInput
	if !decodeBody(w, r, &in) {
		return
	}
	name, sku, categoryID, price, opening, minStock, errMsg := requireProductFields(in)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}

	owned, err := s.categoryOwnedBy(r.Context(), categoryID, userID)
	if err != nil {
		log.Printf("products: create category check: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !owned {
		writeError(w, http.StatusBadRequest, "category does not exist or is not yours")
		return
	}

	var p product
	err = s.pool.QueryRow(r.Context(),
		`INSERT INTO products (user_id, category_id, name, sku, unit_price, opening_stock, current_stock, min_stock)
		 VALUES ($1, $2, $3, $4, $5, $6, $6, $7)
		 RETURNING `+productColumns,
		userID, categoryID, name, sku, price, opening, minStock).Scan(
		&p.ID, &p.CategoryID, &p.Name, &p.SKU, &p.UnitPrice,
		&p.OpeningStock, &p.CurrentStock, &p.MinStock, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation (products_user_sku_unique)
			writeError(w, http.StatusConflict, "sku already exists in your inventory")
			return
		}
		log.Printf("products: create: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

// handleGetProduct returns one own product; foreign or missing IDs return 404
// (no existence leakage, TC-10 pattern).
func (s *Server) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := productID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	p, err := scanProduct(s.pool.QueryRow(r.Context(),
		`SELECT `+productColumns+` FROM products WHERE id = $1 AND user_id = $2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		log.Printf("products: get: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// handleUpdateProduct edits product details (FR-3.2). Absent fields keep
// their stored values; every merged value is re-validated. current_stock is
// never modified here — only stock movements change it. A category change
// must reference the user's own category (AC-07); SKU uniqueness is
// re-checked on edit (Q14).
func (s *Server) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := productID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	var in productInput
	if !decodeBody(w, r, &in) {
		return
	}

	existing, err := scanProduct(s.pool.QueryRow(r.Context(),
		`SELECT `+productColumns+` FROM products WHERE id = $1 AND user_id = $2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		log.Printf("products: update fetch: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Merge: provided fields replace stored values, then validate the whole.
	name, sku := existing.Name, existing.SKU
	categoryID, price := existing.CategoryID, existing.UnitPrice
	opening, minStock := existing.OpeningStock, existing.MinStock
	status := existing.Status
	if in.Name != nil {
		name = strings.TrimSpace(*in.Name)
	}
	if in.SKU != nil {
		sku = strings.TrimSpace(*in.SKU)
	}
	if in.CategoryID != nil {
		categoryID = *in.CategoryID
	}
	if in.UnitPrice != nil {
		price = *in.UnitPrice
	}
	if in.OpeningStock != nil {
		opening = int(*in.OpeningStock)
	}
	if in.MinStock != nil {
		minStock = int(*in.MinStock)
	}
	if in.Status != nil {
		status = *in.Status
	}
	if errMsg := validateProductFields(name, sku, categoryID, price, opening, minStock); errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg)
		return
	}
	if status != "active" && status != "inactive" {
		writeError(w, http.StatusBadRequest, "status must be active or inactive")
		return
	}
	if categoryID != existing.CategoryID {
		owned, err := s.categoryOwnedBy(r.Context(), categoryID, userID)
		if err != nil {
			log.Printf("products: update category check: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !owned {
			writeError(w, http.StatusBadRequest, "category does not exist or is not yours")
			return
		}
	}

	p, err := scanProduct(s.pool.QueryRow(r.Context(),
		`UPDATE products SET category_id = $1, name = $2, sku = $3, unit_price = $4,
		 opening_stock = $5, min_stock = $6, status = $7, updated_at = now()
		 WHERE id = $8 AND user_id = $9
		 RETURNING `+productColumns,
		categoryID, name, sku, price, opening, minStock, status, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			writeError(w, http.StatusConflict, "sku already exists in your inventory")
			return
		}
		log.Printf("products: update: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// handleDeactivateProduct sets the product inactive (AC-14). Products with
// transaction history are never hard-deleted — no delete route exists at all
// (TC-35); deactivation is idempotent.
func (s *Server) handleDeactivateProduct(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := productID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	p, err := scanProduct(s.pool.QueryRow(r.Context(),
		`UPDATE products SET status = 'inactive', updated_at = now()
		 WHERE id = $1 AND user_id = $2
		 RETURNING `+productColumns,
		id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		log.Printf("products: deactivate: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// validateProductFields validates merged product values for update (same
// rules as create; AC-04/05/07 and field sanity).
func validateProductFields(name, sku string, categoryID int64, price float64, opening, minStock int) string {
	if strings.TrimSpace(name) == "" {
		return "product name is required"
	}
	if len(name) > maxProductNameLen {
		return "product name must be at most 200 characters"
	}
	if strings.TrimSpace(sku) == "" {
		return "sku is required"
	}
	if len(sku) > maxProductSKULen {
		return "sku must be at most 100 characters"
	}
	if categoryID <= 0 {
		return "a valid category is required"
	}
	if !validPrice(price) {
		return "unit price must be a non-negative number up to 99999999.99"
	}
	if opening < 0 || float64(opening) > maxStock {
		return "opening stock must be a non-negative whole number"
	}
	if minStock < 0 || float64(minStock) > maxStock {
		return "minimum stock must be a non-negative whole number"
	}
	return ""
}
