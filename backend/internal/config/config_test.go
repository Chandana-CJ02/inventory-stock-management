package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testPassword = "test-password-not-real" // fake value for tests only

func setValidEnv(t *testing.T) {
	t.Helper()
	t.Setenv("PGHOST", "localhost")
	t.Setenv("PGPORT", "5432")
	t.Setenv("PGDATABASE", "testdb")
	t.Setenv("PGUSER", "testuser")
	t.Setenv("PGPASSWORD", testPassword)
	t.Setenv("PGSSLMODE", "")
	t.Setenv("APP_PORT", "")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
}

func missingEnvFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "missing.env")
}

func TestLoadValidDefaults(t *testing.T) {
	setValidEnv(t)
	cfg, err := Load(missingEnvFile(t))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Host != "localhost" || cfg.Database != "testdb" || cfg.User != "testuser" {
		t.Fatalf("unexpected config: host=%q db=%q user=%q", cfg.Host, cfg.Database, cfg.User)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("SSLMode default = %q, want %q", cfg.SSLMode, "disable")
	}
	if cfg.AppPort != "8080" {
		t.Errorf("AppPort default = %q, want %q", cfg.AppPort, "8080")
	}
}

func TestLoadMissingRequired(t *testing.T) {
	for _, tc := range []struct{ name, key string }{
		{"missing host", "PGHOST"},
		{"missing port", "PGPORT"},
		{"missing database", "PGDATABASE"},
		{"missing user", "PGUSER"},
		{"missing password", "PGPASSWORD"},
		{"missing jwt secret", "JWT_SECRET"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			setValidEnv(t)
			t.Setenv(tc.key, "")
			_, err := Load(missingEnvFile(t))
			if err == nil {
				t.Fatalf("Load() with empty %s: expected error", tc.key)
			}
			if !strings.Contains(err.Error(), tc.key) {
				t.Fatalf("error %q should name %s", err, tc.key)
			}
			if strings.Contains(err.Error(), testPassword) {
				t.Fatal("error must never contain values")
			}
		})
	}
}

func TestLoadInvalidPort(t *testing.T) {
	setValidEnv(t)
	t.Setenv("PGPORT", "99999")
	if _, err := Load(missingEnvFile(t)); err == nil || !strings.Contains(err.Error(), "PGPORT") {
		t.Fatalf("expected PGPORT error, got %v", err)
	}
}

func TestLoadInvalidSSLMode(t *testing.T) {
	setValidEnv(t)
	t.Setenv("PGSSLMODE", "bogus")
	if _, err := Load(missingEnvFile(t)); err == nil || !strings.Contains(err.Error(), "PGSSLMODE") {
		t.Fatalf("expected PGSSLMODE error, got %v", err)
	}
}

func TestLoadInvalidAppPort(t *testing.T) {
	setValidEnv(t)
	t.Setenv("APP_PORT", "not-a-number")
	if _, err := Load(missingEnvFile(t)); err == nil || !strings.Contains(err.Error(), "APP_PORT") {
		t.Fatalf("expected APP_PORT error, got %v", err)
	}
}

func TestLoadShortJWTSecret(t *testing.T) {
	setValidEnv(t)
	t.Setenv("JWT_SECRET", "too-short")
	_, err := Load(missingEnvFile(t))
	if err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("expected JWT_SECRET length error, got %v", err)
	}
}

func TestDotEnvAppliedForUnsetVars(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.env")
	content := "# comment\nPGHOST=filehost\nPGPORT=6543\nPGDATABASE=filedb\nPGUSER=fileuser\nPGPASSWORD=" +
		testPassword + "\nAPP_PORT=9090\nJWT_SECRET=" + strings.Repeat("f", 32) + "\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	keys := []string{"PGHOST", "PGPORT", "PGDATABASE", "PGUSER", "PGPASSWORD", "PGSSLMODE", "APP_PORT", "JWT_SECRET"}
	for _, k := range keys {
		if err := os.Unsetenv(k); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		for _, k := range keys {
			_ = os.Unsetenv(k)
		}
	})

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Host != "filehost" || cfg.Port != "6543" || cfg.Database != "filedb" || cfg.AppPort != "9090" {
		t.Fatalf("dotenv values not applied: host=%q port=%q db=%q app_port=%q",
			cfg.Host, cfg.Port, cfg.Database, cfg.AppPort)
	}
	if cfg.SSLMode != "disable" {
		t.Errorf("SSLMode default = %q, want disable", cfg.SSLMode)
	}
}

func TestEnvOverridesDotEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.env")
	if err := os.WriteFile(path, []byte("PGHOST=filehost\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	setValidEnv(t) // OS environment is set and must win

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Host != "localhost" {
		t.Fatalf("OS env should take precedence: host = %q, want localhost", cfg.Host)
	}
}
