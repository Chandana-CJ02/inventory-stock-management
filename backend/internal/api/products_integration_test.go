package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// TestProductCRUDIntegration covers TC-21…TC-29, TC-34/TC-35 and the
// cross-user isolation pattern against the real database. Fixtures are
// throwaway rows with unique suffixes; cleanup removes exactly those rows.
func TestProductCRUDIntegration(t *testing.T) {
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
		{&userA, "prod-a-" + suffix + "@example.test", &catA, "CatA-" + suffix},
		{&userB, "prod-b-" + suffix + "@example.test", &catB, "CatB-" + suffix},
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

	type pr struct {
		ID           int64   `json:"id"`
		CategoryID   int64   `json:"category_id"`
		Name         string  `json:"name"`
		SKU          string  `json:"sku"`
		UnitPrice    float64 `json:"unit_price"`
		OpeningStock int     `json:"opening_stock"`
		CurrentStock int     `json:"current_stock"`
		MinStock     int     `json:"min_stock"`
		Status       string  `json:"status"`
	}

	createBody := func(name, sku string, catID int64) string {
		return fmt.Sprintf(`{"name":%q,"sku":%q,"category_id":%d,"unit_price":10.5,"opening_stock":10,"min_stock":3}`, name, sku, catID)
	}

	// -- TC-21: create valid product; current stock equals opening stock -----
	skuA := "SKU-A-" + suffix
	rr := do("POST", "/api/v1/products", createBody("Widget-"+suffix, skuA, catA), tokenA)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: code = %d, want 201 (body: %s)", rr.Code, rr.Body.String())
	}
	var p pr
	decodeTestJSON(t, rr, &p)
	if p.CurrentStock != 10 || p.OpeningStock != 10 || p.Status != "active" || p.UnitPrice != 10.5 {
		t.Fatalf("created product unexpected: %+v", p)
	}
	productPath := fmt.Sprintf("/api/v1/products/%d", p.ID)

	// -- AC-07 / TC-27: category must exist and belong to the user -----------
	rr = do("POST", "/api/v1/products", createBody("NoCat-"+suffix, "SKUNC-"+suffix, 999999999), tokenA)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("nonexistent category: code = %d, want 400", rr.Code)
	}
	rr = do("POST", "/api/v1/products", createBody("ForeignCat-"+suffix, "SKUFC-"+suffix, catB), tokenA)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("foreign category: code = %d, want 400", rr.Code)
	}

	// -- TC-22/23/24: field validation ---------------------------------------
	for _, tc := range []struct {
		name string
		body string
	}{
		{"negative price", `{"name":"P","sku":"X1","category_id":` + fmt.Sprint(catA) + `,"unit_price":-1,"opening_stock":0,"min_stock":0}`},
		{"fractional opening stock", `{"name":"P","sku":"X2","category_id":` + fmt.Sprint(catA) + `,"unit_price":1,"opening_stock":2.5,"min_stock":0}`},
		{"negative opening stock", `{"name":"P","sku":"X3","category_id":` + fmt.Sprint(catA) + `,"unit_price":1,"opening_stock":-1,"min_stock":0}`},
		{"negative min stock", `{"name":"P","sku":"X4","category_id":` + fmt.Sprint(catA) + `,"unit_price":1,"opening_stock":0,"min_stock":-2}`},
		{"missing name", `{"sku":"X5","category_id":` + fmt.Sprint(catA) + `,"unit_price":1,"opening_stock":0,"min_stock":0}`},
		{"missing sku", `{"name":"P","category_id":` + fmt.Sprint(catA) + `,"unit_price":1,"opening_stock":0,"min_stock":0}`},
	} {
		t.Run("validation/"+tc.name, func(t *testing.T) {
			rr := do("POST", "/api/v1/products", tc.body, tokenA)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("code = %d, want 400 (body: %s)", rr.Code, rr.Body.String())
			}
		})
	}

	// -- TC-25: duplicate SKU within the user's inventory --------------------
	rr = do("POST", "/api/v1/products", createBody("Dup-"+suffix, skuA, catA), tokenA)
	if rr.Code != http.StatusConflict {
		t.Fatalf("duplicate sku: code = %d, want 409", rr.Code)
	}

	// -- TC-26: same SKU allowed for a different user -------------------------
	rr = do("POST", "/api/v1/products", createBody("Other-"+suffix, skuA, catB), tokenB)
	if rr.Code != http.StatusCreated {
		t.Fatalf("same sku other user: code = %d, want 201", rr.Code)
	}

	// -- TC-10/TC-11 pattern: cross-user get/update/deactivate are 404 --------
	for _, tc := range []struct{ method, path, body string }{
		{"GET", productPath, ""},
		{"PUT", productPath, `{"name":"Hijacked"}`},
		{"POST", productPath + "/deactivate", ""},
	} {
		rr := do(tc.method, tc.path, tc.body, tokenB)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("cross-user %s: code = %d, want 404", tc.method, rr.Code)
		}
	}

	// -- TC-29: update persists; current_stock is never touched by edits -----
	rr = do("PUT", productPath,
		fmt.Sprintf(`{"name":"Widget2-%s","unit_price":12,"min_stock":5,"category_id":%d}`, suffix, catA), tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("update: code = %d, want 200 (body: %s)", rr.Code, rr.Body.String())
	}
	rr = do("GET", productPath, "", tokenA)
	decodeTestJSON(t, rr, &p)
	if p.Name != "Widget2-"+suffix || p.UnitPrice != 12 || p.MinStock != 5 {
		t.Fatalf("update did not persist: %+v", p)
	}
	if p.CurrentStock != 10 {
		t.Fatalf("current_stock changed by an edit: %d, want 10 (edits never move stock)", p.CurrentStock)
	}

	// -- AC-14 / TC-34: deactivate (idempotent) -------------------------------
	rr = do("POST", productPath+"/deactivate", "", tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("deactivate: code = %d, want 200", rr.Code)
	}
	rr = do("GET", productPath, "", tokenA)
	decodeTestJSON(t, rr, &p)
	if p.Status != "inactive" {
		t.Fatalf("status = %q, want inactive", p.Status)
	}
	rr = do("POST", productPath+"/deactivate", "", tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("deactivate (idempotent): code = %d, want 200", rr.Code)
	}

	// -- Q14: reactivation via edit -------------------------------------------
	rr = do("PUT", productPath, `{"status":"active"}`, tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("reactivate: code = %d, want 200", rr.Code)
	}
	rr = do("GET", productPath, "", tokenA)
	decodeTestJSON(t, rr, &p)
	if p.Status != "active" {
		t.Fatalf("status = %q, want active after reactivation", p.Status)
	}

	// -- TC-35: no delete path for products -----------------------------------
	rr = do("DELETE", productPath, "", tokenA)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("delete product: code = %d, want 405 (products are never hard-deleted)", rr.Code)
	}
}
