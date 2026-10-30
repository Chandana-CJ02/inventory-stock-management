package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// TestDashboardIntegration covers TC-49…TC-51 (FR-5.1…FR-5.3): exact totals,
// the low-stock list (active products at or below minimum, boundary included),
// the five most recent transactions, empty-dashboard edge case, and per-user
// isolation. Fixtures are self-cleaning throwaway rows.
func TestDashboardIntegration(t *testing.T) {
	s, pool, do := integrationEnv(t)
	ctx := context.Background()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)

	var userA, userB int64
	id, err := s.testInsertUser(ctx, "dash-a-"+suffix+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	userA = id
	id, err = s.testInsertUser(ctx, "dash-b-"+suffix+"@example.test")
	if err != nil {
		t.Fatal(err)
	}
	userB = id
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

	type summary struct {
		Totals struct {
			ActiveProducts int `json:"active_products"`
			Categories     int `json:"categories"`
			LowStock       int `json:"low_stock"`
		} `json:"totals"`
		LowStock []struct {
			ID           int64 `json:"id"`
			CurrentStock int   `json:"current_stock"`
			MinStock     int   `json:"min_stock"`
		} `json:"low_stock"`
		Recent []struct {
			ID        int64 `json:"id"`
			ProductID int64 `json:"product_id"`
			Quantity  int   `json:"quantity"`
		} `json:"recent_transactions"`
	}
	summaryFor := func(token string) summary {
		rr := do("GET", "/api/v1/dashboard/summary", "", token)
		if rr.Code != http.StatusOK {
			t.Fatalf("summary: code = %d, want 200 (body: %s)", rr.Code, rr.Body.String())
		}
		var out summary
		decodeTestJSON(t, rr, &out)
		return out
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

	// -- Empty dashboard edge case (fresh user B, before any data) ------------
	sb := summaryFor(tokenB)
	if sb.Totals.ActiveProducts != 0 || sb.Totals.Categories != 0 || sb.Totals.LowStock != 0 {
		t.Fatalf("empty dashboard totals = %+v, want zeros", sb.Totals)
	}
	if len(sb.LowStock) != 0 || len(sb.Recent) != 0 {
		t.Fatalf("empty dashboard lists = %+v / %+v, want empty", sb.LowStock, sb.Recent)
	}

	// -- Fixtures for user A ---------------------------------------------------
	var catA1, catA2 int64
	for _, c := range []struct {
		cid  *int64
		name string
	}{
		{&catA1, "DashCat1-" + suffix},
		{&catA2, "DashCat2-" + suffix},
	} {
		if err := pool.QueryRow(ctx,
			`INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id`,
			userA, c.name).Scan(c.cid); err != nil {
			t.Fatalf("fixture category: %v", err)
		}
	}
	mkProduct := func(name, sku string, cat int64, opening, min int64) int64 {
		rr := do("POST", "/api/v1/products",
			fmt.Sprintf(`{"name":%q,"sku":%q,"category_id":%d,"unit_price":2,"opening_stock":%d,"min_stock":%d}`,
				name, sku, cat, opening, min), tokenA)
		if rr.Code != http.StatusCreated {
			t.Fatalf("fixture product %s: code = %d (body: %s)", name, rr.Code, rr.Body.String())
		}
		var created struct {
			ID int64 `json:"id"`
		}
		decodeTestJSON(t, rr, &created)
		return created.ID
	}
	p1 := mkProduct("Dash-OK-"+suffix, "DOK-"+suffix, catA1, 10, 3)         // healthy
	p2 := mkProduct("Dash-LOW-"+suffix, "DLOW-"+suffix, catA1, 2, 5)        // low: 2 <= 5
	p4 := mkProduct("Dash-BOUNDARY-"+suffix, "DBOUND-"+suffix, catA2, 5, 5) // low boundary: 5 <= 5
	p3 := mkProduct("Dash-INACTIVE-"+suffix, "DIN-"+suffix, catA1, 0, 0)    // will be deactivated

	rr := do("POST", fmt.Sprintf("/api/v1/products/%d/deactivate", p3), "", tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("deactivate P3: code = %d", rr.Code)
	}

	// -- TC-49: totals (active=3: P1,P2,P4; categories=2; low=2: P2,P4) -------
	sa := summaryFor(tokenA)
	if sa.Totals.ActiveProducts != 3 || sa.Totals.Categories != 2 || sa.Totals.LowStock != 2 {
		t.Fatalf("TC-49 totals = %+v, want active=3 categories=2 low=2", sa.Totals)
	}
	lowIDs := map[int64]bool{}
	for _, l := range sa.LowStock {
		lowIDs[l.ID] = true
	}
	if len(sa.LowStock) != 2 || !lowIDs[p2] || !lowIDs[p4] || lowIDs[p1] || lowIDs[p3] {
		t.Fatalf("TC-50 low-stock list = %+v, want exactly {P2, P4} (P4 boundary included, P3 excluded as inactive)", sa.LowStock)
	}

	// -- TC-51: five most recent transactions ---------------------------------
	// Seven movements on P1: OUT2, IN3, OUT4, OUT1, IN2, OUT4, IN1 (final 5).
	for _, m := range []struct {
		typ string
		qty int
	}{
		{"OUT", 2}, {"IN", 3}, {"OUT", 4}, {"OUT", 1}, {"IN", 2}, {"OUT", 4}, {"IN", 1},
	} {
		rr := do("POST", "/api/v1/transactions",
			fmt.Sprintf(`{"product_id":%d,"type":%q,"quantity":%d}`, p1, m.typ, m.qty), tokenA)
		if rr.Code != http.StatusCreated {
			t.Fatalf("fixture movement %s %d: code = %d", m.typ, m.qty, rr.Code)
		}
	}
	if got := stockOf(p1); got != 5 {
		t.Fatalf("P1 balance after fixtures = %d, want 5", got)
	}

	// Collect the 7 movement IDs in insertion order (m1..m7).
	rows, err := pool.Query(ctx,
		`SELECT id FROM stock_transactions WHERE product_id = $1 ORDER BY id`, p1)
	if err != nil {
		t.Fatal(err)
	}
	var allIDs []int64
	for rows.Next() {
		var mid int64
		if err := rows.Scan(&mid); err != nil {
			t.Fatal(err)
		}
		allIDs = append(allIDs, mid)
	}
	rows.Close()
	if len(allIDs) != 7 {
		t.Fatalf("expected 7 movements on P1, got %d", len(allIDs))
	}
	want := []int64{allIDs[6], allIDs[5], allIDs[4], allIDs[3], allIDs[2]} // newest 5, DESC

	sa = summaryFor(tokenA)
	if len(sa.Recent) != 5 {
		t.Fatalf("TC-51 recent = %d rows, want exactly 5", len(sa.Recent))
	}
	for i, wantID := range want {
		if sa.Recent[i].ID != wantID {
			t.Fatalf("TC-51 recent[%d].id = %d, want %d (newest first)", i, sa.Recent[i].ID, wantID)
		}
		if sa.Recent[i].ProductID != p1 {
			t.Fatalf("TC-51 recent[%d] belongs to product %d, want %d", i, sa.Recent[i].ProductID, p1)
		}
	}

	// -- AC-02: B's dashboard still shows only B's own data -------------------
	sb = summaryFor(tokenB)
	if sb.Totals.ActiveProducts != 0 || sb.Totals.Categories != 0 || sb.Totals.LowStock != 0 || len(sb.Recent) != 0 {
		t.Fatalf("AC-02: B's dashboard leaked A's data: %+v", sb)
	}

	// -- B gets data; dashboard reflects it independently ---------------------
	rr = do("POST", "/api/v1/categories", fmt.Sprintf(`{"name":"DashB-%s"}`, suffix), tokenB)
	if rr.Code != http.StatusCreated {
		t.Fatalf("B category: code = %d", rr.Code)
	}
	var catB struct {
		ID int64 `json:"id"`
	}
	decodeTestJSON(t, rr, &catB)
	rr = do("POST", "/api/v1/products",
		fmt.Sprintf(`{"name":"DashBProd-%s","sku":"DB-%s","category_id":%d,"unit_price":1,"opening_stock":4,"min_stock":6}`, suffix, suffix, catB.ID), tokenB)
	if rr.Code != http.StatusCreated {
		t.Fatalf("B product: code = %d", rr.Code)
	}
	sb = summaryFor(tokenB)
	if sb.Totals.ActiveProducts != 1 || sb.Totals.Categories != 1 || sb.Totals.LowStock != 1 || len(sb.LowStock) != 1 {
		t.Fatalf("B dashboard after data = %+v / low=%+v, want 1/1/1", sb.Totals, sb.LowStock)
	}
}
