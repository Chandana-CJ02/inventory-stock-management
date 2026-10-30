# Known Issues — Inventory and Stock Management System

| Field | Value |
|---|---|
| Last updated | 2026-10-28 (T-104 executed — backend scaffold built and validated) |
| Current issue count | **3 recorded, 0 open** — ENV-001/ENV-002 closed (verified), ENV-003 closed (mitigated) |

## Purpose

An honest, running log of every defect, limitation, environment finding, and accepted gap discovered during the project. Issues are recorded **regardless of whether they get fixed**; a submission with a documented, understood issue list is more credible than a silently incomplete one. Anything discovered during T-101 (environment checks) onward lands here first.

## Rules

- Every defect found during validation (T-301…T-303) is logged, fixed or not.
- Status changes only when actually true; `Fixed` requires a re-run, `Closed` requires a re-run **with passing evidence**.
- `Won't fix` requires a written justification and, if it affects an AC, explicit user sign-off before submission.

## Issue register

| ID | Date | Severity | Component | Description | Linked test/task | Status | Resolution / notes |
|---|---|---|---|---|---|---|---|
| ENV-001 | 2026-10-27 | Critical | Environment / Go toolchain | Go is not installed (or not discoverable): not on PATH (`where go` → not found); absent from `C:\Program Files\Go`, `C:\Go`, and scoop apps. Blocks all backend tasks. | T-101 → blocked T-104, T-105, T-106 | **Closed (verified 2026-10-28)** | Resolved: Go installed by user. Verified 2026-10-28 by user PowerShell **and** in-session probe: `go version` → `go1.27.3 windows/amd64`; binary at `C:\Program Files\Go\bin\go.exe`; on PATH. No evidence file captured (interactive environment check — see Environment findings ▸ Recheck) |
| ENV-002 | 2026-10-27 | Critical | Environment / PostgreSQL | PostgreSQL is not installed (or not discoverable): `psql` not on PATH (`where psql` → not found); absent from `C:\Program Files\PostgreSQL\*` and scoop apps. Blocks database setup and all persistence tasks. | T-101 → blocked T-102, T-105, T-106 | **Closed (verified 2026-10-28)** | Resolved: PostgreSQL 17.11 installed by user. Verified 2026-10-28 by user PowerShell **and** in-session probe: service `postgresql-x64-17` Running; port 5432 LISTENING; `psql (PostgreSQL) 17.11` at `C:\Program Files\PostgreSQL\17\bin\psql.exe` (bin dir in user PATH). Connection to the default `postgres` DB verified by user PowerShell only. **App database, dev credentials, and env config NOT created — T-102 remainder** |
| ENV-003 | 2026-10-28 | Low | Environment / port conflict | TCP port 8080 is occupied on this machine by another local service (a Jenkins instance answered `GET /healthz` with HTTP 403 + login redirect); the Go dev server could not bind (`listen tcp :8080: bind: Only one usage of each socket address …`). | T-104 | **Closed (mitigated 2026-10-28)** | Dev port switched to **8081** via `APP_PORT` in local `.env` and `.env.example`; code default remains 8080 (configurable — this is exactly what env-based config supports). Database connectivity unaffected |
| DEF-001 | 2026-10-29 | Low | Backend test design | `TestVerifyRejectsBadTokens/tampered_signature` was flaky (~1/64 fresh `go test` invocations): the test flipped the token's **final** base64url character, which carries zero-padding bits — an A→B flip there leaves the decoded 32 signature bytes unchanged, so the "tampered" token still verified. `-count=N` in-process repeats share time-quantized tokens, masking the flake. | T-305 / `backend/internal/auth/auth_test.go` | **Closed (verified 2026-10-29)** | Fixed test-only (production auth code untouched): flip a **mid-segment** character (data bits only — the flip always changes the decoded signature bytes). Re-run evidence: pre-fix failure caught live (fresh invocation #9); post-fix 200 in-process runs + 80 + 40 fresh invocations + full `go test ./...` suite pass. Evidence: `evidence/DEF-001_flaky-jwt-fix.txt` |

## Severity definitions

| Severity | Meaning |
|---|---|
| Critical | Breaks an explicit AC or loses/corrupts data (e.g., negative stock, cross-user leak) — blocks submission |
| High | Feature unusable or wrong results, workaround absent |
| Medium | Feature degraded; workaround exists |
| Low | Cosmetic or edge-case inconvenience |

## Status values

`Open` → `In progress` → `Fixed (pending re-test)` → `Closed (verified)` / `Won't fix (justified)`.

## Environment findings

T-101 executed 2026-10-27 (read-only checks only — no installation, no configuration changes):

| Tool | Status | Version | Location / evidence |
|---|---|---|---|
| Node.js | Available | v23.11.0 | `C:\Program Files\nodejs\node.exe` (`node --version`) |
| npm | Available | 10.9.2 | `C:\Program Files\nodejs\npm.cmd` (`npm --version`) |
| Git | Available | 2.43.0.windows.1 | `C:\Program Files\Git\cmd\git.exe` (`git --version`) |
| curl | Available | 8.4.0 | `C:\Program Files\Git\mingw64\bin\curl.exe` (also `C:\Windows\System32\curl.exe`) |
| Go | **Missing** (ENV-001) | — | `go version` → command not found; `where go` → not found; no install in `C:\Program Files\Go`, `C:\Go`, scoop |
| PostgreSQL (psql) | **Missing** (ENV-002) | — | `psql --version` → command not found; `where psql` → not found; no install in `C:\Program Files\PostgreSQL\*`, scoop |

Check scope: PATH lookup, `--version` probes, standard installation directories, and scoop locations. A deliberately portable/non-standard installation elsewhere on disk cannot be fully ruled out by these read-only checks. PostgreSQL connectivity (T-102) was **not** attempted — it depends on T-101 and is blocked by ENV-002.

### Recheck — 2026-10-28 (ENV-001/ENV-002 verified resolved)

The 2026-10-27 record above stands as what was true at that time. On 2026-10-28 the user installed Go and PostgreSQL 17 and supplied PowerShell verification (PROMPT-005); every claim was independently re-probed in-session (read-only) **except** the database connection test, which was not re-run to avoid starting T-102.

| Item | Verified result | Verified by |
|---|---|---|
| Go | `go version` → `go1.27.3 windows/amd64`; binary `C:\Program Files\Go\bin\go.exe`; on PATH | User PowerShell **and** in-session probe, 2026-10-28 |
| psql client | `psql (PostgreSQL) 17.11`; executable `C:\Program Files\PostgreSQL\17\bin\psql.exe`; bin dir in user PATH | User PowerShell **and** in-session probe, 2026-10-28 |
| PostgreSQL server | Windows service `postgresql-x64-17` — **Running**; port 5432 **LISTENING** (IPv4 + IPv6) | User PowerShell **and** in-session probe, 2026-10-28 |
| Connectivity | `psql -U postgres -h localhost -p 5432 -d postgres` succeeded; `SELECT version();` → PostgreSQL 17.11 on 64-bit Windows | **User PowerShell only** (not re-run in-session — T-102 scope, left pending) |
| Node.js / npm / Git | v23.11.0 / 10.9.2 / 2.43.0.windows.1 — unchanged | User PowerShell 2026-10-28 (matches 2026-10-27 session checks) |

**Evidence note:** no evidence files exist for these environment checks. The verification basis is the quoted command output supplied by the user plus the in-session probes recorded in this section. Capture-as-you-go evidence (per `docs/evidence-plan.md`) applies to application testing from T-109 onward, not to these interactive environment checks.

**Still pending — deliberately not performed (excluded by task guardrails):** application database creation, dev credentials, and backend env configuration (T-102 remainder); schema/migrations (T-105); all application implementation. T-102 status was updated to **Not started (unblocked)** in `docs/sprint-plan.md`; no other task statuses changed.

### T-102 execution — 2026-10-28 (application database and role created)

Setup was performed manually by the user in their own interactive psql session (password entered only via `\password`; never shared or stored outside their machine). **Provenance is labeled per result — nothing is called Claude-observed unless it was:**

| Item | Result | Provenance |
|---|---|---|
| Pre-flight check: `inventory_db` / `inventory_app` absent | `(0 rows)` for both | User-reported |
| `CREATE ROLE inventory_app LOGIN` → `\password inventory_app` → `CREATE DATABASE inventory_db OWNER inventory_app` | `CREATE ROLE` / silent success / `CREATE DATABASE` | User-reported |
| Connection as `inventory_app` to `inventory_db`; `SELECT current_database(), current_user;` | `inventory_db \| inventory_app` | User-reported |
| Role privilege flags | `rolsuper=f, rolcreatedb=f, rolcreaterole=f, rolreplication=f, rolbypassrls=f, rolcanlogin=t` | User-reported |
| Claude-side read-only connection re-check (`psql -w … -U inventory_app -d inventory_db`) | **Blocked** — `fe_sendauth: no password supplied` (no `pgpass.conf` exists; `-w` refuses to prompt by design). An authentication-mechanism gap for automation, not a setup failure | Claude-observed |
| Config files `.gitignore` / `.env.example` / `.env` created; `.env.example` placeholders only; `.env` carries a clearly-marked local password placeholder (real secret exists only on the user's machine after they edit it) | Verified (contents inspected with the password line masked) | Claude-observed |
| Git-side `.env` ignore check (`git check-ignore`) | **Deferred to T-103** — no repository exists yet | Documented deferral |

No open issues introduced by T-102. Optional follow-up (not required for T-102 completion): the user may add a `pgpass.conf` line for `inventory_app` to enable Claude-side re-verification and future automation.

### T-103 execution — 2026-10-28 (repository initialized; deferred V5 Git-side ignore check completed)

All results below are **Claude-observed** (actual command outputs):

| Item | Result |
|---|---|
| Pre-check: existing repository | None — `git status` → `fatal: not a git repository`; no `.git` directory |
| Initialization | `git init -b main` → branch `main` |
| **Deferred V5 check (from T-102)** | **PASSED** — `git check-ignore -v .env` → `.gitignore:2:.env  .env` (exit 0); `git status --short --ignored` lists `!! .env`; `.env` never staged or tracked |
| Staged / committed files | 15 — reviewed before commit; staged-content secret scan: 0 matches |
| Initial commit | `bd3b5c0` — "Initialize repository: Phase 1 planning artifacts + T-101/T-102 environment setup" (full hash `bd3b5c0030ab2b686fcdc9e35404d1ba396d877b`, 15 files, 1697 insertions) |
| Note | CRLF conversion warnings at staging (`core.autocrlf=true`) — informational, expected on Windows; no action taken |

The T-102 deferral recorded above ("Git-side `.env` ignore check → T-103") is hereby **completed**. The post-commit documentation updates recording T-103 (this entry and the sprint-plan status) are intentionally uncommitted and will be included in the next approved task's commit.
