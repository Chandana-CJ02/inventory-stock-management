// Package auth provides password hashing (bcrypt) and JWT session-token
// handling for the API.
package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost // 10

// dummyHash is compared against during logins for unknown accounts so that
// response timing does not reveal whether an email is registered.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("timing-equalizer-not-a-real-password"), bcryptCost)

// HashPassword returns the bcrypt hash of the password. The hash is safe to
// store; the plaintext is never persisted or logged.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword reports whether password matches the stored bcrypt hash.
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// VerifyDummy burns the same bcrypt cost as VerifyPassword for accounts that
// do not exist, keeping login timing uniform.
func VerifyDummy(password string) {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
}
