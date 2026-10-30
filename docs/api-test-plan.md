# API Test Plan — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.0 (planning phase) |
| Date | 2026-10-27 |
| Execution status | **Nothing executed yet** — endpoints do not exist. This plan defines *how* API-level testing will run during T-301…T-303. |
| Related | `docs/test-cases.md` (case definitions), `docs/evidence-plan.md` (recording), `docs/architecture.md` §6 (endpoint design) |

## 1. Purpose and scope

Verify every REST endpoint's behavior directly at the HTTP level, independent of the UI — covering authentication, authorization (cross-user), validation, business-rule and error responses, immutability of transactions, and concurrency. Backend validation must hold even where the frontend also validates (AC-16).

## 2. Environment and tooling

| Item | Plan |
|---|---|
| Base URL | `http://localhost:8080` [Proposed port, Q12 — confirmed at T-104] |
| Tool | `curl` from Git Bash (availability verified read-only at T-101; curl ships with Git Bash). A GUI client (Postman/Insomnia) is optional if already installed [Proposed]. |
| Auth | `Authorization: Bearer <JWT>` header captured from the login response at test time |
| Data reset | Re-apply migrations to a fresh database between major suites if state gets polluted [Proposed] |

No test scripts are written in the planning phase; helper commands are created only during execution tasks.

## 3. Test data setup [Proposed values]

| Fixture | Value |
|---|---|
| User A | `owner-a@example.test` — owns categories *Electronics*, *Stationery* |
| User B | `owner-b@example.test` — owns category *Electronics* (same name as A's — exercises per-user uniqueness, TC-15) |
| Product A1 | SKU `SKU-A1`, stock 10, min 3, price 10.50, category Electronics (A) |
| Product A2 | SKU `SKU-A2`, stock 2, min 5, price 5.00, category Stationery (A) — **low stock by design** |
| Product B1 | SKU `SKU-A1` (same SKU as A's — exercises per-user SKU scope, TC-26), stock 4, min 1 |
| Concurrency target | Product A1 driven to a known boundary stock (5) before TC-45 |

This fixture makes dashboard expectations exact: active=3 (before deactivation), categories A=2, low-stock=1 (A2).

## 4. Endpoint scenario matrix

Every row is executed as an actual HTTP request; expected HTTP codes follow `docs/architecture.md` §6. Linked test-case IDs define the exact payloads.

| Endpoint | No auth | Invalid/expired token | Cross-user access (B on A's data) | Happy path | Invalid payload | Not found |
|---|---|---|---|---|---|---|
| POST /api/v1/auth/register | n/a | n/a | n/a | TC-01 → 2xx | TC-02/TC-03 → 400/409 | n/a |
| POST /api/v1/auth/login | n/a | n/a | n/a | TC-05 → 2xx + JWT | TC-06 → 401 | n/a |
| POST /api/v1/auth/logout | 401 | 401 | n/a | 2xx | — | n/a |
| GET/POST /api/v1/categories | 401 (TC-07) | 401 (TC-08) | lists never contain foreign rows (TC-12) | TC-13 → 201 | TC-14/TC-20 → 400/409 | — |
| GET/PUT/DELETE /api/v1/categories/{id} | 401 | 401 | TC-10/TC-11 → 404 | TC-16/TC-17 → 2xx | — | unknown id → 404 |
| DELETE /api/v1/categories/{id} (guarded) | 401 | 401 | 404 | TC-17 → 2xx | **TC-18 → 409/400 + message** | — |
| GET /api/v1/products?search=&category_id= | 401 | 401 | only B's rows (TC-12) | TC-30…TC-33 → 200 | — | — |
| POST /api/v1/products | 401 | 401 | foreign category id → 400 (TC-27) | TC-21 → 201 | TC-22…TC-25, TC-28 → 400/409 | — |
| GET/PUT /api/v1/products/{id} | 401 | 401 | TC-10/TC-11 → 404 | TC-29 → 200 | TC-22…TC-24 → 400 | unknown id → 404 |
| POST /api/v1/products/{id}/deactivate | 401 | 401 | 404 | TC-34 → 200 | — | — |
| POST /api/v1/transactions | 401 | 401 | foreign product id → 404 (AC-02) | TC-36/TC-37 → 201 | TC-38…TC-40 → 400/409 | — |
| GET /api/v1/transactions?product_id=&type= | 401 | 401 | only B's rows (TC-12) | TC-42/TC-43 → 200 | — | — |
| PUT/DELETE /api/v1/transactions/{id} | — | — | — | **must NOT exist: expect 404/405 (TC-44)** | — | — |
| GET /api/v1/dashboard/summary | 401 | 401 | only B's numbers (TC-12) | TC-49…TC-51 → 200 | — | — |

## 5. Concurrency test method ⚠ (high-risk — AC-11/AC-12)

Documented method; the concrete commands are produced during T-303, not before:

1. Set product A1 stock to exactly **5** (drive it with recorded movements and read it back from the DB).
2. Fire **10 simultaneous** `POST /api/v1/transactions` Stock Out qty 1 (e.g., curl background jobs + `wait`, or an equivalent parallel launcher available in the verified toolchain).
3. **Expected:** exactly 5 responses success, 5 responses failure (insufficient stock); DB stock reads **0**; exactly 5 new transaction rows; **no negative value at any point** (TC-45).
4. Repeat with a mixed parallel batch of IN and OUT; verify `final stock == opening + ΣIN − ΣOUT` computed from the transactions table (TC-46).
5. Evidence: full terminal transcript + before/after `SELECT current_stock` outputs + success/fail counts (`docs/evidence-plan.md` §5).

## 6. Recording rules

- Each request's command and response (status + body) are saved to `evidence/` at execution time, per the naming convention in `docs/evidence-plan.md`.
- A case is marked Pass/Fail only after execution; failures are logged in `docs/known-issues.md` with the failing request attached.
- No result is written in advance. `Not Run` stays until the request actually fires.

## 7. Entry and exit criteria

- **Entry (T-301):** all endpoints from `docs/architecture.md` §6 implemented; migrations applied; two users + fixtures creatable via API.
- **Exit:** every matrix cell executed or explicitly blocked with a known-issue entry; all executed cases have recorded actual results and evidence; AC coverage confirmed in `docs/traceability-matrix.md`.
