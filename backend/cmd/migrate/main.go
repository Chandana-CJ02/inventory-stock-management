// Command migrate applies pending SQL migrations to the configured database.
//
// Usage (run from the backend/ directory):
//
//	go run ./cmd/migrate           apply all pending migrations
//	go run ./cmd/migrate -verify   print a read-only schema summary
//
// Each migration is applied at most once, inside its own transaction, and is
// recorded in schema_migrations. Re-running is safe: already-applied versions
// are skipped. Credentials come from the environment / .env and are never
// logged.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"inventory-backend/internal/config"
	"inventory-backend/internal/database"
	"inventory-backend/migrations"
)

type migration struct {
	version int
	name    string
	sql     string
}

func main() {
	verify := flag.Bool("verify", false, "print a read-only schema summary instead of migrating")
	flag.Parse()
	if err := run(*verify); err != nil {
		log.Fatal(err)
	}
}

func run(verify bool) error {
	cfg, err := config.Load(".env", "../.env")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	connCfg, err := pgx.ParseConfig(database.DSN(cfg))
	if err != nil {
		return fmt.Errorf("migrate: parse config: %w", err)
	}
	// Migration files contain multiple SQL statements (and trigger functions
	// with embedded semicolons), so the simple query protocol is required.
	connCfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgx.ConnectConfig(ctx, connCfg)
	if err != nil {
		return fmt.Errorf("migrate: connect: %w", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	if verify {
		return verifySchema(ctx, conn)
	}
	return migrate(ctx, conn)
}

func migrate(ctx context.Context, conn *pgx.Conn) error {
	if _, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    BIGINT PRIMARY KEY,
			name       TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("migrate: ensure schema_migrations: %w", err)
	}

	list, err := loadMigrations()
	if err != nil {
		return err
	}

	var current int64
	if err := conn.QueryRow(ctx,
		`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&current); err != nil {
		return fmt.Errorf("migrate: read current version: %w", err)
	}

	applied, skipped := 0, 0
	for _, m := range list {
		if int64(m.version) <= current {
			skipped++
			continue
		}
		label := fmt.Sprintf("%04d_%s", m.version, m.name)
		tx, err := conn.Begin(ctx)
		if err != nil {
			return fmt.Errorf("migrate: begin %s: %w", label, err)
		}
		if _, err := tx.Exec(ctx, m.sql); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migrate: apply %s: %w", label, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`,
			m.version, m.name); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("migrate: record %s: %w", label, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("migrate: commit %s: %w", label, err)
		}
		log.Printf("migrate: applied %s", label)
		applied++
	}
	log.Printf("migrate: done (%d applied, %d already current)", applied, skipped)
	return nil
}

func loadMigrations() ([]migration, error) {
	entries, err := migrations.FS.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("migrate: read embedded migrations: %w", err)
	}
	byVersion := make(map[int]migration)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		base := strings.TrimSuffix(name, ".up.sql")
		idx := strings.Index(base, "_")
		if idx <= 0 {
			return nil, fmt.Errorf("migrate: unexpected file name %q (want NNNN_name.up.sql)", name)
		}
		version, err := strconv.Atoi(base[:idx])
		if err != nil || version <= 0 {
			return nil, fmt.Errorf("migrate: bad version number in file name %q", name)
		}
		sqlBytes, err := migrations.FS.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("migrate: read %q: %w", name, err)
		}
		byVersion[version] = migration{version: version, name: base[idx+1:], sql: string(sqlBytes)}
	}

	out := make([]migration, 0, len(byVersion))
	for _, m := range byVersion {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })
	return out, nil
}

// verifySchema prints a read-only summary of the applied schema. It only reads
// catalogs; no data or objects are modified.
func verifySchema(ctx context.Context, conn *pgx.Conn) error {
	printRows := func(title, query string) error {
		fmt.Println("verify:", title)
		rows, err := conn.Query(ctx, query)
		if err != nil {
			return fmt.Errorf("verify %s: %w", title, err)
		}
		defer rows.Close()
		for rows.Next() {
			var a, b string
			if err := rows.Scan(&a, &b); err != nil {
				return err
			}
			fmt.Printf("  %-24s %s\n", a, b)
		}
		return rows.Err()
	}

	if err := printRows("tables:",
		`SELECT table_name, table_type FROM information_schema.tables
		 WHERE table_schema = 'public' AND table_type = 'BASE TABLE' ORDER BY table_name`); err != nil {
		return err
	}
	if err := printRows("constraints per table:",
		`SELECT c.relname, COUNT(*) FROM pg_constraint con
		 JOIN pg_class c ON con.conrelid = c.oid
		 WHERE c.relnamespace = 'public'::regnamespace
		 GROUP BY c.relname ORDER BY c.relname`); err != nil {
		return err
	}
	if err := printRows("indexes:",
		`SELECT tablename, indexname FROM pg_indexes
		 WHERE schemaname = 'public'
		 ORDER BY tablename, indexname`); err != nil {
		return err
	}
	if err := printRows("stock_transactions triggers:",
		`SELECT tgname, 'enabled' FROM pg_trigger
		 WHERE tgrelid = 'stock_transactions'::regclass AND NOT tgisinternal
		 ORDER BY tgname`); err != nil {
		return err
	}
	// TC-04 evidence helper: proves passwords are stored as bcrypt hashes
	// ('$2'-prefixed) without ever printing the hash itself.
	if err := printRows("users (bcrypt hash format):",
		`SELECT email, (password_hash LIKE '$2%')::text FROM users ORDER BY id`); err != nil {
		return err
	}

	fmt.Println("verify: applied migrations")
	rows, err := conn.Query(ctx, `SELECT version, name FROM schema_migrations ORDER BY version`)
	if err != nil {
		return fmt.Errorf("verify applied migrations: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var version int64
		var name string
		if err := rows.Scan(&version, &name); err != nil {
			return err
		}
		fmt.Printf("  %04d            %s\n", version, name)
	}
	return rows.Err()
}
