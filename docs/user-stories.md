# User Stories — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.0 (planning phase) |
| Date | 2026-10-27 |
| Source | `docs/requirements.md` (traces to `Chandana 1.docx`) |
| Status | All stories are planned; none implemented |

## Roles

| Role | Description | Source |
|---|---|---|
| Store Owner | A registered user who manages exactly their own inventory (categories, products, stock movements). | [Explicit] The assignment names a "store owner" and requires per-user inventory isolation. |

No other role exists in the assignment — no admin, staff, or manager role is defined. Adding one would be scope creep. [Explicit absence]

## Priority scheme

MoSCoW. Every story below traces to an explicit FR/AC, so all are **Must**. `Rank` gives the build order within Must. ⚠ marks stories carrying **high-risk** implementation/test areas (atomic stock updates, negative-stock prevention under concurrency, cross-user isolation, API-level testing).

## Story summary

| ID | Pri | Rank | Role | Story statement | Dependencies | Test cases |
|---|---|---|---|---|---|---|
| US-01 | Must | 1 | Store Owner | I want to register an account so that I can create my own private inventory. | — | TC-01, TC-02, TC-03, TC-04 |
| US-02 | Must | 2 | Store Owner | I want to log in with my credentials so that I can access my inventory securely. | US-01 | TC-05, TC-06, TC-07, TC-08 |
| US-03 | Must | 3 | Store Owner | I want to log out so that my session ends on the device I am using. | US-02 | TC-09 |
| US-04 | Must | 4 | Store Owner | I want my products, categories, and transactions to be visible and changeable only by me, so that my business data stays private. | US-02 | TC-10, TC-11, TC-12 ⚠ |
| US-05 | Must | 5 | Store Owner | I want to create, view, update, delete, and list product categories so that I can organize my products. | US-02 | TC-13, TC-16, TC-17, TC-19 |
| US-06 | Must | 6 | Store Owner | I want category names to be unique in my inventory so that my catalog stays unambiguous. | US-05 | TC-14, TC-15 |
| US-07 | Must | 7 | Store Owner | I want to be prevented from deleting a category that still has active products so that I cannot orphan my products. | US-05, US-08 (test data) | TC-18 |
| US-08 | Must | 8 | Store Owner | I want to add products with name, SKU, category, unit price, opening stock, and minimum stock level so that my catalog reflects reality. | US-05 | TC-21…TC-28 |
| US-09 | Must | 9 | Store Owner | I want to view and edit product details so that I can keep product information current. | US-08 | TC-29 |
| US-10 | Must | 10 | Store Owner | I want to search products by name or SKU and filter by category so that I can find products quickly. | US-08 | TC-30, TC-31, TC-32 |
| US-11 | Must | 11 | Store Owner | I want to see the current available stock on each product so that I always know what is on hand. | US-08 | TC-33 |
| US-12 | Must | 12 | Store Owner | I want to record Stock In movements so that received goods increase my available stock automatically. | US-08 | TC-36, TC-39, TC-48 ⚠ |
| US-13 | Must | 13 | Store Owner | I want to record Stock Out movements, blocked when quantity exceeds available stock, so that I never oversell stock I do not have. | US-12 | TC-37, TC-38, TC-40 ⚠ |
| US-14 | Must | 14 | Store Owner | I want the system to keep stock correct even when several operations happen at once, so that concurrent activity never corrupts my numbers. | US-12, US-13 | TC-45, TC-46, TC-47 ⚠ |
| US-15 | Must | 15 | Store Owner | I want a transaction history filterable by product and type so that I can audit every stock change. | US-12 | TC-41, TC-42, TC-43 |
| US-16 | Must | 16 | Store Owner | I want transaction records to be immutable so that my audit trail cannot be rewritten; corrections happen as new transactions. | US-12 | TC-44 |
| US-17 | Must | 17 | Store Owner | I want a dashboard with totals, low-stock products, and the five most recent transactions so that I can see my inventory health at a glance. | US-05, US-08, US-12 | TC-49, TC-50, TC-51 |
| US-18 | Must | 18 | Store Owner | I want to deactivate a product that has transaction history instead of deleting it, so that history stays intact without cluttering active stock. | US-12 | TC-34, TC-35 |
| US-19 | Must | 19 | Store Owner | I want clear, meaningful messages when something fails, so that I know what went wrong and what to do next. | cross-cutting | TC-52, TC-53 |

## Acceptance criteria per story

Story-level ACs restate the mapped requirement-level ACs from `docs/requirements.md` §4; the AC IDs in parentheses are the authoritative source.

- **US-01** — Valid registration succeeds; duplicate email rejected with a message (AC-01). Blank required fields are blocked in the UI before submission (AC-03). The stored password is a hash, never plaintext (AC-01).
- **US-02** — Valid credentials return a token and grant access (AC-01). Wrong credentials are rejected without a token. Protected endpoints reject missing (TC-07) or invalid/expired (TC-08) tokens.
- **US-03** — After logout the client token is cleared and protected views are inaccessible until re-login.
- **US-04 ⚠** — User B cannot read (TC-10) or modify (TC-11) any of user A's records by ID, and no list endpoint ever returns A's data to B (TC-12). Verified explicitly with two accounts (high-risk area: cross-user isolation / IDOR).
- **US-05** — Category create, view, update, delete work; categories display in a list (FR-2.1, FR-2.4). Empty name blocked in UI (AC-03).
- **US-06** — Duplicate name within the same user's inventory is rejected (FR-2.2); the same name is allowed for a different user (scope of uniqueness is per-user).
- **US-07** — Deleting a category with active products is blocked with a meaningful message (AC-08). See Q3 [Proposed] for the stricter "any product references" rule.
- **US-08** — Product creation accepts all six fields (FR-3.1); rejects negative price (AC-04), negative/non-integer opening and minimum stock (AC-05), duplicate SKU within the user's inventory (AC-06), and missing/invalid category (AC-07). Required-field gaps are blocked in the UI and by the API (AC-03, AC-16).
- **US-09** — Product edits persist and are reflected in lists (FR-3.2).
- **US-10** — Search by name (FR-3.3, match mode per Q7 [Proposed]) and SKU; category filter narrows the list (FR-3.4).
- **US-11** — Current available stock is displayed per product and updates after movements (FR-3.5).
- **US-12 ⚠** — Stock In with a positive integer quantity increases stock and writes a complete transaction record (FR-4.1, FR-4.2, FR-4.3, AC-09). Zero/negative/non-integer quantities are rejected. Optional remarks handled per TC-48. Stock change and transaction row commit atomically (AC-12).
- **US-13 ⚠** — Stock Out decreases stock (AC-09); requests exceeding available stock are rejected with stock left unchanged (AC-10).
- **US-14 ⚠** — Parallel stock operations never produce negative or inconsistent stock (AC-11); final stock always equals opening + ΣStockIn − ΣStockOut from the transaction table (AC-12). High-risk area: verified by deliberately concurrent requests (TC-45, TC-46).
- **US-15** — Every transaction row contains product, type, quantity, optional remarks, timestamp (FR-4.2); history filters by product (FR-4.4) and by type.
- **US-16** — No API allows editing or deleting an existing transaction (AC-13); corrections are made by recording new transactions.
- **US-17** — Dashboard shows total active products, total categories, low-stock count (FR-5.1); lists products with stock at or below minimum (FR-5.2); shows the five most recent transactions (FR-5.3).
- **US-18** — A product with history can be deactivated, not deleted (AC-14); its history remains intact; it drops out of active lists and active counts.
- **US-19** — Failed operations show meaningful frontend messages (AC-15); backend validation errors are surfaced readably even when the request bypassed the UI (AC-16).

## Notes

- Ranks 1–19 are the intended implementation order and align with the sprint tasks in `docs/sprint-plan.md`.
- Test case IDs refer to `docs/test-cases.md`; all are **Not Run** as of 2026-10-27.
- Open questions Q1–Q14 in `docs/requirements.md` §9 affect several stories (US-07 → Q3; US-08 → Q13/Q14; US-10 → Q7; US-17 → Q2) and are flagged there.
