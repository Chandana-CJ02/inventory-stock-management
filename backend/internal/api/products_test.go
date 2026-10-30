package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"inventory-backend/internal/auth"
)

func TestProductRoutesRequireAuth(t *testing.T) {
	mux := newNilPoolServer(t)
	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/api/v1/products"},
		{http.MethodPost, "/api/v1/products"},
		{http.MethodGet, "/api/v1/products/1"},
		{http.MethodPut, "/api/v1/products/1"},
		{http.MethodPost, "/api/v1/products/1/deactivate"},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusUnauthorized {
				t.Fatalf("code = %d, want 401", rr.Code)
			}
		})
	}
}

func TestProductIDValidation(t *testing.T) {
	mux := newNilPoolServer(t)
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, err := tokens.Issue(1, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"/api/v1/products/abc",
		"/api/v1/products/0",
	} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("code = %d, want 400", rr.Code)
			}
		})
	}
}

// No DELETE route exists for products (TC-35: products with history are never
// hard-deleted). The mux must answer 405 for the known path.
func TestProductDeleteMethodNotAllowed(t *testing.T) {
	mux := newNilPoolServer(t)
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, _ := tokens.Issue(1, "user@example.test")
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("code = %d, want 405 (products must have no delete path)", rr.Code)
	}
}

func TestCreateProductValidationFailsClosed(t *testing.T) {
	mux := newNilPoolServer(t) // nil pool: validation must reject before any DB access
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, _ := tokens.Issue(1, "user@example.test")

	for _, tc := range []struct {
		name string
		body string
	}{
		{"missing name", `{"sku":"S1","category_id":1,"unit_price":1,"opening_stock":0,"min_stock":0}`},
		{"blank name", `{"name":"   ","sku":"S1","category_id":1,"unit_price":1,"opening_stock":0,"min_stock":0}`},
		{"missing sku", `{"name":"P","category_id":1,"unit_price":1,"opening_stock":0,"min_stock":0}`},
		{"missing category", `{"name":"P","sku":"S1","unit_price":1,"opening_stock":0,"min_stock":0}`},
		{"missing price", `{"name":"P","sku":"S1","category_id":1,"opening_stock":0,"min_stock":0}`},
		{"negative price", `{"name":"P","sku":"S1","category_id":1,"unit_price":-0.01,"opening_stock":0,"min_stock":0}`},
		{"negative opening stock", `{"name":"P","sku":"S1","category_id":1,"unit_price":1,"opening_stock":-1,"min_stock":0}`},
		{"fractional opening stock", `{"name":"P","sku":"S1","category_id":1,"unit_price":1,"opening_stock":2.5,"min_stock":0}`},
		{"negative min stock", `{"name":"P","sku":"S1","category_id":1,"unit_price":1,"opening_stock":0,"min_stock":-3}`},
		{"non-integer min stock", `{"name":"P","sku":"S1","category_id":1,"unit_price":1,"opening_stock":0,"min_stock":0.5}`},
		{"malformed json", `{bad`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("code = %d, want 400 (body: %s)", rr.Code, rr.Body.String())
			}
		})
	}
}
