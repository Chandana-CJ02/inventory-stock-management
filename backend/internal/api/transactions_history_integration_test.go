package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// TestTransactionHistoryIntegration covers T-204 / TC-42, TC-43: history
// listing filtered by product and type, ordering, paging, and cross-user
// isolation — through the full routing + middleware + handler + DB stack.
func TestTransactionHistoryIntegration(t *testing.T) {
	s, pool, do := integrationEnv(t)
	ctx := context.Background()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)

	var userA, userB, catA, catB, pA, pB int64
	id, err := s.testInsertUser(ctx, "hist-a-"+suffix+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	userA = id
	id, err = s.testInsertUser(ctx, "hist-b-"+suffix+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	userB = id
	for _, c := range []struct {
		uid  int64
		cid  *int64
		name string
	}{
		{userA, &catA, "HistCatA-" + suffix},
		{userB, &catB, "HistCatB-" + suffix},
	} {
		if err := pool.QueryRow(ctx,
			`INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id`,
			c.uid, c.name).Scan(c.cid); err != nil {
			t.Fatalf("fixture category: %v", err)
		}
	}
	for _, f := range []struct {
		uid  int64
		cid  int64
		name string
		pid  *int64
	}{
		{userA, catA, "HistP-A-" + suffix, &pA},
		{userB, catB, "HistP-B-" + suffix, &pB},
	} {
		if err := pool.QueryRow(ctx,
			`INSERT INTO products (user_id, category_id, name, sku, unit_price, opening_stock, current_stock, min_stock)
			 VALUES ($1, $2, $3, $4, 1, 0, 0, 0) RETURNING id`,
			f.uid, f.cid, f.name, "HIST-"+f.name).Scan(f.pid); err != nil {
			t.Fatalf("fixture product: %v", err)
		}
	}
	t.Cleanup(func() {
		ids := fmt.Sprintf("{%d,%d}", userA, userB)
		_, _ = pool.Exec(ctx, `DELETE FROM stock_transactions WHERE user_id = ANY($1)`, ids)
		_, _ = pool.Exec(ctx, `DELETE FROM products WHERE user_id = ANY($1)`, ids)
		_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE user_id = ANY($1)`, ids)
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = ANY($1)`, ids)
	})

	tokenA, _, err := s.tokens.Issue(userA, "a@example.test")
	if err != nil {
		t.Fatal(err)
	}
	tokenB, _, err := s.tokens.Issue(userB, "b@example.test")
	if err != nil {
		t.Fatal(err)
	}

	// Movements for A's product: IN 5, OUT 2, IN 3 (final balance 6).
	for _, m := range []struct {
		typ string
		qty int
	}{
		{"IN", 5}, {"OUT", 2}, {"IN", 3},
	} {
		rr := do("POST", "/api/v1/transactions",
			fmt.Sprintf(`{"product_id":%d,"type":%q,"quantity":%d}`, pA, m.typ, m.qty), tokenA)
		if rr.Code != http.StatusCreated {
			t.Fatalf("fixture movement %s %d: code = %d (body: %s)", m.typ, m.qty, rr.Code, rr.Body.String())
		}
	}
	// One movement for B's product.
	rr := do("POST", "/api/v1/transactions",
		fmt.Sprintf(`{"product_id":%d,"type":"IN","quantity":4}`, pB), tokenB)
	if rr.Code != http.StatusCreated {
		t.Fatalf("fixture movement B: code = %d", rr.Code)
	}

	// -- TC-42: filter by product --------------------------------------------
	rr = do("GET", fmt.Sprintf("/api/v1/transactions?product_id=%d", pA), "", tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("history by product: code = %d", rr.Code)
	}
	var out struct {
		Transactions []transaction `json:"transactions"`
	}
	decodeTestJSON(t, rr, &out)
	if len(out.Transactions) != 3 {
		t.Fatalf("history by product: got %d rows, want 3", len(out.Transactions))
	}
	for _, t2 := range out.Transactions {
		if t2.ProductID != pA {
			t.Fatalf("history contains a foreign product row: %+v", t2)
		}
	}

	// -- TC-43: filter by type -------------------------------------------------
	rr = do("GET", "/api/v1/transactions?type=OUT", "", tokenA)
	decodeTestJSON(t, rr, &out)
	if rr.Code != http.StatusOK || len(out.Transactions) != 1 || out.Transactions[0].Type != "OUT" || out.Transactions[0].Quantity != 2 {
		t.Fatalf("history by type OUT: code=%d rows=%+v", rr.Code, out.Transactions)
	}

	// -- combined product + type filter ----------------------------------------
	rr = do("GET", fmt.Sprintf("/api/v1/transactions?product_id=%d&type=IN", pA), "", tokenA)
	decodeTestJSON(t, rr, &out)
	if len(out.Transactions) != 2 {
		t.Fatalf("combined filter: got %d IN rows for the product, want 2", len(out.Transactions))
	}

	// -- newest first (created_at DESC) ----------------------------------------
	rr = do("GET", fmt.Sprintf("/api/v1/transactions?product_id=%d", pA), "", tokenA)
	decodeTestJSON(t, rr, &out)
	if out.Transactions[0].Quantity != 3 || out.Transactions[2].Quantity != 5 {
		t.Fatalf("ordering wrong: %+v", out.Transactions)
	}

	// -- AC-02: B's history contains only B's movement --------------------------
	rr = do("GET", "/api/v1/transactions", "", tokenB)
	decodeTestJSON(t, rr, &out)
	if len(out.Transactions) != 1 || out.Transactions[0].ProductID != pB {
		t.Fatalf("isolation: B saw %+v", out.Transactions)
	}
	rr = do("GET", fmt.Sprintf("/api/v1/transactions?product_id=%d", pA), "", tokenB)
	decodeTestJSON(t, rr, &out)
	if len(out.Transactions) != 0 {
		t.Fatalf("isolation: B queried A's product history and got %d rows, want 0", len(out.Transactions))
	}

	// -- limit (Q8) --------------------------------------------------------------
	rr = do("GET", fmt.Sprintf("/api/v1/transactions?product_id=%d&limit=2", pA), "", tokenA)
	decodeTestJSON(t, rr, &out)
	if len(out.Transactions) != 2 {
		t.Fatalf("limit=2: got %d rows, want 2", len(out.Transactions))
	}

	// -- invalid filters rejected -------------------------------------------------
	for _, q := range []string{"?type=BAD", "?product_id=abc", "?limit=0", "?limit=9999"} {
		rr := do("GET", "/api/v1/transactions"+q, "", tokenA)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("GET %s: code = %d, want 400", q, rr.Code)
		}
	}
}
