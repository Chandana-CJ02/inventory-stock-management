package api

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

type contextKey int

const userIDKey contextKey = 0

// UserID returns the authenticated user's ID stored by requireAuth. Protected
// handlers must scope every database query by this ID — never by any
// client-supplied value (AC-02).
func UserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}

// requireAuth rejects requests without a valid Bearer token (401) and stores
// the authenticated user's ID in the request context.
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(authz, prefix) {
			writeError(w, http.StatusUnauthorized, "missing or malformed Authorization header")
			return
		}
		claims, err := s.tokens.Verify(strings.TrimSpace(strings.TrimPrefix(authz, prefix)))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		id, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil || id <= 0 {
			writeError(w, http.StatusUnauthorized, "invalid token subject")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, id)))
	})
}
