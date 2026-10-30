package api

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

// maxRemarksLen caps optional remarks (documented cap; unspecified by the
// assignment). maxStockInt mirrors the INTEGER column bound used for the
// overflow guard below.
const maxRemarksLen = 500

const maxStockInt = int(^uint32(0) >> 1) // 2147483647

type transaction struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	Type         string    `json:"type"`
	Quantity     int       `json:"quantity"`
	Remarks      *string   `json:"remarks"`
	CreatedAt    time.Time `json:"created_at"`
	CurrentStock int       `json:"current_stock"` // product balance after this movement
}

// positiveWhole validates movement quantities: a whole number > 0 that cannot
// overflow the INTEGER stock column (AC-09).
func positiveWhole(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0) && v > 0 && v <= float64(maxStockInt) && v == math.Trunc(v)
}

type movementInput struct {
	ProductID *int64   `json:"product_id"`
	Type      *string  `json:"type"`
	Quantity  *float64 `json:"quantity"`
	Remarks   *string  `json:"remarks"`
}

// handleCreateMovement records a Stock In/Out and applies it to the product
// balance ATOMICALLY (AC-12): the conditional balance UPDATE (AC-10/AC-11)
// and the immutable movement INSERT commit together or not at all.
//
// Concurrency: PostgreSQL serializes concurrent UPDATEs on the product row —
// the second request blocks until the first commits, then re-evaluates the
// sufficiency guard against the committed value. Losers get no row and are
// rejected with 409. AC-11 therefore holds under real concurrent load.
func (s *Server) handleCreateMovement(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var in movementInput
	if !decodeBody(w, r, &in) {
		return
	}
	if in.ProductID == nil || *in.ProductID <= 0 {
		writeError(w, http.StatusBadRequest, "a valid product is required")
		return
	}
	if in.Type == nil || (*in.Type != "IN" && *in.Type != "OUT") {
		writeError(w, http.StatusBadRequest, "type must be IN or OUT")
		return
	}
	if in.Quantity == nil || !positiveWhole(*in.Quantity) {
		writeError(w, http.StatusBadRequest, "quantity must be a positive whole number")
		return
	}
	if in.Remarks != nil && len(*in.Remarks) > maxRemarksLen {
		writeError(w, http.StatusBadRequest, "remarks must be at most 500 characters")
		return
	}
	qty := int(*in.Quantity)
	delta := qty
	if *in.Type == "OUT" {
		delta = -qty
	}

	ctx := r.Context()

	// Fast-path state check: gives precise 404/409 semantics up front; the
	// transaction's conditional UPDATE below remains the race-safe authority.
	var status string
	var current int
	err := s.pool.QueryRow(ctx,
		`SELECT status, current_stock FROM products WHERE id = $1 AND user_id = $2`,
		*in.ProductID, userID).Scan(&status, &current)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	if err != nil {
		log.Printf("movements: product fetch: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if status != "active" {
		writeError(w, http.StatusConflict, "product is not active")
		return
	}
	if *in.Type == "OUT" && current < qty {
		writeError(w, http.StatusConflict, "insufficient stock")
		return
	}
	if *in.Type == "IN" && current > maxStockInt-qty {
		writeError(w, http.StatusConflict, "stock limit would be exceeded")
		return
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Printf("movements: begin: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Conditional balance update: the row lock plus the re-evaluated guard
	// make concurrent insufficient OUTs fail here instead of corrupting stock.
	// Two statement variants keep every parameter in a typed context (an
	// untyped `$6 - $5` fails with SQLSTATE 42725).
	var balance int
	var upd string
	var updArgs []any
	if *in.Type == "IN" {
		maxBalance := maxStockInt - qty // reject IN that would overflow INTEGER
		upd = `UPDATE products SET current_stock = current_stock + $1
		       WHERE id = $2 AND user_id = $3 AND status = 'active' AND current_stock <= $4
		       RETURNING current_stock`
		updArgs = []any{delta, *in.ProductID, userID, maxBalance}
	} else {
		upd = `UPDATE products SET current_stock = current_stock - $1
		       WHERE id = $2 AND user_id = $3 AND status = 'active' AND current_stock >= $1
		       RETURNING current_stock`
		updArgs = []any{qty, *in.ProductID, userID}
	}
	err = tx.QueryRow(ctx, upd, updArgs...).Scan(&balance)
	if errors.Is(err, pgx.ErrNoRows) { // race lost: balance changed under us
		_ = tx.Rollback(ctx)
		writeError(w, http.StatusConflict, "insufficient stock")
		return
	}
	if err != nil {
		log.Printf("movements: balance update: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var (
		id        int64
		createdAt time.Time
	)
	err = tx.QueryRow(ctx,
		`INSERT INTO stock_transactions (user_id, product_id, type, quantity, remarks)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		userID, *in.ProductID, *in.Type, qty, in.Remarks).Scan(&id, &createdAt)
	if err != nil {
		log.Printf("movements: insert: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if err := tx.Commit(ctx); err != nil {
		log.Printf("movements: commit: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, transaction{
		ID: id, ProductID: *in.ProductID, Type: *in.Type, Quantity: qty,
		Remarks: in.Remarks, CreatedAt: createdAt, CurrentStock: balance,
	})
}

// handleListMovements serves the stock transaction history (FR-4.4, T-204):
// only the authenticated user's rows, newest first, filterable by product and
// type, with a bounded page size (Q8). Rows are immutable (AC-13) — history
// can only grow.
func (s *Server) handleListMovements(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	q := r.URL.Query()

	var productFilter int64
	if raw := q.Get("product_id"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || v <= 0 {
			writeError(w, http.StatusBadRequest, "invalid product_id filter")
			return
		}
		productFilter = v
	}
	typeFilter := q.Get("type")
	if typeFilter != "" && typeFilter != "IN" && typeFilter != "OUT" {
		writeError(w, http.StatusBadRequest, "type must be IN or OUT")
		return
	}
	limit := 100
	if raw := q.Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 1 || v > 500 {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 500")
			return
		}
		limit = v
	}

	query := `SELECT id, product_id, type, quantity, remarks, created_at
	          FROM stock_transactions WHERE user_id = $1`
	args := []any{userID}
	if productFilter != 0 {
		args = append(args, productFilter)
		query += fmt.Sprintf(" AND product_id = $%d", len(args))
	}
	if typeFilter != "" {
		args = append(args, typeFilter)
		query += fmt.Sprintf(" AND type = $%d", len(args))
	}
	args = append(args, limit)
	query += fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", len(args))

	rows, err := s.pool.Query(r.Context(), query, args...)
	if err != nil {
		log.Printf("movements: list: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer rows.Close()

	list := []transaction{}
	for rows.Next() {
		var t transaction
		if err := rows.Scan(&t.ID, &t.ProductID, &t.Type, &t.Quantity, &t.Remarks, &t.CreatedAt); err != nil {
			log.Printf("movements: list scan: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		list = append(list, t)
	}
	if err := rows.Err(); err != nil {
		log.Printf("movements: list rows: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"transactions": list})
}
