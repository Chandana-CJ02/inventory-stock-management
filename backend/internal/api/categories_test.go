package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"inventory-backend/internal/auth"
)

// newNilPoolServer returns a Server whose pool is nil. Only code paths that
// never touch the database (auth rejection, input validation) may be exercised
// with it — a pass-through to the DB would fail the test.
func newNilPoolServer(t *testing.T) *http.ServeMux {
	t.Helper()
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	return New(nil, tokens)
}

func TestCategoryRoutesRequireAuth(t *testing.T) {
	mux := newNilPoolServer(t)
	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/api/v1/categories"},
		{http.MethodPost, "/api/v1/categories"},
		{http.MethodGet, "/api/v1/categories/1"},
		{http.MethodPut, "/api/v1/categories/1"},
		{http.MethodDelete, "/api/v1/categories/1"},
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

func TestCategoryIDValidation(t *testing.T) {
	mux := newNilPoolServer(t)
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, err := tokens.Issue(1, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		"/api/v1/categories/abc",
		"/api/v1/categories/0",
		"/api/v1/categories/-3",
	} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("code = %d, want 400 (body: %s)", rr.Code, rr.Body.String())
			}
		})
	}
}

func TestCategoryNameValidation(t *testing.T) {
	mux := newNilPoolServer(t)
	tokens := auth.NewTokenManager(testSecret, time.Hour)
	token, _, err := tokens.Issue(1, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		body string
	}{
		{"empty name", `{"name":""}`},
		{"whitespace name", `{"name":"   "}`},
		{"missing name", `{}`},
		{"too long name", `{"name":"` + strings.Repeat("x", 101) + `"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", strings.NewReader(tc.body))
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
