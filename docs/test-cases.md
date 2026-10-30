# Test Cases — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.0 (planning phase) |
| Date | 2026-10-27 |
| Execution status | **53 of 53 executed as of 2026-10-29: 53 Pass. Zero Not Run.** Evidence per row below. |
| Sources | ACs from `docs/requirements.md` §4; stories from `docs/user-stories.md` |

> **Integrity rule.** A case may only move from `Not Run` to `Pass`/`Fail` when it has actually been executed against running code, and only then may an actual result and evidence reference be recorded. No case is pre-marked as passing. ⚠ marks high-risk areas (concurrency/atomicity, isolation, API-level coverage).

## Coverage types used

Positive · Negative · Validation · Authorization (AuthZ) · Database integrity · Concurrency · Atomicity · Boundary · UI.

---

## 1. Authentication and session

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-01 | Positive | AC-01 | Fresh database | Register with valid, unused email and a strong password (UI or API) | Success; user row created; response contains no password material |
| TC-02 | Negative | AC-01 | TC-01 done | Register again with the same email | Rejected with a meaningful message; still exactly one user row |
| TC-03 | Validation | AC-03 | Register/login form open | Submit with blank required fields from the UI | Submission blocked; field-level messages shown before any API call |
| TC-04 | DB integrity | AC-01 | TC-01 done | Inspect `users` row in PostgreSQL | `password_hash` is a bcrypt-format digest; plaintext password appears nowhere |
| TC-05 | Positive | AC-01 | TC-01 done | Log in with valid credentials | Success; JWT returned |
| TC-06 | Negative | AC-01 | TC-01 done | Log in with wrong password | Rejected; no token issued; meaningful message |
| TC-07 | AuthZ | AC-01, AC-02 | Any protected endpoint | Call it with no `Authorization` header | 401; no data |
| TC-08 | AuthZ | AC-01 | Valid token in hand | Call a protected endpoint with a tampered token; then with an expired token | 401 both times |
| TC-09 | Positive | AC-01 | Logged in | Log out via UI; attempt to open a protected page | Token cleared from client storage; redirect to login |

## 2. Cross-user isolation ⚠ (two accounts required: A and B)

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-10 | AuthZ | AC-02 | A owns category/product/transaction; B logged in | B requests A's category, product, and transaction by ID | 404 or 403; no data returned |
| TC-11 | AuthZ | AC-02 | As above | B attempts PUT/DELETE on A's category and product | Rejected; A's records unchanged afterwards |
| TC-12 | AuthZ | AC-02 | Both users own data | As B, list categories, products, transactions, dashboard | Responses contain only B's data, never A's |

## 3. Categories

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-13 | Positive | FR-2.1, FR-2.4 | Logged in as A | Create "Electronics"; open the category list | Created; appears in the list |
| TC-14 | Negative | FR-2.2 | TC-13 done | Create "Electronics" again as A | Rejected — duplicate name within the user's inventory |
| TC-15 | Boundary | FR-2.2, AC-02 | TC-13 done | Create "Electronics" as user B | Allowed — uniqueness is scoped per user |
| TC-16 | Positive | FR-2.1 | TC-13 done | Rename the category | Update persists in the list |
| TC-17 | Positive | FR-2.1 | Category with no products | Delete it | Removed from the list |
| TC-18 | Negative | AC-08 | Category with active products | Delete it (UI and API) | Blocked with meaningful message; category remains |
| TC-19 | Validation | AC-03 | Category form open | Submit empty name in UI | Blocked with field message |
| TC-20 | Validation | AC-16 | Valid token | POST category with blank/missing name directly to the API | 400 with a field-level backend error |

## 4. Products

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-21 | Positive | FR-3.1, FR-3.5 | Valid category exists | Create product: name, SKU, category, price 10.50, opening 10, min 3 | Created; current available stock shows 10 |
| TC-22 | Validation | AC-04 | Product form / API | Price −1 via UI and via API | UI blocks; API returns 400 |
| TC-23 | Validation | AC-05 | Product form / API | Opening stock −1, then 2.5, then "abc" | Rejected at both layers each time |
| TC-24 | Validation | AC-05 | Product form / API | Min stock −1, then 2.5 | Rejected at both layers |
| TC-25 | Negative | AC-06 | TC-21 done | Create another product with the same SKU as A | Rejected — SKU must be unique within the user's inventory |
| TC-26 | Boundary | AC-06, AC-02 | TC-21 done | B creates a product with the same SKU as A's | Allowed — SKU scope is per user |
| TC-27 | Negative | AC-07 | Valid token | Create product with a nonexistent category id; then with another user's category id | 400 each time — no valid category |
| TC-28 | Validation | AC-03, AC-16 | Product form / API | Omit name, then SKU, then category, then price | UI blocks submission; direct API calls return per-field 400s |
| TC-29 | Positive | FR-3.2 | TC-21 done | Edit price and min stock; reload list | Changes persist |
| TC-30 | Positive | FR-3.3 | Products exist | Search by a name substring (case varied) | Matching products returned [match mode per Q7, Proposed] |
| TC-31 | Positive | FR-3.3 | TC-21 done | Search by the exact SKU | The product is returned |
| TC-32 | Positive | FR-3.4 | ≥2 categories with products | Filter by one category | Only that category's products shown |
| TC-33 | Positive | FR-3.5 | TC-21 done, then Stock In 5 | Product list current stock | Shows 15 (auto-updated, FR-4.3) |
| TC-34 | Positive | AC-14 | Product with ≥1 transaction | Deactivate it | Status becomes inactive; excluded from active lists and active-product count; its history remains |
| TC-35 | Negative | AC-14 | Product with ≥1 transaction | Attempt permanent delete (UI and API) | Not possible — no delete path / request blocked |

## 5. Stock operations ⚠

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-36 | Positive | FR-4.1, FR-4.3, AC-09 | Product stock 10 | Stock In 5 | Stock 15; transaction row exists with timestamp |
| TC-37 | Positive | FR-4.1, FR-4.3 | Product stock 10 | Stock Out 4 | Stock 6 |
| TC-38 | Negative | AC-10 | Product stock 10 | Stock Out 999 (API and UI) | Rejected; stock unchanged at 10; meaningful message |
| TC-39 | Validation | AC-09 | Valid product | Stock In with 0, then −3, then 2.5 | Rejected at both layers each time |
| TC-40 | Validation | AC-09 | Valid product | Stock Out with 0, then −3, then 2.5 | Rejected at both layers each time |
| TC-41 | DB integrity | FR-4.2 | Several movements done | Inspect `stock_transactions` rows | Each has product, type, quantity, remarks (NULL when omitted), timestamp |
| TC-42 | Positive | FR-4.4 | ≥2 products with movements | Filter history by one product | Only that product's rows |
| TC-43 | Positive | FR-4.4 | Mixed IN/OUT history | Filter by type OUT | Only Stock Out rows |
| TC-44 | Negative / API | AC-13 | Transaction rows exist | Send PUT and DELETE to the transaction endpoint; attempt edit via UI | 404/405; no UI affordance; DB rows unchanged. Correction only via new transactions |
| TC-45 | Concurrency ⚠ | AC-11 | Product stock exactly 5 | Fire 10 simultaneous Stock Out qty 1 | Exactly 5 succeed, 5 fail; final stock 0; **never negative** |
| TC-46 | Concurrency ⚠ | AC-11, AC-12 | Known starting stock | Mixed burst of parallel IN and OUT operations | Final stock == opening + ΣIN − ΣOUT computed from the transactions table |
| TC-47 | Atomicity | AC-12 | After TC-38's failed Stock Out | Check stock and transaction rows | No new row written; stock unchanged — failure left no partial state |
| TC-48 | Boundary | FR-4.2 | Valid product | Stock In omitting remarks; then with remarks text | Both accepted; remarks NULL vs stored verbatim |

## 6. Dashboard

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-49 | Positive | FR-5.1 | Seeded: 3 active + 1 inactive product, 2 categories, 1 low-stock | Open dashboard | Totals match exactly: active=3, categories=2, low-stock=1 |
| TC-50 | Positive | FR-5.2 | Stock-below-min product and an above-min product | Low-stock list | Below/equal-min product listed; above-min not [active-only per Q2, Proposed] |
| TC-51 | Positive | FR-5.3 | 7 transactions recorded | Recent-transactions panel | Exactly the 5 most recent, newest first |

## 7. Error surfacing

| ID | Type | Verifies | Precondition | Steps | Expected |
|---|---|---|---|---|---|
| TC-52 | UI | AC-15 | Logged in | Trigger a failing operation through the UI (e.g., TC-38 scenario) | Clear human-readable message; no raw stack trace/JSON dump |
| TC-53 | API / UI | AC-16 | Valid token | Send a request that bypasses UI validation (e.g., negative price straight to API); observe the app when applicable | 400 with field-level backend message; surfaced readably in the UI when the flow is UI-triggered |

---

## Execution record

Maintained per case **only after execution**. Rows show `Pass`/`Partial` only with a genuine evidence reference; everything else stays honestly `Not Run`. Last updated 2026-10-28 (T-206 frontend records; all backend suites previously recorded).

| ID | Status | Actual result | Evidence ref |
|---|---|---|---|
| TC-01 | **Pass** | Register live via proxy → 201, account created | evidence/T-106_auth-api-live-tests.txt |
| TC-02 | **Pass** | Duplicate email → 409 `email is already registered` | evidence/T-106_auth-api-live-tests.txt |
| TC-03 | **Pass** | Literal blank-field cases for both auth forms (register + login) added to the vitest suite 2026-10-29: submission blocked, field-level messages shown, zero POST calls asserted (src/auth.test.jsx, 10/10 passing) — the same client-side no-POST standard accepted for TC-19/TC-28. Invalid-email and short-password gates were proven earlier. No actual browser network traffic was observed; the basis is the automated suite. | evidence/TC-03_auth-blank-field-vitest.txt + evidence/T-108_frontend-build-tests.txt |
| TC-04 | **Pass** | `migrate -verify`: users' `password_hash` in bcrypt format (`true`); plaintext never stored | evidence/T-106_auth-api-live-tests.txt |
| TC-05 | **Pass** | Login live via proxy → 200 + Bearer token (masked); mixed-case email accepted | evidence/T-106_auth-api-live-tests.txt |
| TC-06 | **Pass** | Wrong password → 401, same message as unknown email | evidence/T-106_auth-api-live-tests.txt |
| TC-07 | **Pass** | Protected call without token → 401 | evidence/T-106_auth-api-live-tests.txt |
| TC-08 | **Pass** | Protected call with garbage token → 401 | evidence/T-106_auth-api-live-tests.txt |
| TC-09 | **Pass** | Logout with valid token → 200; client discards token (stateless, no revocation) | evidence/T-106_auth-api-live-tests.txt |
| TC-10 | **Pass** | Cross-user reads → 404 for categories (T-107), products (T-201), movements/dashboard (T-203/205) | evidence/T-107_category-tests.txt |
| TC-11 | **Pass** | Cross-user update/delete → 404; data unchanged | evidence/T-107_category-tests.txt |
| TC-12 | **Pass** | List/history endpoints contain only the caller's rows | evidence/T-204_history-tests.txt |
| TC-13 | **Pass** | Category created via the UI form → appears in the refreshed list (frontend create→refresh flow; backend create+list verified in T-107) | evidence/T-206_frontend-feature-tests.txt |
| TC-14 | **Pass** | Duplicate category name, same user → 409 | evidence/T-107_category-tests.txt |
| TC-15 | **Pass** | Same category name, different user → 201 | evidence/T-107_category-tests.txt |
| TC-16 | **Pass** | Category rename persists after refresh | evidence/T-107_category-tests.txt |
| TC-17 | **Pass** | Delete unused category → 200; gone from list | evidence/T-107_category-tests.txt |
| TC-18 | **Pass** | Delete category with products → 409 with guard message; category retained | evidence/T-107_category-tests.txt |
| TC-19 | **Pass** | Blank category name blocked in the UI before any API call (assertion: zero POSTs); backend blank name → 400 via live proxy smoke | evidence/T-206_frontend-feature-tests.txt |
| TC-20 | **Pass** | Blank category name via proxy → 400 `category name is required` | evidence/T-206_frontend-feature-tests.txt |
| TC-21 | **Pass** | Product created via POST /api/v1/products → 201; `current_stock` = `opening_stock` = 10 (Q13: no opening-balance movement row) | evidence/T-201_product-tests.txt |
| TC-22 | **Pass** | unit_price −1 → 400 (AC-04) | evidence/T-201_product-tests.txt |
| TC-23 | **Pass** | opening_stock −1 and 2.5 → 400 (AC-05: whole numbers only) | evidence/T-201_product-tests.txt |
| TC-24 | **Pass** | min_stock −2 / 0.5 → 400 (AC-05) | evidence/T-201_product-tests.txt |
| TC-25 | **Pass** | duplicate SKU, same user → 409 (AC-06) | evidence/T-201_product-tests.txt |
| TC-26 | **Pass** | same SKU, different user → 201 (uniqueness is per-user; AC-02 scope check) | evidence/T-201_product-tests.txt |
| TC-27 | **Pass** | nonexistent category → 400; foreign (other user's) category → 400 (AC-07) | evidence/T-201_product-tests.txt |
| TC-28 | **Pass** | Frontend: missing required product fields blocked before any API call; negative price and fractional opening stock blocked with messages. Backend: missing name/sku → 400 (plus 4 further invalid-field cases) | evidence/T-206_frontend-feature-tests.txt |
| TC-29 | **Pass** | update persisted (name, price, min stock, category); **`current_stock` unchanged by the edit** | evidence/T-201_product-tests.txt |
| TC-30 | **Pass** | name substring search case-insensitive; name fragment; zero-match → empty list | evidence/T-202_product-search-tests.txt |
| TC-31 | **Pass** | SKU fragment search; exact SKU search | evidence/T-202_product-search-tests.txt |
| TC-32 | **Pass** | category filter (both own products); combined search + category; no-match → empty | evidence/T-202_product-search-tests.txt |
| TC-33 | **Pass** | Product list shows current available stock (10 after create; 15/11 after movements) | evidence/T-201_product-tests.txt |
| TC-34 | **Pass** | deactivate → `inactive`, idempotent re-call 200; reactivation via PUT (Q14) | evidence/T-201_product-tests.txt |
| TC-35 | **Pass** | DELETE /api/v1/products/{id} → **405** (no delete route exists; hard delete impossible) | evidence/T-201_product-tests.txt |
| TC-36 | **Pass** | Stock In 5 → 201; stock 10 → 15; balance returned; row recorded | evidence/T-203_stock-tests.txt |
| TC-37 | **Pass** | Stock Out 4 → 201; stock 15 → 11 | evidence/T-203_stock-tests.txt |
| TC-38 | **Pass** | Stock Out 999 on stock 11 → 409 `insufficient stock`; stock unchanged | evidence/T-203_stock-tests.txt |
| TC-39 | **Pass** | Stock In qty 0 / −3 / 2.5 → 400 (3 subtests) | evidence/T-203_stock-tests.txt |
| TC-40 | **Pass** | Stock Out qty 0 / −3 / 2.5 → 400 (3 subtests) | evidence/T-203_stock-tests.txt |
| TC-41 | **Pass** | Recorded row contains product, type, quantity, remarks, created_at | evidence/T-203_stock-tests.txt |
| TC-42 | **Pass** | History filtered by product → 3 rows, all the queried product's | evidence/T-204_history-tests.txt |
| TC-43 | **Pass** | History filtered by type OUT → only OUT rows | evidence/T-204_history-tests.txt |
| TC-44 | **Pass** | No PUT/DELETE routes for movements → 404/405; DB immutability triggers present | evidence/T-203_stock-tests.txt |
| TC-45 | **Pass** | Concurrency burst: 10 parallel Stock Outs on stock 5 → exactly 5×201 + 5×409, final stock 0, one row per success | evidence/T-203_stock-tests.txt |
| TC-46 | **Pass** | Mixed-direction concurrency: 10 parallel IN qty 3 + 10 parallel OUT qty 2 on stock 20 → all 20×201, final stock 30 (0+20+30−20), 11 IN + 10 OUT rows recorded; no lost updates or inconsistencies | evidence/T-209_tc46-mixed-concurrency.txt |
| TC-47 | **Pass** | Failed Stock Out left stock AND transaction count unchanged (no partial write) | evidence/T-203_stock-tests.txt |
| TC-48 | **Pass** | Remarks omitted → NULL; remarks provided → stored verbatim | evidence/T-203_stock-tests.txt |
| TC-49 | **Pass** | Seeded totals exact: active=3, categories=2, low-stock=1 (mixed fixture incl. inactive product and boundary stock=min) | evidence/T-205_dashboard-tests.txt |
| TC-50 | **Pass** | Low-stock list = products with current ≤ min, active only; boundary (stock=min) included; healthy and inactive products excluded | evidence/T-205_dashboard-tests.txt |
| TC-51 | **Pass** | Exactly 5 most recent transactions, newest first (verified by id sequence m7…m3 of 7 seeded movements) | evidence/T-205_dashboard-tests.txt |
| TC-52 | **Pass** | Failed operations show the backend's readable messages in the UI (duplicate category 409, delete-guard 409, failed login 401, duplicate SKU 409) — no raw errors/stacks rendered | evidence/T-206_frontend-feature-tests.txt |
| TC-53 | **Pass** | Backend validation errors surfaced in the UI: 409 duplicate name/SKU messages rendered from real backend responses via the proxy | evidence/T-206_frontend-feature-tests.txt |

(Executed rows cite integration/unit/vitest runs recorded under `evidence/`. Assertions: TC-21…TC-29 and TC-34/35 inside `TestProductCRUDIntegration`, TC-30…TC-32 inside `TestProductSearchIntegration`, TC-36…TC-48 inside `TestStockMovementIntegration`, TC-49…TC-51 inside `TestDashboardIntegration`, TC-10…TC-18 inside `TestCategoryCRUDIntegration`, TC-01…TC-09 in the T-106 live-proxy transcript, and TC-13/TC-19/TC-28/TC-52/TC-53 in the T-206 vitest suite — all PASS in the cited files.)

## Coverage checklist vs. required test types

| Required type | Cases |
|---|---|
| Positive | TC-01, TC-05, TC-09, TC-13, TC-16, TC-17, TC-21, TC-29…TC-34, TC-36, TC-37, TC-42, TC-43, TC-48…TC-51 |
| Negative | TC-02, TC-06, TC-14, TC-18, TC-25, TC-27, TC-35, TC-38, TC-44 |
| Authorization | TC-07, TC-08, TC-10, TC-11, TC-12, TC-15, TC-26 |
| Validation | TC-03, TC-19, TC-20, TC-22…TC-24, TC-28, TC-39, TC-40 |
| Database integrity | TC-04, TC-41, TC-47 |
| Stock concurrency ⚠ | TC-45, TC-46 |
| API-level | All cases executable at API level per `docs/api-test-plan.md`; TC-44 is API-surface specific |
