package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDashboardRouteRequiresAuth(t *testing.T) {
	mux := newNilPoolServer(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard/summary", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("code = %d, want 401", rr.Code)
	}
}
