package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestStockMovementIntegration covers TC-36…TC-41, TC-47/TC-48 and the
// cross-user pattern against the real database, finishing with TC-45: a
// genuine concurrent burst over real HTTP against a boundary-stock product.
func TestStockMovementIntegration(t *testing.T) {
	s, pool, do := integrationEnv(t)
	ctx := context.Background()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)

	var userA, userB, catA, p1, p3 int64
	id, err := s.testInsertUser(ctx, "stock-a-"+suffix+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	userA = id
	id, err = s.testInsertUser(ctx, "stock-b-"+suffix+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	userB = id
	if err := pool.QueryRow(ctx,
		`INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id`,
		userA, "StockCat-"+suffix).Scan(&catA); err != nil {
		t.Fatal(err)
	}
	mkProduct := func(name, stock, min string) int64 {
		var pid int64
		if err := pool.QueryRow(ctx,
			`INSERT INTO products (user_id, category_id, name, sku, unit_price, opening_stock, current_stock, min_stock)
			 VALUES ($1, $2, $3, $4, 1, $5, $5, $6) RETURNING id`,
			userA, catA, name, "STK-"+name+"-"+suffix, stock, min).Scan(&pid); err != nil {
			t.Fatalf("fixture product: %v", err)
		}
		return pid
	}
	p1 = mkProduct("P1", "10", "3") // TC-36/37/38/41/48
	p3 = mkProduct("P3", "7", "1")  // deactivation check

	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM stock_transactions WHERE user_id = ANY($1)`, fmt.Sprintf("{%d,%d}", userA, userB))
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE user_id = ANY($1)`, fmt.Sprintf("{%d,%d}", userA, userB))
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE user_id = ANY($1)`, fmt.Sprintf("{%d,%d}", userA, userB))
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = ANY($1)`, fmt.Sprintf("{%d,%d}", userA, userB))
	})

	tokenA, _, err := s.tokens.Issue(userA, "a@example.test")
	if err != nil {
		t.Fatal(err)
	}
	tokenB, _, err := s.tokens.Issue(userB, "b@example.test")
	if err != nil {
		t.Fatal(err)
	}

	move := func(token, body string) *httptest.ResponseRecorder {
		return do("POST", "/api/v1/transactions", body, token)
	}
	stockOf := func(productID int64) int {
		rr := do("GET", fmt.Sprintf("/api/v1/products/%d", productID), "", tokenA)
		if rr.Code != http.StatusOK {
			t.Fatalf("product fetch: code = %d", rr.Code)
		}
		var p struct {
			CurrentStock int `json:"current_stock"`
		}
		decodeTestJSON(t, rr, &p)
		return p.CurrentStock
	}
	movementCount := func(productID int64, typ string) int {
		var n int
		if err := pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM stock_transactions WHERE product_id = $1 AND type = $2`,
			productID, typ).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}

	// -- TC-36: Stock In 5 on P1 (10 -> 15), remarks stored ------------------
	rr := move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":"IN","quantity":5,"remarks":"restock"}`, p1))
	if rr.Code != http.StatusCreated {
		t.Fatalf("stock in: code = %d, want 201 (body: %s)", rr.Code, rr.Body.String())
	}
	var mv struct {
		ID           int64  `json:"id"`
		ProductID    int64  `json:"product_id"`
		Type         string `json:"type"`
		Quantity     int    `json:"quantity"`
		CurrentStock int    `json:"current_stock"`
	}
	decodeTestJSON(t, rr, &mv)
	if mv.Type != "IN" || mv.Quantity != 5 || mv.CurrentStock != 15 || mv.ProductID != p1 {
		t.Fatalf("stock in response unexpected: %+v", mv)
	}
	if got := stockOf(p1); got != 15 {
		t.Fatalf("stock after IN = %d, want 15 (FR-4.3)", got)
	}

	// -- TC-37: Stock Out 4 on P1 (15 -> 11) ----------------------------------
	rr = move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":"OUT","quantity":4}`, p1))
	if rr.Code != http.StatusCreated {
		t.Fatalf("stock out: code = %d, want 201", rr.Code)
	}
	decodeTestJSON(t, rr, &mv)
	if mv.CurrentStock != 11 {
		t.Fatalf("stock after OUT = %d, want 11", mv.CurrentStock)
	}

	// -- TC-38 / TC-47: insufficient Stock Out rejected; nothing recorded -----
	before := movementCount(p1, "OUT")
	rr = move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":"OUT","quantity":999}`, p1))
	if rr.Code != http.StatusConflict {
		t.Fatalf("insufficient out: code = %d, want 409", rr.Code)
	}
	var e struct {
		Error string `json:"error"`
	}
	decodeTestJSON(t, rr, &e)
	if e.Error != "insufficient stock" {
		t.Fatalf("error message = %q", e.Error)
	}
	if got := stockOf(p1); got != 11 {
		t.Fatalf("stock changed by a failed OUT: %d, want 11 (AC-10)", got)
	}
	if after := movementCount(p1, "OUT"); after != before {
		t.Fatalf("failed OUT wrote rows: %d -> %d (atomicity violated)", before, after)
	}

	// -- TC-39/TC-40: invalid quantities rejected for both types --------------
	for _, tc := range []struct {
		typ string
		qty string
	}{
		{"IN", "0"}, {"IN", "-3"}, {"IN", "2.5"},
		{"OUT", "0"}, {"OUT", "-3"}, {"OUT", "2.5"},
	} {
		t.Run("invalid qty "+tc.typ+" "+tc.qty, func(t *testing.T) {
			rr := move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":%q,"quantity":%s}`, p1, tc.typ, tc.qty))
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("code = %d, want 400", rr.Code)
			}
		})
	}

	// -- TC-41: recorded rows contain every required field --------------------
	var (
		typ     string
		qty     int
		remarks *string
		created string
		count   int
	)
	if err := pool.QueryRow(ctx,
		`SELECT type, quantity, remarks, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SSOF'), COUNT(*) OVER ()
		 FROM stock_transactions WHERE product_id = $1 AND remarks IS NOT NULL LIMIT 1`, p1).Scan(
		&typ, &qty, &remarks, &created, &count); err != nil {
		t.Fatalf("TC-41 row fetch: %v", err)
	}
	if typ != "IN" || qty != 5 || remarks == nil || *remarks != "restock" || created == "" {
		t.Fatalf("TC-41 row incomplete: type=%s qty=%d remarks=%v created=%q", typ, qty, remarks, created)
	}

	// -- TC-48: remarks optional (NULL when omitted) ---------------------------
	rr = move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":"OUT","quantity":1}`, p1))
	if rr.Code != http.StatusCreated {
		t.Fatalf("out without remarks: code = %d", rr.Code)
	}
	var nullRemarks *string
	if err := pool.QueryRow(ctx,
		`SELECT remarks FROM stock_transactions WHERE product_id = $1 AND remarks IS NULL LIMIT 1`,
		p1).Scan(&nullRemarks); err != nil {
		t.Fatalf("omitted remarks not stored as NULL: %v", err)
	}

	// -- AC-02: user B cannot move user A's product ----------------------------
	rr = move(tokenB, fmt.Sprintf(`{"product_id":%d,"type":"OUT","quantity":1}`, p1))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("cross-user movement: code = %d, want 404", rr.Code)
	}

	// -- inactive product rejected ---------------------------------------------
	rr = do("POST", fmt.Sprintf("/api/v1/products/%d/deactivate", p3), "", tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("deactivate P3: code = %d", rr.Code)
	}
	rr = move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":"IN","quantity":1}`, p3))
	if rr.Code != http.StatusConflict {
		t.Fatalf("movement on inactive product: code = %d, want 409", rr.Code)
	}

	// -- TC-45 ⚠: concurrent Stock Outs at the boundary (real HTTP burst) ------
	// Fresh product at 0; IN 5 brings it to the boundary of exactly 5. Then
	// 10 simultaneous real-HTTP OUTs of 1: exactly 5 must succeed, 5 must
	// conflict, final stock 0, and one recorded row per successful movement.
	pBoundary := mkProduct("Boundary", "0", "0")
	rr = move(tokenA, fmt.Sprintf(`{"product_id":%d,"type":"IN","quantity":5}`, pBoundary))
	if rr.Code != http.StatusCreated {
		t.Fatalf("boundary top-up: code = %d, want 201", rr.Code)
	}
	if got := stockOf(pBoundary); got != 5 {
		t.Fatalf("boundary stock = %d, want exactly 5 before the burst", got)
	}
	insBefore := movementCount(pBoundary, "IN")
	outsBefore := movementCount(pBoundary, "OUT")

	// the burst: 10 concurrent real HTTP requests against a real server
	mux := New(pool, s.tokens)
	ts := httptest.NewServer(mux)
	defer ts.Close()
	client := &http.Client{Timeout: 15 * time.Second}

	const burst = 10
	var (
		mu      sync.Mutex
		ok201   int
		conf409 int
		other   int
	)
	var wg sync.WaitGroup
	for i := 0; i < burst; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/v1/transactions",
				strings.NewReader(fmt.Sprintf(`{"product_id":%d,"type":"OUT","quantity":1}`, pBoundary)))
			if err != nil {
				mu.Lock()
				other++
				mu.Unlock()
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+tokenA)
			res, err := client.Do(req)
			if err != nil {
				mu.Lock()
				other++
				mu.Unlock()
				return
			}
			defer res.Body.Close()
			mu.Lock()
			defer mu.Unlock()
			switch res.StatusCode {
			case http.StatusCreated:
				ok201++
			case http.StatusConflict:
				conf409++
			default:
				other++
			}
		}()
	}
	wg.Wait()

	if ok201 != 5 || conf409 != 5 || other != 0 {
		t.Fatalf("concurrent burst: 201=%d 409=%d other=%d; want exactly 5/5/0 (AC-11)", ok201, conf409, other)
	}
	if got := stockOf(pBoundary); got != 0 {
		t.Fatalf("stock after burst = %d, want 0 (never negative, AC-11)", got)
	}
	if outs := movementCount(pBoundary, "OUT"); outs != outsBefore+5 {
		t.Fatalf("recorded OUT rows = %d, want %d (one row per successful movement only)", outs, outsBefore+5)
	}
	if ins := movementCount(pBoundary, "IN"); ins != insBefore {
		t.Fatalf("recorded IN rows = %d, want %d (the burst contains only OUTs)", ins, insBefore)
	}
	t.Logf("TC-45 concurrency burst: %d created, %d conflicted, final stock 0, movements recorded correctly", ok201, conf409)
}
