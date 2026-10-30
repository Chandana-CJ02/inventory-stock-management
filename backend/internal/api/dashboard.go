package api

import (
	"log"
	"net/http"
	"time"
)

type dashboardTotals struct {
	ActiveProducts int `json:"active_products"`
	Categories     int `json:"categories"`
	LowStock       int `json:"low_stock"`
}

type lowStockProduct struct {
	ID           int64  `json:"id"`
	CategoryID   int64  `json:"category_id"`
	Name         string `json:"name"`
	SKU          string `json:"sku"`
	CurrentStock int    `json:"current_stock"`
	MinStock     int    `json:"min_stock"`
}

type recentTransaction struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Type      string    `json:"type"`
	Quantity  int       `json:"quantity"`
	Remarks   *string   `json:"remarks"`
	CreatedAt time.Time `json:"created_at"`
}

// handleDashboardSummary serves the dashboard (FR-5.1…FR-5.3, US-17): totals
// for the authenticated user's inventory, the low-stock list, and the five
// most recent stock transactions. Every query is scoped by the JWT-context
// user ID, so a dashboard can only ever show its owner's data (AC-02).
//
// Low stock means an ACTIVE product whose current stock is at or below its
// minimum (FR-5.2; active-only interpretation per docs/requirements.md Q2).
// The partial index products_low_stock_idx covers this predicate.
func (s *Server) handleDashboardSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	ctx := r.Context()

	var totals dashboardTotals
	err := s.pool.QueryRow(ctx, `
		SELECT
		  (SELECT COUNT(*) FROM products  WHERE user_id = $1 AND status = 'active'),
		  (SELECT COUNT(*) FROM categories WHERE user_id = $1),
		  (SELECT COUNT(*) FROM products  WHERE user_id = $1 AND status = 'active' AND current_stock <= min_stock)
		`, userID).Scan(&totals.ActiveProducts, &totals.Categories, &totals.LowStock)
	if err != nil {
		log.Printf("dashboard: totals: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	low := []lowStockProduct{}
	lowRows, err := s.pool.Query(ctx, `
		SELECT id, category_id, name, sku, current_stock, min_stock
		FROM products
		WHERE user_id = $1 AND status = 'active' AND current_stock <= min_stock
		ORDER BY (current_stock - min_stock), name`, userID)
	if err != nil {
		log.Printf("dashboard: low stock: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer lowRows.Close()
	for lowRows.Next() {
		var p lowStockProduct
		if err := lowRows.Scan(&p.ID, &p.CategoryID, &p.Name, &p.SKU, &p.CurrentStock, &p.MinStock); err != nil {
			log.Printf("dashboard: low stock scan: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		low = append(low, p)
	}
	if err := lowRows.Err(); err != nil {
		log.Printf("dashboard: low stock rows: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	recent := []recentTransaction{}
	txRows, err := s.pool.Query(ctx, `
		SELECT id, product_id, type, quantity, remarks, created_at
		FROM stock_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT 5`, userID)
	if err != nil {
		log.Printf("dashboard: recent: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	defer txRows.Close()
	for txRows.Next() {
		var t recentTransaction
		if err := txRows.Scan(&t.ID, &t.ProductID, &t.Type, &t.Quantity, &t.Remarks, &t.CreatedAt); err != nil {
			log.Printf("dashboard: recent scan: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		recent = append(recent, t)
	}
	if err := txRows.Err(); err != nil {
		log.Printf("dashboard: recent rows: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"totals":              totals,
		"low_stock":           low,
		"recent_transactions": recent,
	})
}
