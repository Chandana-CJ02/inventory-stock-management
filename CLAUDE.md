# CLAUDE.md — Operating Rules for Claude Code in This Repository

Inventory and Stock Management System — a 3-day assignment. The assignment document **`Chandana 1.docx`** (project root) is the source of truth for requirements; nothing may be claimed, added, or implemented beyond it.

## Project facts

- **Stack (mandatory, from the assignment):** React.js functional components + hooks · Go REST APIs · PostgreSQL with relationships and constraints · JWT auth + hashed passwords.
- **Deadline:** 3 days, 8–10 h/day (user-imposed; the assignment document itself states no deadline). Schedule lives in `docs/sprint-plan.md`.
- **Phase status:** Implementation complete (2026-10-29): Days 1–3 tasks T-101…T-307 done; T-308 (GitHub push + submission checklist) in progress pending the user's GitHub repository. All 53 test cases Pass (evidence in `evidence/`); 33 frontend vitest tests and the full Go suite pass; `go vet`/`gofmt` clean.

## Document map

| Path | Contents |
|---|---|
| `docs/requirements.md` | Explicit FRs/ACs (AC-01…AC-16), NFRs, open questions Q1–Q14 with [Proposed] defaults |
| `docs/user-stories.md` | US-01…US-19 with priorities, dependencies, TC mappings |
| `docs/sprint-plan.md` | Day 1–3 task board (T-1xx/T-2xx/T-3xx), estimates, risks, fallback |
| `docs/architecture.md` | Component/data/auth design; [Explicit] vs [Proposed] decision log (§10) |
| `docs/test-cases.md` | TC-01…TC-53 definitions + execution record (all Not Run) |
| `docs/api-test-plan.md` | Endpoint scenario matrix, fixtures, concurrency method |
| `docs/evidence-plan.md` | Evidence naming, inventory, capture workflow |
| `docs/traceability-matrix.md` | AC → FR → US → TC → task → evidence chain |
| `docs/known-issues.md` | Honest defect/environment log (4 recorded: ENV-001…003 + DEF-001 — all closed) |
| `prompts/raw-prompt-log.md` | Verbatim user-prompt history (append-only) |

## Operating rules

1. **Plan before implementation.** No code until the user approves the specific task. One approved task at a time.
2. **Inspect before changing.** Read relevant files (and this repo's docs) before edits. Never modify files blindly.
3. **Stay in scope.** Work only in this project directory. Avoid unrelated directories, destructive commands, and unnecessary broad scans. No feature not in `docs/requirements.md` — proposed behavior exists only via the labeled open questions there.
4. **No environment changes without approval.** Do not install packages, initialize frameworks, start servers, or alter system configuration unless the active approved task says so. Environment checks (T-101/T-102) are read-only.
5. **Validate before claiming.** After each task: build/run what the task's validation step defines, inspect the diff, and only then mark the sprint-plan task done. `git status`/`git diff` before any commit.
6. **Status honesty — absolute.** Never claim a file was created, a test passed, or a requirement verified unless actually verified in this session. Unexecuted tests stay `Not Run`. Unknown → say unknown.
7. **Traceability discipline.** Every code change that affects behavior updates: the relevant TC rows (after execution), `docs/traceability-matrix.md`, and `docs/sprint-plan.md` statuses. New defects go to `docs/known-issues.md`.
8. **Prompt log discipline.** After each new user prompt in a session, append it verbatim to `prompts/raw-prompt-log.md` per that file's method (metadata separate, no secrets, `[unavailable]` markers for gaps — never reconstruct).
9. **Evidence discipline.** Evidence is captured at execution time per `docs/evidence-plan.md`. No evidence, no `Pass`.
10. **Secrets.** Never commit credentials, JWT secrets, or connection strings. Secrets live in git-ignored env config.
11. **Proposed vs explicit.** When implementing a [Proposed] decision from `docs/requirements.md` §9 / `docs/architecture.md` §10, note it in the task output so the user can override.
