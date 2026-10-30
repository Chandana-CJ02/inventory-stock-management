package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned for any token that fails verification: bad
// signature, wrong issuer, expired, or malformed. It deliberately does not
// distinguish causes to callers.
var ErrInvalidToken = errors.New("invalid or expired token")

// TokenManager issues and verifies HS256-signed JWTs. The secret is loaded
// from configuration and never logged, persisted, or embedded in binaries.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
	issuer string
}

// NewTokenManager returns a manager using secret for HS256 signatures and the
// given token time-to-live (docs/architecture.md Q6: ~24 h).
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl, issuer: "inventory-api"}
}

// Claims are the token payload fields used by the API.
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// Issue creates a signed token for the user and returns it with its expiry
// time. Tokens and the secret must never be logged.
func (m *TokenManager) Issue(userID int64, email string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("issue token: %w", err)
	}
	return signed, expiresAt, nil
}

// Verify parses and validates a token: HS256 only (blocks algorithm-confusion
// attacks), correct issuer, and a required, non-expired expiry.
func (m *TokenManager) Verify(tokenString string) (Claims, error) {
	var claims Claims
	tok, err := jwt.ParseWithClaims(tokenString, &claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
			}
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !tok.Valid {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}
