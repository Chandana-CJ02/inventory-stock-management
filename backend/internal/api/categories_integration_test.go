package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"inventory-backend/internal/auth"
	"inventory-backend/internal/config"
	"inventory-backend/internal/database"
)

// integrationEnv skips unless DB integration tests are requested, then wires a
// real Server (routing + middleware + handlers + pool) against inventory_db.
// findEnvFile walks up from the test's working directory (the package dir)
// until it finds the repository-root .env file.
func findEnvFile(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Skipf("working directory unavailable: %v", err)
	}
	for i := 0; i < 8; i++ {
		p := filepath.Join(dir, ".env")
		if _, err := os.Stat(p); err == nil {
			return p
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Skip(".env not found above the test working directory")
	return ""
}

func integrationEnv(t *testing.T) (*Server, *pgxpool.Pool, func(method, path, body, token string) *httptest.ResponseRecorder) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run category integration tests against inventory_db")
	}
	envPath := findEnvFile(t)
	cfg, err := config.Load(envPath)
	if err != nil {
		t.Skipf("config unavailable: %v", err)
	}
	pool, err := database.NewPool(context.Background(), cfg)
	if err != nil {
		t.Fatalf("integration pool: %v", err)
	}
	t.Cleanup(pool.Close)

	tokens := auth.NewTokenManager(cfg.JWTSecret, time.Hour)
	server := &Server{pool: pool, tokens: tokens, start: time.Now()}
	mux := New(pool, tokens)
	do := func(method, path, body, token string) *httptest.ResponseRecorder {
		var rd io.Reader
		if body != "" {
			rd = strings.NewReader(body)
		}
		req := httptest.NewRequest(method, path, rd)
		req.Header.Set("Content-Type", "application/json")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}
	return server, pool, do
}

func decodeTestJSON(t *testing.T, rr *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rr.Body.Bytes(), v); err != nil {
		t.Fatalf("invalid JSON response: %v — body: %s", err, rr.Body.String())
	}
}

// testInsertUser creates a throwaway user row directly (schema-level fixture;
// the hash value never authenticates anything). Same package: s.pool is used.
func (s *Server) testInsertUser(ctx context.Context, email string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, 'integration-test-not-a-login-hash') RETURNING id`,
		email).Scan(&id)
	return id, err
}

// TestCategoryCRUDIntegration covers TC-13…TC-18 and the cross-user isolation
// cases TC-10/TC-11 against the real database. It creates its own throwaway
// users/categories/products with unique suffixes and removes exactly those
// rows afterwards; pre-existing data is never touched.
func TestCategoryCRUDIntegration(t *testing.T) {
	s, pool, do := integrationEnv(t)
	ctx := context.Background()
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)

	// -- fixtures: two throwaway users ------------------------------------
	var userA, userB int64
	for _, u := range []struct {
		id    *int64
		email string
	}{
		{&userA, "cat-a-" + suffix + "@example.test"},
		{&userB, "cat-b-" + suffix + "@example.test"},
	} {
		id, err := s.testInsertUser(ctx, u.email)
		if err != nil {
			t.Fatalf("fixture user %s: %v", u.email, err)
		}
		*u.id = id
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

	// -- TC-13: create + list ---------------------------------------------
	nameA := "Electronics-" + suffix
	rr := do("POST", "/api/v1/categories", `{"name":"`+nameA+`"}`, tokenA)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: code = %d, want 201 (body: %s)", rr.Code, rr.Body.String())
	}
	var created category
	decodeTestJSON(t, rr, &created)
	if created.Name != nameA || created.ID <= 0 {
		t.Fatalf("created = %+v", created)
	}

	rr = do("GET", "/api/v1/categories", "", tokenA)
	var list struct {
		Categories []category `json:"categories"`
	}
	decodeTestJSON(t, rr, &list)
	if rr.Code != http.StatusOK || len(list.Categories) != 1 || list.Categories[0].Name != nameA {
		t.Fatalf("list: code=%d categories=%+v", rr.Code, list.Categories)
	}

	// -- TC-14: duplicate name within the same user's inventory ------------
	rr = do("POST", "/api/v1/categories", `{"name":"`+nameA+`"}`, tokenA)
	if rr.Code != http.StatusConflict {
		t.Fatalf("duplicate: code = %d, want 409", rr.Code)
	}

	// -- TC-15: same name allowed for a different user ---------------------
	rr = do("POST", "/api/v1/categories", `{"name":"`+nameA+`"}`, tokenB)
	if rr.Code != http.StatusCreated {
		t.Fatalf("same name other user: code = %d, want 201", rr.Code)
	}

	// -- TC-10: cross-user read is 404 --------------------------------------
	rr = do("GET", fmt.Sprintf("/api/v1/categories/%d", created.ID), "", tokenB)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("cross-user get: code = %d, want 404", rr.Code)
	}

	// -- TC-11: cross-user update and delete are 404 ------------------------
	rr = do("PUT", fmt.Sprintf("/api/v1/categories/%d", created.ID), `{"name":"Hijacked"}`, tokenB)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("cross-user put: code = %d, want 404", rr.Code)
	}
	rr = do("DELETE", fmt.Sprintf("/api/v1/categories/%d", created.ID), "", tokenB)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("cross-user delete: code = %d, want 404", rr.Code)
	}

	// -- TC-16: update persists ---------------------------------------------
	newName := "Audio-" + suffix
	rr = do("PUT", fmt.Sprintf("/api/v1/categories/%d", created.ID), `{"name":"`+newName+`"}`, tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("update: code = %d, want 200 (body: %s)", rr.Code, rr.Body.String())
	}
	rr = do("GET", fmt.Sprintf("/api/v1/categories/%d", created.ID), "", tokenA)
	var fetched category
	decodeTestJSON(t, rr, &fetched)
	if fetched.Name != newName {
		t.Fatalf("update did not persist: name = %q", fetched.Name)
	}

	// -- TC-14 variant: rename onto the user's own other category name -------
	secondName := "Stationery-" + suffix
	rr = do("POST", "/api/v1/categories", `{"name":"`+secondName+`"}`, tokenA)
	if rr.Code != http.StatusCreated {
		t.Fatalf("second create: code = %d, want 201", rr.Code)
	}
	rr = do("PUT", fmt.Sprintf("/api/v1/categories/%d", created.ID), `{"name":"`+secondName+`"}`, tokenA)
	if rr.Code != http.StatusConflict {
		t.Fatalf("rename onto sibling name: code = %d, want 409", rr.Code)
	}

	// -- validation through the full stack ----------------------------------
	rr = do("POST", "/api/v1/categories", `{"name":"   "}`, tokenA)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("blank name: code = %d, want 400", rr.Code)
	}
	rr = do("GET", "/api/v1/categories/abc", "", tokenA)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("bad id: code = %d, want 400", rr.Code)
	}
	rr = do("GET", "/api/v1/categories/999999999", "", tokenA)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("missing id: code = %d, want 404", rr.Code)
	}

	// -- TC-18 / AC-08: deletion blocked while a product references it -------
	var productID int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO products (user_id, category_id, name, sku, unit_price, opening_stock, current_stock, min_stock)
		 VALUES ($1, $2, 'GuardProduct', $3, 1.00, 0, 0, 0) RETURNING id`,
		userA, created.ID, "GUARD-"+suffix).Scan(&productID); err != nil {
		t.Fatalf("guard product fixture: %v", err)
	}
	rr = do("DELETE", fmt.Sprintf("/api/v1/categories/%d", created.ID), "", tokenA)
	if rr.Code != http.StatusConflict {
		t.Fatalf("delete with product: code = %d, want 409", rr.Code)
	}
	if _, err := pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, productID); err != nil {
		t.Fatalf("guard product cleanup: %v", err)
	}

	// -- TC-17: deletion succeeds once unused --------------------------------
	rr = do("DELETE", fmt.Sprintf("/api/v1/categories/%d", created.ID), "", tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("delete: code = %d, want 200", rr.Code)
	}
	rr = do("GET", fmt.Sprintf("/api/v1/categories/%d", created.ID), "", tokenA)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("get after delete: code = %d, want 404", rr.Code)
	}
}
