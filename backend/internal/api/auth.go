package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"inventory-backend/internal/auth"
)

// maxBodyBytes bounds request bodies defensively.
const maxBodyBytes = 8192

// emailPattern is intentionally conservative (local@domain.tld). The unique
// index is case-insensitive, so emails are normalized to lower case before
// storage and lookup (docs/requirements.md Q4).
var emailPattern = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handleRegister creates a new account. Passwords are stored only as bcrypt
// hashes; responses never include hash material (AC-01, TC-04).
func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if !emailPattern.MatchString(email) || len(email) > 254 {
		writeError(w, http.StatusBadRequest, "a valid email address is required")
		return
	}
	// 8 characters is a documented minimum (unspecified by the assignment);
	// 72 is bcrypt's input limit — longer input is silently truncated by
	// bcrypt, so it is rejected instead.
	if len(req.Password) < 8 || len(req.Password) > 72 {
		writeError(w, http.StatusBadRequest, "password must be between 8 and 72 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("register: hash: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	var id int64
	var createdAt time.Time
	err = s.pool.QueryRow(r.Context(),
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`,
		email, hash).Scan(&id, &createdAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation (users_email_lower_unique)
			writeError(w, http.StatusConflict, "email is already registered")
			return
		}
		log.Printf("register: insert: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         id,
		"email":      email,
		"created_at": createdAt,
	})
}

// handleLogin verifies credentials and returns a signed JWT. Unknown accounts
// and wrong passwords produce the same 401 message, with bcrypt timing
// equalized for unknown accounts, so the endpoint cannot be used to enumerate
// registered emails.
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req credentialsRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	var (
		id          int64
		storedEmail string
		hash        string
	)
	err := s.pool.QueryRow(r.Context(),
		`SELECT id, email, password_hash FROM users WHERE lower(email) = $1`,
		email).Scan(&id, &storedEmail, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		auth.VerifyDummy(req.Password) // timing equalizer
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		log.Printf("login: select: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}
	if !auth.VerifyPassword(hash, req.Password) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, expiresAt, err := s.tokens.Issue(id, storedEmail)
	if err != nil {
		log.Printf("login: issue token: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token":      token,
		"token_type": "Bearer",
		"expires_at": expiresAt,
	})
}

// handleLogout ends the client session. The API issues stateless JWTs without
// server-side revocation (approved architecture, Q5): logout means the client
// discards the token, and the short expiry (24 h) bounds the exposure of a
// leaked one. This limitation is documented rather than hidden — no "logged
// out" token is actually invalid until it expires.
func (s *Server) handleLogout(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "logged out; discard the token client-side (stateless JWTs are not server-revoked)",
	})
}
