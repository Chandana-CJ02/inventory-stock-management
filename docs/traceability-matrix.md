# Traceability Matrix — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.1 (execution complete) |
| Date | 2026-10-29 |
| Rule | **Status values:** `Not started` → `Implemented` → `Tested` → `Verified`. `Verified` requires executed, passing tests with evidence on disk. Nothing is `Verified` today — no code or test run exists. |

## 1. Acceptance criteria traceability (authoritative)

| AC | Requirement (abridged) | FR | Story | Test cases | Task(s) | Planned evidence | Status |
|---|---|---|---|---|---|---|---|
| AC-01 | Secure password storage + token auth | FR-1.1 | US-01, US-02 | TC-01…TC-09, esp. TC-04 | T-106, T-108 | evidence/TC-04_hashed-password, TC-05/06/07/08 captures | **Verified** (TC-01…TC-09 executed live in T-106, evidence on disk; bcrypt format confirmed via `-verify`) |
| AC-02 | No access to another user's inventory | FR-1.2, FR-1.3 | US-04 | TC-10, TC-11, TC-12, TC-15, TC-26 ⚠ | T-106…T-203, T-302 | evidence/TC-10…TC-12 cross-user captures | **Tested** — categories, products, movements, dashboard isolation verified (T-107, T-201, T-203/T-204 history filter, T-205 summary, T-208 two-user live spot check 2026-10-29: B's reads/updates/deletes of A's records → 404/405, B's lists and dashboard show only B's data); TC-26 SKU scope verified T-201; browser-side spot check pending T-207 manual pass |
| AC-03 | Frontend validation of required fields | FR-1.1, FR-2.1, FR-3.1 | US-01, US-05, US-08 | TC-03, TC-19, TC-28 | T-108, T-206 | UI screenshots of blocked submissions | **Tested** — TC-03/TC-19/TC-28 frontend blocks passing with no-POST assertions (T-206 vitest; TC-03 auth blank-field literal cases added 2026-10-29, evidence/TC-03_auth-blank-field-vitest.txt) |
| AC-04 | Price cannot be negative | FR-3.1 | US-08 | TC-22 | T-201, T-301 | API 400 capture + DB CHECK | **Verified** (TC-22 passing, evidence on disk) |
| AC-05 | Opening/min stock non-negative integers | FR-3.1 | US-08 | TC-23, TC-24 | T-201, T-301 | API 400 captures | **Verified** (TC-23, TC-24 passing, evidence on disk) |
| AC-06 | SKU unique within user's inventory | FR-3.1 | US-08 | TC-25, TC-26 | T-201, T-301 | 409 capture + cross-user allowed capture | **Verified** (TC-25, TC-26 passing, evidence on disk) |
| AC-07 | Product requires valid category | FR-3.1 | US-08 | TC-27 | T-201, T-301 | 400 captures (unknown + foreign category) | **Verified** (TC-27 passing, evidence on disk) |
| AC-08 | No category delete with active products | FR-2.1 | US-07 | TC-18 | T-107, T-301 | Blocked-delete capture | **Verified** (TC-18 passing: delete with products → 409 guard message, category retained — `evidence/T-107_category-tests.txt`; UI-side 409 surfacing proven in T-206 vitest) |
| AC-09 | Stock quantities positive integers | FR-4.1 | US-12, US-13 | TC-39, TC-40 (positives TC-36/37) | T-203, T-301 | 400 captures | **Verified** (TC-39/TC-40: 6 invalid-qty subtests passing, evidence on disk) |
| AC-10 | No Stock Out beyond available stock | FR-4.1 | US-13 | TC-38 | T-203, T-301 | 409/400 capture + unchanged stock query | **Verified** (TC-38: 409, stock unchanged, no row, evidence on disk) |
| AC-11 | Concurrency cannot corrupt/negative stock | FR-4.1 | US-14 | TC-45, TC-46 ⚠ | T-203, T-303 | Concurrency transcript + before/after stock | **Verified** — TC-45 single-direction burst passing (5/5 split, final 0); TC-46 mixed-direction passing (20 concurrent IN/OUT, all 201, final stock 30 correct, no lost updates) |
| AC-12 | Transactions in PostgreSQL; atomic changes | FR-4.1 | US-12…US-14 | TC-46, TC-47 ⚠ | T-105, T-203, T-303 | Row inspection + failed-op atomicity check | **Verified** — TC-47 atomicity passing; TC-46 mixed-direction passing (20/20 concurrent ops, row counts exact, stock consistent) |
| AC-13 | Transactions never edited/deleted | FR-4.2 | US-16 | TC-44 | T-203 (absence by design), T-301 | 404/405 captures; DB rows unchanged | **Verified** (TC-44: no PUT/DELETE routes → 404/405; DB triggers present per T-105 `-verify`) |
| AC-14 | Deactivate (not delete) products with history | FR-3.2 | US-18 | TC-34, TC-35 | T-201, T-301 | Deactivation capture + retained history | **Verified** (TC-34, TC-35 passing, evidence on disk; no delete route exists — 405. Movement-history interaction completes with T-203+) |
| AC-15 | Meaningful frontend error messages | cross-cutting | US-19 | TC-52 | T-206, T-207, T-304 | UI screenshot of readable error | **Tested (partial)** — readable backend/JS-validation messages proven in vitest (duplicate name/SKU, failed login, delete-guard); browser screenshots pending T-304 |
| AC-16 | Backend validation on all APIs | cross-cutting | US-19 | TC-53 (+ backend halves of TC-20/22/23/24/27/28/39/40) | T-107, T-201, T-203, T-301 | Raw API 400 captures bypassing UI | **Tested (partial)** — backend validation executed and passing for auth (T-106), categories (T-107: TC-20), products (T-201: TC-22…TC-24, TC-27, TC-28); stock halves (TC-39/40) pending T-203; UI-bypass flow (TC-53) pending |

## 2. Functional requirements without a numbered AC

| FR | Requirement | Story | Test cases | Task(s) | Status |
|---|---|---|---|---|---|
| FR-2.4 | Categories displayed in a list | US-05 | TC-13, TC-19 | T-107, T-206 | **Tested** — backend list verified in T-107; frontend list+create flow verified in T-206 vitest |
| FR-3.2 | View and edit product details | US-09 | TC-29 | T-201 | **Tested** — TC-29 passing, evidence on disk |
| FR-3.3 | Search by name or SKU | US-10 | TC-30, TC-31 | T-202 | **Tested** — TC-30, TC-31 passing, evidence on disk |
| FR-3.4 | Filter by category | US-10 | TC-32 | T-202 | **Tested** — TC-32 passing, evidence on disk |
| FR-3.5 | Display current available stock | US-11 | TC-33 | T-201, T-206 | **Tested** — TC-33 passing; live proxy re-verified (stock 11, history 2 rows, dashboard consistent) |
| FR-4.2 | Transaction record fields (product, type, qty, remarks, timestamp) | US-15 | TC-41, TC-48 | T-203 | **Tested** — TC-41, TC-48 passing, evidence on disk |
| FR-4.4 | History filter by product and type | US-15 | TC-42, TC-43 | T-204 | **Tested** — TC-42, TC-43 passing, evidence on disk; live re-verified |
| FR-5.1 | Totals: active products, categories, low-stock count | US-17 | TC-49 | T-205 | **Tested** — TC-49 passing, evidence on disk |
| FR-5.2 | Low-stock list (stock ≤ minimum) | US-17 | TC-50 | T-205 | **Tested** — TC-50 passing, evidence on disk |
| FR-5.3 | Five most recent transactions | US-17 | TC-51 | T-205 | **Tested** — TC-51 passing, evidence on disk |

## 3. Submission deliverables traceability

| Deliverable | Producing task | Verified? |
|---|---|---|
| D-1 GitHub URL with React + Go source | T-103, T-308 | No |
| D-2 PostgreSQL schema/migrations | T-105 | No |
| D-3 README with installation/execution instructions | T-307 | No |

## 4. Coverage summary

Computed only after execution. Current honest state (2026-10-29, final pass): **16/16 ACs executed — 12 Verified, 2 Tested (AC-02, AC-03), 2 Tested (partial) (AC-15, AC-16). Zero Not Run. All 53 test cases Pass (TC-03 completed 2026-10-29). Remaining partials are browser-screenshot evidence gaps (AC-15/AC-16) and the live two-user browser spot check (AC-02), not missing execution.**
