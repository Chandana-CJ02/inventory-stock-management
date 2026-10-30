package auth

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "unit-test-secret-0123456789abcdef0123456789abcdef"

func TestHashVerifyRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("hash %q is not a bcrypt hash", hash)
	}
	if !VerifyPassword(hash, "correct horse battery staple") {
		t.Fatal("VerifyPassword() = false for the correct password")
	}
	if VerifyPassword(hash, "wrong password") {
		t.Fatal("VerifyPassword() = true for a wrong password")
	}
}

func TestHashesAreSalted(t *testing.T) {
	a, _ := HashPassword("same-password")
	b, _ := HashPassword("same-password")
	if a == b {
		t.Fatal("two hashes of the same password are identical — salt missing")
	}
}

func TestTokenIssueVerifyRoundTrip(t *testing.T) {
	m := NewTokenManager(testSecret, time.Hour)
	token, expiresAt, err := m.Issue(42, "user@example.test")
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if time.Until(expiresAt) <= 0 {
		t.Fatal("expiresAt is not in the future")
	}
	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if claims.Email != "user@example.test" {
		t.Errorf("claims.Email = %q", claims.Email)
	}
	if id, err := strconv.ParseInt(claims.Subject, 10, 64); err != nil || id != 42 {
		t.Errorf("claims.Subject = %q, want 42", claims.Subject)
	}
}

func TestVerifyRejectsBadTokens(t *testing.T) {
	m := NewTokenManager(testSecret, time.Hour)

	token, _, err := m.Issue(1, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}

	// Wrong-secret token.
	other := NewTokenManager("another-secret-0123456789abcdef0123456789abcdef", time.Hour)
	wrongSig, _, err := other.Issue(1, "user@example.test")
	if err != nil {
		t.Fatal(err)
	}

	// Tampered signature (flip a mid-segment character). The final base64url
	// character carries zero-padding bits, so flipping it can leave the decoded
	// signature bytes unchanged and make this check pass vacuously; a mid
	// character carries only data bits, so the flip always alters the bytes.
	sigStart := strings.LastIndex(token, ".") + 1
	mid := sigStart + (len(token)-sigStart)/2
	flip := byte('A')
	if token[mid] == 'A' {
		flip = 'B'
	}
	tampered := token[:mid] + string(flip) + token[mid+1:]

	for _, tc := range []struct {
		name  string
		token string
	}{
		{"garbage", "not.a.jwt"},
		{"empty", ""},
		{"wrong secret", wrongSig},
		{"tampered signature", tampered},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := m.Verify(tc.token); err == nil {
				t.Fatalf("Verify(%q) accepted an invalid token", tc.name)
			}
		})
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	expired := NewTokenManager(testSecret, -time.Minute) // already expired on issue
	token, _, err := expired.Issue(1, "user@example.test")
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if _, err := NewTokenManager(testSecret, time.Hour).Verify(token); err == nil {
		t.Fatal("expired token accepted")
	}
}

func TestVerifyRejectsForeignSigningMethod(t *testing.T) {
	// Build a token with a different algorithm but the same secret.
	claims := Claims{Email: "user@example.test"}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour))
	tok := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := tok.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	m := NewTokenManager(testSecret, time.Hour)
	if _, err := m.Verify(signed); err == nil {
		t.Fatal("HS512 token accepted although only HS256 is allowed")
	}
}
