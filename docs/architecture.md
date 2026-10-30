# Architecture — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.0 (design for planning — nothing built yet) |
| Date | 2026-10-27 |
| Mandatory stack source | `Chandana 1.docx` mandatory-requirements table |
| Related | `docs/requirements.md` (esp. §9 open questions), `docs/sprint-plan.md`, `docs/api-test-plan.md` |

> **Labeling rule.** Anything not explicitly required by the assignment is tagged **[Proposed]**. The mandatory stack itself is [Explicit]: React.js functional components + hooks · Go REST APIs · PostgreSQL with relationships and constraints · JWT auth + hashed passwords.

## 1. High-level shape

```
┌─────────────────────────────┐        ┌────────────────────────────────┐        ┌────────────────────┐
│ React SPA                   │  HTTP  │ Go REST API                    │  SQL   │ PostgreSQL         │
│ (functional components,     │ ─────► │ router → auth middleware →     │ ─────► │ users              │
│  hooks)                     │  JSON  │ handler → service → repository │  txns  │ categories         │
│ pages · AuthContext ·       │ ◄───── │ input validation · JWT check · │ ◄───── │ products           │
│ api client                  │        │ bcrypt compare · user scoping  │        │ stock_transactions │
└─────────────────────────────┘        └────────────────────────────────┘        └────────────────────┘
```

## 2. Component responsibilities

### React frontend [Explicit: React.js, functional components, hooks]

| Unit | Responsibility |
|---|---|
| Auth pages (Register, Login) | Forms with **frontend validation of required fields** (AC-03); call auth APIs; store token [Proposed: localStorage, see §9]. |
| Protected layout | Redirects to login when no/invalid token; logout clears the token (AC-01, FR-1.1). |
| Categories page | List (FR-2.4) + create/update/delete with confirmations and error display. |
| Products page | List with current stock (FR-3.5), search by name/SKU (FR-3.3), filter by category (FR-3.4), add/edit forms with field validation (AC-03), deactivate action. |
| Stock page | Stock In/Out forms (positive-integer checks AC-09), history list with product/type filters (FR-4.4). |
| Dashboard page | Totals, low-stock list (stock ≤ minimum), five most recent transactions (FR-5.1…FR-5.3). |
| API client + AuthContext | Single fetch wrapper attaching the JWT; centralizes 401 handling and error toasts (AC-15). |

### Go backend [Explicit: Go with REST APIs]

| Layer | Responsibility |
|---|---|
| Router | Maps endpoints (§7); returns 404/405 for undefined routes and methods — this is what makes TC-44 (no transaction edit/delete) hold. [Proposed: standard library `net/http` or chi] |
| Auth middleware | Validates the JWT signature and expiry; injects the authenticated `user_id` into the request context; rejects otherwise (401). |
| Handlers | Parse/validate input **per endpoint** (AC-16); map service results to HTTP codes; produce structured error messages (AC-15 surfacing). |
| Services | Business rules: uniqueness checks, category-delete guard (AC-08), sufficiency check for Stock Out (AC-10). |
| Repository | All SQL; **every query filters by the context `user_id`** (AC-02 / NFR-03); owns transactions for the stock engine. |
| Migrations | Plain SQL migration files applied by a small runner [Proposed]; the migration set itself is submission deliverable D-2. |

### PostgreSQL [Explicit: with appropriate relationships and constraints]

Stores all state. Constraints are the last line of defense behind API validation (§5).

## 3. Data model

All tables get `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`; products/categories also `updated_at`. Timestamps are stored in UTC [Proposed].

**users** — [Proposed fields; the assignment defines no user attributes beyond registration/login]

| Column | Type | Constraints |
|---|---|---|
| id | BIGINT / UUID | PK |
| email | TEXT | NOT NULL, **unique** (case-insensitive) [Proposed identifier choice] |
| password_hash | TEXT | NOT NULL — bcrypt digest only, never plaintext (AC-01) |

**categories**

| Column | Type | Constraints |
|---|---|---|
| id | BIGINT | PK |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| name | TEXT | NOT NULL; **UNIQUE (user_id, lower(trim(name)))** — per-user uniqueness (FR-2.2), case-insensitive per Q4 [Proposed] |

**products**

| Column | Type | Constraints |
|---|---|---|
| id | BIGINT | PK |
| user_id | BIGINT | NOT NULL, FK → users(id) |
| category_id | BIGINT | NOT NULL, FK → categories(id) — makes "no valid category" impossible at DB level (AC-07) |
| name | TEXT | NOT NULL |
| sku | TEXT | NOT NULL; **UNIQUE (user_id, lower(trim(sku)))** (AC-06) |
| unit_price | NUMERIC(10,2) | NOT NULL, CHECK (unit_price >= 0) (AC-04) |
| opening_stock | INTEGER | NOT NULL, CHECK (opening_stock >= 0) (AC-05) |
| current_stock | INTEGER | NOT NULL, CHECK (current_stock >= 0) (AC-05/AC-11 backstop) |
| min_stock | INTEGER | NOT NULL, CHECK (min_stock >= 0) (AC-05) |
| status | TEXT | NOT NULL, CHECK (status IN ('active','inactive')), default 'active' — supports AC-14 and FR-5.1 "active products" |

`current_stock` is a maintained counter initialized to `opening_stock` (Q13 [Proposed]: no initial Stock In row is written).

**stock_transactions**

| Column | Type | Constraints |
|---|---|---|
| id | BIGINT | PK |
| user_id | BIGINT | NOT NULL, FK → users(id) — denormalized from product for cheap per-user scoping [Proposed] |
| product_id | BIGINT | NOT NULL, FK → products(id) |
| type | TEXT | NOT NULL, CHECK (type IN ('IN','OUT')) |
| quantity | INTEGER | NOT NULL, CHECK (quantity > 0) (AC-09 at DB level) |
| remarks | TEXT | NULL — optional per FR-4.2 |
| created_at | TIMESTAMPTZ | NOT NULL — the required timestamp |

Immutability (AC-13): **no API endpoint** exists to update or delete rows. [Proposed optional hardening] a `BEFORE UPDATE OR DELETE` trigger raising an exception, guaranteeing immutability even against manual SQL.

**Relationships:** users 1—N categories; users 1—N products; categories 1—N products; users 1—N stock_transactions; products 1—N stock_transactions. Indexes on `products(user_id)`, `products(user_id, lower(sku))`, `stock_transactions(user_id, created_at DESC)`, `stock_transactions(product_id)` [Proposed].

## 4. Authentication and authorization

- **Registration:** validate input (AC-16) → hash password with **bcrypt** [Proposed algorithm; "hashed passwords" is Explicit] → insert user. Duplicate email → meaningful 409/400 (AC-15).
- **Login:** fetch user by email → `bcrypt.CompareHashAndPassword` → issue **JWT** signed with a secret from an environment variable [Proposed secret management — never hardcoded], expiry ~24 h (Q6 [Proposed]). Claims: `sub` (user id), `exp`, `iat`.
- **Authorization model:** single role; the token's `sub` **is** the tenant boundary. Middleware rejects missing/expired/invalid tokens with 401 (AC-01). Every repository call takes `user_id` from the token context — never from client-supplied payloads (AC-02, NFR-03).
- **Cross-user object access:** requested IDs are always fetched with `WHERE id = $1 AND user_id = $2`; a miss returns **404** [Proposed: 404 over 403 to avoid existence leakage] — verified by TC-10…TC-12.
- **Logout:** [Proposed per Q5] client discards the token; no server-side blacklist; short expiry bounds exposure. Endpoint included to satisfy FR-1.1.

## 5. Validation model (defense in depth)

| Layer | What it enforces | ACs |
|---|---|---|
| Frontend | Required fields present; non-negative price; non-negative integer opening/min stock; positive integer quantities; format errors — blocks submission early | AC-03, AC-15 |
| Backend handlers | Re-validates **everything** the frontend checks, on every request (AC-16); returns field-level 400 messages | AC-04…AC-10, AC-16 |
| Database | FKs, uniqueness, CHECK constraints (§3) — final backstop even against bugs | AC-04…AC-07, AC-09, NFR-09 |

## 6. API boundaries

Prefix `/api/v1` [Proposed]. Full scenario matrix lives in `docs/api-test-plan.md`.

| Method & path | Auth | Purpose / key rules |
|---|---|---|
| GET /healthz | no | Liveness [Proposed] |
| POST /api/v1/auth/register | no | Create user; unique email; hashed password |
| POST /api/v1/auth/login | no | Verify credentials; issue JWT |
| POST /api/v1/auth/logout | yes | Client-side session end [Proposed semantics, Q5] |
| GET /api/v1/categories | yes | List own categories (FR-2.4) |
| POST /api/v1/categories | yes | Create; per-user unique name (FR-2.2) |
| GET/PUT/DELETE /api/v1/categories/{id} | yes | Own category only; DELETE guarded by AC-08 (+Q3 rule) |
| GET /api/v1/products?search=&category_id= | yes | Own products; current stock included; search name/SKU (FR-3.3–FR-3.5) |
| POST /api/v1/products | yes | Valid own category (AC-07); SKU unique per user (AC-06); price ≥ 0 (AC-04); opening/min stock non-negative integers (AC-05) |
| GET/PUT /api/v1/products/{id} | yes | Own product only (FR-3.2); SKU uniqueness re-checked on edit (Q14) |
| POST /api/v1/products/{id}/deactivate | yes | AC-14 [Proposed endpoint shape] |
| POST /api/v1/transactions | yes | type IN/OUT; quantity positive integer (AC-09); OUT guarded by sufficiency (AC-10); atomic with stock update (AC-12) |
| GET /api/v1/transactions?product_id=&type= | yes | Own history, filtered (FR-4.4) |
| PUT/DELETE /api/v1/transactions/{id} | — | **Intentionally absent** (AC-13); requests must 404/405 (TC-44) |
| GET /api/v1/dashboard/summary | yes | Totals (active products), category count, low-stock count + list (stock ≤ min), 5 most recent transactions (FR-5.1…FR-5.3) |

## 7. Transaction safety and concurrency ⚠ (high-risk design area)

**Requirement:** AC-11 (no incorrect/negative stock under concurrency) and AC-12 (atomic stock changes in PostgreSQL).

**Designed mechanism** [Proposed — the assignment requires the outcome, not the mechanism]:

1. Open a DB transaction.
2. `INSERT INTO stock_transactions (...)` — movement record.
3. `UPDATE products SET current_stock = current_stock ± $qty WHERE id = $pid AND user_id = $uid AND status = 'active' AND ($type = 'IN' OR current_stock >= $qty) RETURNING current_stock`.
   - No returned row → rollback: insufficient stock (AC-10), inactive product, or foreign product (AC-02).
   - The conditional `WHERE` takes the row lock and re-checks the guard atomically, so two concurrent Stock Outs cannot both pass the same last unit.
4. Commit. Steps 2–3 commit together or not at all (AC-12).
5. `CHECK (current_stock >= 0)` (§3) is a DB-level backstop against any coding error.

Serializable isolation is **not** required for this pattern — the single conditional UPDATE serializes access to the row; this is stated deliberately so the simpler correct design is chosen over a heavier one. Verified empirically by TC-45 (boundary burst) and TC-46 (mixed In/Out sum check), plus TC-47 (failed op leaves no row and no change).

## 8. Dashboard computation

Single endpoint returning: `COUNT(products WHERE status='active')`, `COUNT(categories)`, low-stock = products with `current_stock <= min_stock` [Proposed: active only, Q2], recent transactions = `ORDER BY created_at DESC LIMIT 5` scoped to the user.

## 9. Security risks and mitigations

| Risk | Mitigation | Residual |
|---|---|---|
| Plaintext/weak password storage | bcrypt via reputable Go package [Proposed]; never log passwords | Low |
| JWT secret leakage | Secret only via env var; strong random value; never committed; `.env` git-ignored | Ops discipline |
| Token theft (XSS) with localStorage [Proposed, Q11] | No `dangerouslySetInnerHTML`; httpOnly-cookie variant listed as alternative | Accepted for this scope |
| IDOR / cross-user access (AC-02) | `user_id` always from token context; every query scoped; 404 on foreign IDs; TC-10…TC-12 | Low |
| SQL injection | Parameterized queries only; no string-built SQL | Low |
| Negative/incorrect stock under race (AC-11) | §7 mechanism + DB CHECK | Low, tested |
| Audit tampering (AC-13) | No write endpoints; optional DB trigger [Proposed] | Low |
| CORS misconfiguration | Dev CORS restricted to the frontend origin only [Proposed, Q12] | Low |
| Secrets in repo | `.env` in `.gitignore`; README documents required env vars without values | Ops discipline |

## 10. Explicit-vs-proposed decision log

| Decision | Status |
|---|---|
| React functional components + hooks; Go REST; PostgreSQL; JWT + hashed passwords | **Explicit** (assignment) |
| bcrypt, JWT 24 h, secret via env, localStorage token, 404 for foreign IDs, logout = client discard | Proposed (Q5/Q6/Q11) |
| Per-user case-insensitive uniqueness (SKU, category name) | Proposed (Q4) |
| `current_stock` maintained counter; opening stock creates no movement row | Proposed (Q13) |
| status active/inactive; deactivate-only when history exists; hard delete only with zero history; reactivation allowed | Proposed (Q2/Q14; AC-14 is Explicit) |
| Category delete blocked when *any* product references it (stricter than AC-08's "active") | Proposed (Q3) |
| `/api/v1` prefix, endpoint shapes, deactivation endpoint, health endpoint | Proposed |
| Substring search; NUMERIC(10,2); no pagination; optional limit on history | Proposed (Q7/Q8/Q9) |
| Go router choice, migration file mechanism, Vite, plain CSS, dev ports/CORS | Proposed (Q10/Q12) |
| PostgreSQL driver `github.com/jackc/pgx/v5` + `pgxpool` (max 5 conns, bounded startup ping); Go module path `inventory-backend`; stdlib `net/http` ServeMux with Go 1.22+ method routing (`GET /healthz`) | **Chosen at T-104** — a PG driver is essential and pgx/v5 is the maintained standard; module rename is trivial if a GitHub-hosted path is preferred later |
| `.env` loading: minimal in-repo parser in `internal/config` (no godotenv); OS environment variables take precedence over file values; server reads `.env` or `../.env` | **Chosen at T-104** — avoids an extra dependency for a simple KEY=VALUE format |
| Health endpoint `GET /healthz`: 200 + `{"status":"ok","database":"up"}` when the pool pings, 503 + `{"status":"degraded","database":"down"}` on failed/slow ping (2 s budget) | **Chosen at T-104** |
| HTTP timeouts (read header 5 s / read 10 s / write 15 s / idle 60 s), graceful shutdown on SIGINT/SIGTERM (5 s budget), `APP_PORT` env var (code default 8080; **local dev uses 8081** — 8080 occupied on this machine, see known-issues ENV-003) | **Chosen at T-104** |
| Migrations: numbered plain-SQL files (`backend/migrations/NNNN_name.{up,down}.sql`) embedded in the binary, applied by `backend/cmd/migrate` (idempotent, per-migration transaction, `schema_migrations` bookkeeping, `-verify` read-only summary; up-only — down files exist for manual dev resets); **composite foreign keys** `(category_id, user_id) → categories(id, user_id)` and `(product_id, user_id) → products(id, user_id)` enforce per-user isolation at the database level (defense in depth for AC-02/AC-07); DB-level immutability triggers on `stock_transactions` (AC-13) | **Chosen at T-105** — implements the §3 data model; the composite FKs and triggers harden constraints already implied by AC-02/AC-07/AC-13 |
| Auth implementation: bcrypt (`golang.org/x/crypto` — Go team's reference implementation; stdlib has none) + HS256 JWTs (`github.com/golang-jwt/jwt/v5`, the maintained standard — hand-rolled verification rejected as a security risk); `JWT_SECRET` env var (required, ≥32 chars, never logged); 24 h expiry (Q6); password policy 8–72 chars (bcrypt's input limit; min length is a documented assumption); login responses identical for unknown emails vs wrong passwords with bcrypt timing equalization (anti-enumeration); logout = client-side token discard, **no server-side revocation of stateless JWTs** (Q5 — documented limitation); uniform `{"error": message}` responses (AC-15/16) | **Chosen at T-106** |
| Frontend: Vite 8 + React 19 (functional components + hooks only), plain CSS, no router/UI framework yet (multi-page navigation arrives with later tasks); dev CORS avoided via **Vite proxy** `/api → :8081` (same-origin requests; production would use one reverse proxy); session token + email in `localStorage` (Q11 trade-off), 401 responses clear the local session; frontend tests via Vitest + Testing Library + jsdom rendering the real App with mocked fetch | **Chosen at T-108** |

## 11. Intentionally not designed (scope restraint)

Deployment/hosting, CI pipelines, pagination, role hierarchies, notifications, multi-currency, UI component library — none appear in the assignment (`docs/requirements.md` §2). They will not be designed unless the user requests them.
