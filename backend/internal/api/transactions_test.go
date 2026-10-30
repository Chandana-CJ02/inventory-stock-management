package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"inventory-backend/internal/auth"
)

// T-203 registers POST /api/v1/transactions; T-204 adds the history GET.
// Both require auth.
func TestMovementRoutesRequireAuth(t *testing.T) {
	mux := newNilPoolServer(t)
	for _, tc := range []struct{ method, path string }{
		{http.MethodPost, "/api/v1/transactions"},
		{http.MethodGet, "/api/v1/transactions"},
		{http.MethodGet, "/api/v1/transactions?product_id=1&type=OUT"},
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

// TC-44: stock movements are immutable — no PUT/DELETE routes exist, so any
// mutation attempt yields 404 (no such route) and can never touch history.
func TestMovementMutationMethodsHaveNoRoutes(t *testing.T) {
	mux := newNilPoolServer(t)
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, _ := tokens.Issue(1, "user@example.test")
	for _, method := range []string{http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/v1/transactions/1", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusNotFound {
				t.Fatalf("code = %d, want 404 (no mutation route may exist)", rr.Code)
			}
		})
	}
}

func TestCreateMovementValidationFailsClosed(t *testing.T) {
	mux := newNilPoolServer(t) // nil pool: validation must reject before any DB access
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, _ := tokens.Issue(1, "user@example.test")

	for _, tc := range []struct {
		name string
		body string
	}{
		{"missing product", `{"type":"IN","quantity":1}`},
		{"zero product", `{"product_id":0,"type":"IN","quantity":1}`},
		{"missing type", `{"product_id":1,"quantity":1}`},
		{"bad type", `{"product_id":1,"type":"SIDEWAYS","quantity":1}`},
		{"missing quantity", `{"product_id":1,"type":"IN"}`},
		{"zero quantity", `{"product_id":1,"type":"IN","quantity":0}`},
		{"negative quantity", `{"product_id":1,"type":"IN","quantity":-3}`},
		{"fractional quantity", `{"product_id":1,"type":"OUT","quantity":2.5}`},
		{"remarks too long", `{"product_id":1,"type":"IN","quantity":1,"remarks":"` + strings.Repeat("x", 501) + `"}`},
		{"malformed json", `{bad`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", strings.NewReader(tc.body))
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
