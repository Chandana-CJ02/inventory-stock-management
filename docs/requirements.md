# Requirements Specification — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.0 (planning phase — nothing implemented or verified yet) |
| Date | 2026-10-27 |
| Source of truth | `Chandana 1.docx` (project root) |
| Related documents | `docs/user-stories.md`, `docs/architecture.md`, `docs/test-cases.md`, `docs/traceability-matrix.md` |

> **Convention.** Requirements quoted or closely paraphrased from the assignment are tagged **[Explicit]**. Every choice, interpretation, or default that is *not* in the assignment is tagged **[Proposed]** and consolidated in §9. No requirement on this page has been implemented or verified.

---

## 1. Problem statement [Explicit — verbatim from assignment]

> "Develop a web-based Inventory and Stock Management System for a small business. The application should allow a store owner to manage product categories, maintain products, record stock movements, and identify products running low on stock."

> "Expected outcome: A working inventory application demonstrating relational database design, authenticated CRUD APIs, data consistency, and business-rule validation."

## 2. Scope

**In scope [Explicit]:** user registration/authentication with per-user inventory isolation; category management; product management with search/filter; stock in/out with transaction history; low-stock dashboard; PostgreSQL persistence; JWT + hashed-password auth; React functional-component frontend; Go REST backend.

**Out of scope [Explicit absence — not mentioned in the assignment]:** team/organization accounts or role hierarchies (no admin role is defined); suppliers, purchase orders, invoicing, billing; reporting beyond the specified dashboard; deployment/hosting; mobile apps; notifications/emails; multi-currency. Adding any of these would be scope creep and will not be done.

## 3. Functional requirements

### FR-1 — User registration and authentication [Explicit]

| ID | Requirement |
|---|---|
| FR-1.1 | Provide registration, login, and logout functionalities. |
| FR-1.2 | Each registered user manages their own independent inventory. |
| FR-1.3 | Users must only access their own products, categories, and transactions. |

### FR-2 — Category management [Explicit]

| ID | Requirement |
|---|---|
| FR-2.1 | Users can create, view, update, and delete product categories. |
| FR-2.2 | Each category must have a unique name within the user's inventory. |
| FR-2.3 | Example names given: Electronics, Stationery, Accessories. (Examples only, not mandated values.) |
| FR-2.4 | Display categories in a list. |

### FR-3 — Product management [Explicit]

| ID | Requirement |
|---|---|
| FR-3.1 | Users can add products with: product name; SKU (unique product code); category; unit price; opening stock quantity; minimum stock level. |
| FR-3.2 | Users can view and edit product details. |
| FR-3.3 | Products can be searched by name or SKU. |
| FR-3.4 | Products can be filtered by category. |
| FR-3.5 | Display the current available stock for each product. |

### FR-4 — Stock management [Explicit]

| ID | Requirement |
|---|---|
| FR-4.1 | Two stock operations: Stock In (add quantities to available stock) and Stock Out (reduce quantities from available stock). |
| FR-4.2 | Each transaction must record the product, transaction type, quantity, optional remarks, and timestamp. |
| FR-4.3 | Available stock must update automatically after successful transactions. |
| FR-4.4 | Display stock transaction history with filtering by product and transaction type. |

### FR-5 — Dashboard [Explicit]

| ID | Requirement |
|---|---|
| FR-5.1 | Display total active products, total categories, and number of low-stock products. |
| FR-5.2 | Display products where available stock is at or below the minimum stock level. |
| FR-5.3 | Show the five most recent stock transactions. |

## 4. Acceptance criteria [Explicit — numbered AC-01…AC-16 for traceability]

| ID | Criterion (assignment wording) |
|---|---|
| AC-01 | Registration and login must use secure password storage and token-based authentication. |
| AC-02 | Users must not be able to access or modify another user's inventory. |
| AC-03 | Required fields must have frontend validation. |
| AC-04 | Product price cannot be negative. |
| AC-05 | Opening stock and minimum stock must be non-negative integers. |
| AC-06 | SKU must be unique within the user's inventory. |
| AC-07 | A product cannot be created without a valid category. |
| AC-08 | A category cannot be deleted while active products are associated with it. |
| AC-09 | Stock In and Stock Out quantities must be positive integers. |
| AC-10 | Stock Out must not be allowed when the requested quantity exceeds available stock. |
| AC-11 | The backend must prevent concurrent stock operations from creating incorrect or negative stock. |
| AC-12 | Stock transactions must be stored in PostgreSQL, and stock changes must be performed atomically. |
| AC-13 | Existing stock transactions must not be directly edited or deleted. Corrections can be made through new transactions. |
| AC-14 | Products can be deactivated rather than permanently deleted when transaction history exists. |
| AC-15 | Incorrect inputs and failed operations must display meaningful frontend messages. |
| AC-16 | All APIs must perform backend validation even when frontend validation is implemented. |

## 5. Non-functional requirements (derived from explicit ACs / constraints)

| ID | Requirement | Source |
|---|---|---|
| NFR-01 | Passwords stored only as secure hashes (never plaintext). | AC-01, auth constraint |
| NFR-02 | Protected resources require a valid token on every request. | AC-01 |
| NFR-03 | Every read and write is scoped to the authenticated user (defense against IDOR). | AC-02, FR-1.3 |
| NFR-04 | Every API endpoint validates its own input, independent of the frontend. | AC-16 |
| NFR-05 | Stock quantity changes and transaction inserts commit together or not at all. | AC-12 |
| NFR-06 | Concurrent operations can never yield negative or inconsistent stock. | AC-11 |
| NFR-07 | Stock transaction rows are immutable once written. | AC-13 |
| NFR-08 | Failed operations produce clear, human-readable messages in the UI. | AC-15 |
| NFR-09 | Relational integrity via foreign keys, uniqueness, and check constraints. | Mandatory table: "PostgreSQL with appropriate relationships and constraints" |

## 6. Technical constraints [Explicit — mandatory requirements table, verbatim]

| Area | Mandatory requirement |
|---|---|
| Frontend | React.js, functional components, hooks |
| Backend | Go with REST APIs |
| Database | PostgreSQL with appropriate relationships and constraints |
| Authentication | JWT-based authentication and hashed passwords |

## 7. Submission deliverables [Explicit]

| # | Deliverable |
|---|---|
| D-1 | GitHub URL for source code containing the React frontend and Go backend. |
| D-2 | PostgreSQL schema or database migrations. |
| D-3 | README with installation and execution instructions. |

## 8. Classification summary

- FR-1…FR-5, AC-01…AC-16, §6 constraints, §7 deliverables: **explicit** from the assignment.
- §5 NFRs: **derived** — they restate what explicit criteria demand; no new behavior introduced.
- Everything else in this project's docs (frameworks, ports, hashing algorithm, token lifetime, search match mode, status-field mechanics, etc.): **proposed** — see §9 and `docs/architecture.md` §10.

## 9. Open questions and proposed working decisions

None of the items below are specified in the assignment. Each has a proposed default so planning is not blocked; all are labeled **[Proposed]** and remain overridable. No open question has been silently decided.

| # | Question | Proposed working decision [all Proposed] |
|---|---|---|
| Q1 | The assignment document states **no deadline**. The 3-day / 8–10 h-per-day limit comes from the user's instruction, not the docx. | Day 1 assumed to start 2026-10-28; confirm actual due date. |
| Q2 | Products imply an active/inactive status (dashboard says "total active products"; AC-14 says products "can be deactivated"), but status is not in the field list. Hard delete when no history? Reactivation? | `status` ∈ {active, inactive}, default active. Hard delete only when the product has zero transactions; otherwise deactivation only. Reactivation allowed via edit. |
| Q3 | AC-08 blocks deleting a category with **active** products. What if only inactive products reference it (their FK would dangle)? | Block deletion when **any** product references the category (stricter than the literal AC; compliant with it; simplest safe rule). |
| Q4 | Are SKU / category-name uniqueness checks case-sensitive? Trimmed? | Case-insensitive, whitespace-trimmed (functional unique index on normalized value). |
| Q5 | Logout semantics: client-side token discard vs. server-side invalidation? | Client discards token; short token lifetime limits exposure. No server blacklist (not required by any AC). |
| Q6 | Token lifetime / refresh strategy? | JWT expiry ~24 h; no refresh tokens (nothing in the assignment requires them). |
| Q7 | Search match mode: substring vs. exact; case sensitivity? | Case-insensitive substring match on name; exact (case-insensitive) match on SKU also supported. |
| Q8 | Pagination for product lists / transaction history? | Not mentioned in assignment → omitted; transaction history endpoint accepts an optional limit. Revisit only if needed. |
| Q9 | Price precision and currency? | `NUMERIC(10,2)`, non-negative; single implicit currency, display-only. |
| Q10 | Go HTTP router, migration tooling, React build tooling? | Go standard library `net/http` (or chi router), plain SQL migration files applied by a small migrate step, Vite + React. Final choice recorded at implementation. |
| Q11 | Where is the JWT stored on the client? | `localStorage` with XSS risk noted in architecture §9; httpOnly cookie listed as alternative. |
| Q12 | Ports / CORS / UI styling? | Dev ports 8080 (Go) / 5173 (Vite); CORS restricted to the frontend origin in dev; plain CSS, no UI framework. |
| Q13 | Does creating a product with an opening stock generate an initial Stock In transaction row? | No — opening stock initializes the product's stock counter; movements begin with user-recorded transactions. |
| Q14 | Can SKU/category be edited after creation; can a deactivated product be reactivated? | Yes to both; SKU uniqueness re-checked on edit. |

**These decisions need user confirmation or override before/parallel with implementation start** — none blocks writing code, but each changes observable behavior.
