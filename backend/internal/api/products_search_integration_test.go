package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// TestProductSearchIntegration covers T-202 / TC-30…TC-32: search by name and
// SKU, category filter, combined filters, and cross-user isolation of search
// results. Fixtures are self-cleaning throwaway rows.
func TestProductSearchIntegration(t *testing.T) {
	s, pool, do := integrationEnv(t)
	ctx := context.Background()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)

	var userA, userB, catA, catB int64
	for _, f := range []struct {
		uid   *int64
		email string
		cid   *int64
		cname string
	}{
		{&userA, "srch-a-" + suffix + "@example.test", &catA, "SrchCatA-" + suffix},
		{&userB, "srch-b-" + suffix + "@example.test", &catB, "SrchCatB-" + suffix},
	} {
		id, err := s.testInsertUser(ctx, f.email)
		if err != nil {
			t.Fatalf("fixture user: %v", err)
		}
		*f.uid = id
		if err := pool.QueryRow(ctx,
			`INSERT INTO categories (user_id, name) VALUES ($1, $2) RETURNING id`,
			id, f.cname).Scan(f.cid); err != nil {
			t.Fatalf("fixture category: %v", err)
		}
	}
	t.Cleanup(func() {
		ids := fmt.Sprintf("{%d,%d}", userA, userB)
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

	// Fixtures (created through the T-201 API): user A owns an Alpha Widget in
	// catA and a Beta Gadget in catA; user B owns its own Alpha Widget in catB.
	mk := func(name, sku string, cat int64, token string) {
		rr := do("POST", "/api/v1/products",
			fmt.Sprintf(`{"name":%q,"sku":%q,"category_id":%d,"unit_price":5,"opening_stock":1,"min_stock":1}`, name, sku, cat),
			token)
		if rr.Code != http.StatusCreated {
			t.Fatalf("fixture %s: code = %d, want 201 (body: %s)", name, rr.Code, rr.Body.String())
		}
	}
	nameA1, skuA1 := "Alpha Widget-"+suffix, "AW-1-"+suffix
	nameA2, skuA2 := "Beta Gadget-"+suffix, "BG-2-"+suffix
	mk(nameA1, skuA1, catA, tokenA)
	mk(nameA2, skuA2, catA, tokenA)
	mk("Alpha Widget-"+suffix, "AW-1-B-"+suffix, catB, tokenB) // same name/sku prefix as A's, different user

	idsOf := func(token string, query string) []int64 {
		rr := do("GET", "/api/v1/products"+query, "", token)
		if rr.Code != http.StatusOK {
			t.Fatalf("GET %s: code = %d, want 200 (body: %s)", query, rr.Code, rr.Body.String())
		}
		var out struct {
			Products []product `json:"products"`
		}
		decodeTestJSON(t, rr, &out)
		ids := make([]int64, 0, len(out.Products))
		for _, p := range out.Products {
			ids = append(ids, p.ID)
		}
		return ids
	}

	t.Run("TC-30 search by name substring is case-insensitive", func(t *testing.T) {
		ids := idsOf(tokenA, "?search=alpha+widget")
		if len(ids) != 1 {
			t.Fatalf("search 'alpha widget': got %d products, want 1 (A's own)", len(ids))
		}
	})

	t.Run("TC-30 search by name fragment", func(t *testing.T) {
		ids := idsOf(tokenA, "?search="+nameA1[:5])
		if len(ids) != 1 {
			t.Fatalf("search fragment: got %d products, want 1", len(ids))
		}
	})

	t.Run("TC-31 search by SKU", func(t *testing.T) {
		ids := idsOf(tokenA, "?search=aw-1") // matches skuA1 case-insensitively, not skuA2
		if len(ids) != 1 {
			t.Fatalf("sku search: got %d products, want 1", len(ids))
		}
	})

	t.Run("TC-31 search by exact SKU", func(t *testing.T) {
		ids := idsOf(tokenA, "?search="+skuA1)
		if len(ids) != 1 {
			t.Fatalf("exact sku search: got %d products, want 1", len(ids))
		}
	})

	t.Run("TC-32 filter by category", func(t *testing.T) {
		ids := idsOf(tokenA, fmt.Sprintf("?category_id=%d", catA))
		if len(ids) != 2 {
			t.Fatalf("category filter: got %d products, want 2 (both of A's)", len(ids))
		}
	})

	t.Run("TC-32 combined search + category filter", func(t *testing.T) {
		ids := idsOf(tokenA, fmt.Sprintf("?search=widget&category_id=%d", catA))
		if len(ids) != 1 {
			t.Fatalf("combined filter: got %d products, want 1", len(ids))
		}
	})

	t.Run("TC-32 filter with no matches", func(t *testing.T) {
		ids := idsOf(tokenA, "?search=zzz-nothing")
		if len(ids) != 0 {
			t.Fatalf("no-match search: got %d products, want 0", len(ids))
		}
	})

	t.Run("AC-02 search never returns other users' products", func(t *testing.T) {
		ids := idsOf(tokenA, "?search=alpha") // B owns an "Alpha Widget-<suffix>" too
		if len(ids) != 1 {
			t.Fatalf("cross-user search isolation: got %d products, want 1 (only A's own)", len(ids))
		}
	})

	t.Run("invalid category_id filter is rejected", func(t *testing.T) {
		rr := do("GET", "/api/v1/products?category_id=abc", "", tokenA)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("code = %d, want 400", rr.Code)
		}
	})
}
