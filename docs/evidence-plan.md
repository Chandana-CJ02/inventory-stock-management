# Evidence Plan — Inventory and Stock Management System

| Field | Value |
|---|---|
| Version | 1.0 (planning phase) |
| Date | 2026-10-27 |
| Evidence collected so far | **None.** Nothing has been executed; no evidence exists. This plan defines how evidence will be captured and named during T-109, T-208, T-301…T-306. |

## 1. Principles

1. **Capture at execution time.** Evidence is taken while the test runs — never reconstructed afterwards.
2. **Never fabricate.** A test without real captured evidence stays `Not Run`/`Blocked`; no screenshot, transcript, or query output is ever created for something that did not happen.
3. **Every executed test cites its evidence.** `docs/test-cases.md` execution-record rows reference the file(s) here.
4. **No secrets in evidence.** Capture steps mask passwords/tokens where feasible (e.g., truncated JWTs); evidence never includes real credentials.

## 2. Storage layout and naming [Proposed]

```
evidence/
├── TC-01_register-success.png
├── TC-04_users-password-hash.png        # psql output
├── TC-10_cross-user-category-404.txt    # raw request/response
├── TC-45_concurrent-stock-out-terminal.txt
├── TC-45_stock-after-concurrency.png
└── ...
```

Convention: `evidence/TC-<id>_<short-description>.<png|txt>`; suffix `-pN` for multiple files per case. API captures are raw text (command + status + body); UI evidence is a full-page screenshot showing the relevant state.

## 3. Evidence inventory by acceptance criterion (planned artifacts)

| AC | Planned evidence | Source tasks |
|---|---|---|
| AC-01 | Registration/login screenshots; login API capture with JWT (truncated); psql output of `users.password_hash` showing a bcrypt digest (TC-04); 401 captures for missing/invalid token (TC-07/08) | T-106, T-109 |
| AC-02 | Raw captures of B's cross-user GET/PUT/DELETE on A's ids (404s) + B's list responses showing only B's data (TC-10…TC-12) | T-302 |
| AC-03 | Screenshots of UI forms blocking empty/invalid submissions with field messages (TC-03, TC-19, TC-28 UI halves) | T-108, T-206 |
| AC-04…AC-07 | API captures of 400/409 rejections per case; DB constraint listings from migration output | T-201, T-301 |
| AC-08 | Capture of blocked category delete with message (TC-18) | T-301 |
| AC-09, AC-10 | Captures of rejected stock payloads (TC-38…TC-40) + unchanged stock query output | T-203, T-301 |
| AC-11, AC-12 | **Concurrency bundle** (§5): transcript, before/after stock queries, success/fail counts, ΣIN/ΣOUT reconciliation (TC-45…TC-47) | T-203, T-303 |
| AC-13 | Captures showing PUT/DELETE on transactions return 404/405 + unchanged rows (TC-44) | T-301 |
| AC-14 | Deactivation capture; product absent from active list; its history still present (TC-34/35) | T-201, T-301 |
| AC-15, AC-16 | Screenshot of a readable UI error (TC-52); raw 400 from a UI-bypassing request (TC-53) | T-304, T-301 |
| Dashboard FRs | Dashboard screenshot with the §3 fixture numbers visible (TC-49…TC-51) | T-205, T-208 |
| Deliverable D-2 | Migration files + their clean apply output on a fresh DB | T-105 |
| Deliverable D-3 | Final README; clean-clone run-through notes | T-307 |

## 4. Capture workflow

1. Execute the test per `docs/test-cases.md`.
2. Capture immediately to `evidence/` with the §2 naming.
3. Fill the case's row in the execution record: status (`Pass`/`Fail`), actual result, evidence file reference.
4. Update the matching `docs/traceability-matrix.md` row (`Tested` when its cases are executed; `Verified` only when passing with evidence).
5. On failure: open a `docs/known-issues.md` entry linking the failing evidence; re-test after the fix (T-305) and capture again — the failing capture is kept, not deleted.

## 5. Concurrency evidence checklist ⚠ (AC-11/AC-12 — high-risk)

A concurrency claim is only evidenced with **all** of:

- [ ] Before-state: `SELECT current_stock FROM products WHERE id=…` output (driven to the boundary value, e.g., 5)
- [ ] Full terminal transcript of the 10 parallel Stock Out requests (commands + responses)
- [ ] Success/failure counts from the transcript (expected 5/5)
- [ ] After-state stock query output (expected 0; never negative)
- [ ] Transaction-row count and ΣIN−ΣOUT reconciliation for the mixed batch (TC-46)
- [ ] Proof the DB row shows the atomic result (single `stock_transactions` insert set matching successes)

## 6. Gaps

If any executed test cannot produce its planned evidence (tool limitation), the gap is recorded in `docs/known-issues.md` with the substitution used — never silently skipped.
