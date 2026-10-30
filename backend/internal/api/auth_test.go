package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"inventory-backend/internal/auth"
)

const testSecret = "unit-test-secret-0123456789abcdef0123456789abcdef"

func newTestServer(t *testing.T) *Server {
	t.Helper()
	// pool is nil: these tests only exercise paths that do not touch the
	// database (validation failures, middleware, logout).
	return &Server{
		tokens: auth.NewTokenManager(testSecret, time.Hour),
		start:  time.Now(),
	}
}

func TestRequireAuth(t *testing.T) {
	s := newTestServer(t)
	called := false
	protected := s.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		id, ok := UserID(r.Context())
		if !ok || id != 42 {
			t.Errorf("context user id = %d, %v; want 42, true", id, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))

	post := func(authz string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		if authz != "" {
			req.Header.Set("Authorization", authz)
		}
		rr := httptest.NewRecorder()
		protected.ServeHTTP(rr, req)
		return rr
	}

	if rr := post(""); rr.Code != http.StatusUnauthorized || called {
		t.Fatalf("no header: code=%d called=%v, want 401 false", rr.Code, called)
	}
	if rr := post("Token abc"); rr.Code != http.StatusUnauthorized || called {
		t.Fatalf("malformed scheme: code=%d called=%v, want 401 false", rr.Code, called)
	}
	if rr := post("Bearer not.a.jwt"); rr.Code != http.StatusUnauthorized || called {
		t.Fatalf("garbage token: code=%d called=%v, want 401 false", rr.Code, called)
	}

	other := auth.NewTokenManager("another-secret-0123456789abcdef0123456789abcdef", time.Hour)
	wrongSig, _, err := other.Issue(42, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}
	if rr := post("Bearer " + wrongSig); rr.Code != http.StatusUnauthorized || called {
		t.Fatalf("wrong-secret token: code=%d called=%v, want 401 false", rr.Code, called)
	}

	token, _, err := s.tokens.Issue(42, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}
	rr := post("Bearer " + token)
	if rr.Code != http.StatusOK || !called {
		t.Fatalf("valid token: code=%d called=%v, want 200 true", rr.Code, called)
	}
}

func TestExpiredTokenRejectedByMiddleware(t *testing.T) {
	expired := auth.NewTokenManager(testSecret, -time.Minute)
	token, _, err := expired.Issue(1, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}

	s := newTestServer(t)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	s.requireAuth(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("handler must not run for an expired token")
	})).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code = %d, want 401", rr.Code)
	}
}

func TestRegisterValidationFailsClosed(t *testing.T) {
	s := newTestServer(t) // nil pool: a pass-through to the DB would panic, failing the test
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/register", s.handleRegister)

	for _, tc := range []struct {
		name string
		body string
		want int
	}{
		{"bad email", `{"email":"not-an-email","password":"GoodPassword123"}`, 400},
		{"empty email", `{"email":"","password":"GoodPassword123"}`, 400},
		{"short password", `{"email":"a@b.co","password":"short"}`, 400},
		{"missing password", `{"email":"a@b.co"}`, 400},
		{"malformed json", `{nope`, 400},
		{"empty body", ``, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != tc.want {
				t.Fatalf("code = %d, want %d (body: %s)", rr.Code, tc.want, rr.Body.String())
			}
			var resp map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("response is not JSON: %v", err)
			}
			if resp["error"] == "" {
				t.Fatal("error response has no message")
			}
		})
	}
}

func TestLoginValidationFailsClosed(t *testing.T) {
	s := newTestServer(t) // nil pool
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth/login", s.handleLogin)

	for _, tc := range []struct {
		name string
		body string
		want int
	}{
		{"missing email", `{"password":"whatever1"}`, 400},
		{"missing password", `{"email":"a@b.co"}`, 400},
		{"malformed json", `{`, 400},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)
			if rr.Code != tc.want {
				t.Fatalf("code = %d, want %d", rr.Code, tc.want)
			}
		})
	}
}
