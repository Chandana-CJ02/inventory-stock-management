# Raw Prompt Log

Verbatim, chronological record of the user prompts issued to Claude Code for this project.

## Recording method

1. **Append-only.** New prompts are added at the end. Existing entries are never edited, reworded, or reordered.
2. **Chronological.** Entries are numbered `PROMPT-NNN` in the order they were received.
3. **Metadata kept separate from text.** Each entry has a metadata table (index, date received, project phase, tool/model, notes) followed by the exact prompt text in a fenced block. No commentary is inserted inside the fenced block.
4. **Verbatim fidelity.** Text is transcribed from the session transcript with original markdown preserved. Byte-level details (trailing whitespace, line-wrap) are not guaranteed; where any portion of the original text would be unavailable, it is marked `[unavailable]` rather than reconstructed. For the entries below, the full text was available and nothing was reconstructed.
5. **No secrets.** Credentials, tokens, API keys, passwords, or connection strings are never recorded. If a future prompt contains one, it is replaced with `[REDACTED]` and the redaction is noted in that entry's metadata.
6. **Scope.** Only user-authored prompt text is logged. Assistant replies, tool output, and system-injected notices are not prompts; a system notice attached to a prompt is mentioned in metadata only, never merged into the prompt text.
7. **Order repair (2026-10-29).** Entry order was corrected on 2026-10-29: PROMPT-022 and PROMPT-026 were relocated to their chronological positions (text unchanged, no renumbering) and six deferred prompts were appended verbatim (029–034) from conversation records. Appended numbering follows repair order, not exact receipt order; each appended entry's metadata records its receipt context. No prompt text was altered.
8. **Duplicate removal (2026-10-29).** Two verified double-logs were deleted with user approval: a byte-identical second copy of PROMPT-038, and the PROMPT-035 "T-208 Day 2 integration checkpoint instruction" entry whose metadata and verbatim body were byte-identical to PROMPT-036 apart from its ID. The distinct PROMPT-035 (concise audit request) and PROMPT-036 entries are preserved untouched. Result: 38 entries, IDs 001–038 each appearing exactly once. No prompt text was altered or reconstructed. Prompts received after PROMPT-038 during the finalization phase are not logged (deferred backlog, disclosed here rather than reconstructed).

---

## PROMPT-001 — Phase 1: planning-only instruction

| Field | Value |
|---|---|
| Received | 2026-10-27 |
| Phase | Phase 1 — planning and requirements analysis |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Message arrived with a system-injected attribution reminder prepended. That notice is not user-authored prompt text and is intentionally not reproduced here. |

```text
# Inventory and Stock Management System — Phase 1: Planning Only

We are starting a new assignment in:

`E:\Inventory-Stock-Management`

The assignment document is `Chandana 1.docx` in the project root. Treat this document as the source of truth.

## Strict scope boundary

**This phase is planning and requirements analysis only. Do not implement application features or write application code.**

Do not create React components, Go handlers, API implementations, database migration files, or executable scripts. Do not install packages, initialize frameworks, start servers, or modify system configuration.

## Step 1: Inspect and understand

1. Read `Chandana 1.docx` completely. Use an available document-reading method and verify that the extracted text covers the entire document.
2. Inspect the current project directory before making any changes.
3. Identify all explicit functional requirements, technical requirements, constraints, acceptance criteria, and submission deliverables.
4. Distinguish requirements explicitly stated in the assignment from recommendations or assumptions.
5. Do not invent requirements or claim that the assignment specifies something it does not.
6. If any requirement is ambiguous, record it in an open-questions section instead of silently deciding.

## Step 2: Prepare the SDLC planning artifacts

After inspecting the assignment, propose the following documents:

- `docs/requirements.md`
- `docs/user-stories.md`
- `docs/sprint-plan.md`
- `docs/architecture.md`
- `docs/test-cases.md`
- `docs/known-issues.md`
- `docs/traceability-matrix.md`
- `docs/api-test-plan.md`
- `docs/evidence-plan.md`
- `prompts/raw-prompt-log.md`
- `CLAUDE.md`
- `README.md` outline only; do not write implementation instructions that have not been verified.

Before creating these files, present a concise proposed outline and explain how each artifact will be used. Do not silently expand the scope beyond the assignment.

## Step 3: Planning requirements

### User stories
Give each story a unique ID, priority, user role, story statement, acceptance criteria, dependencies, and mapped test case IDs.

### Sprint plan
Plan within the three-day deadline, with 8–10 available hours per day. Include task IDs, estimates, dependencies, deliverables, validation steps, risks, and honest status fields. Reserve time for integration, API testing, defect fixes, evidence, and GitHub submission. Clearly flag any work that may not fit the deadline.

### Architecture
Use the required stack from the assignment:
- Frontend: React.js functional components and hooks.
- Backend: Go REST APIs.
- Database: PostgreSQL.
- Authentication: JWT and securely hashed passwords.

Describe component responsibilities, data relationships, authentication and authorization, API boundaries, transaction safety, and relevant security risks. Label any proposed choices that are not explicitly specified in the assignment.

### Testing and evidence
Map every acceptance criterion to test cases. Include positive, negative, authorization, validation, database integrity, stock concurrency, and API-level testing where relevant. Define how actual results and evidence will be recorded. Do not mark unexecuted tests as passed.

### Raw prompt log
Specify a method for preserving prompts verbatim, in chronological order. Keep metadata separate from the exact prompt text. Never record secrets or credentials.

### Claude Code operating rules
Plan before implementation. Work on one approved task at a time. Inspect relevant files before changing them. Avoid unrelated directories, destructive commands, and unnecessary broad scans. Validate changes and inspect diffs before accepting a task.

## Step 4: Report before writing

First report:
1. Assignment summary.
2. Explicit requirements and deliverables.
3. Proposed document outlines.
4. Key risks and unresolved questions.
5. Proposed three-day plan.
6. Any limitations in reading the assignment.

Then STOP and wait for my approval before creating the planning documents or implementing anything.

Do not claim a file has been created, a test has passed, or a requirement has been verified unless you have actually verified it.
```

---

## PROMPT-002 — Phase 1 approval and planning-document creation instruction

| Field | Value |
|---|---|
| Received | 2026-10-27 |
| Phase | Phase 1 — planning-document creation (approved) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | None. |

```text
Your Phase 1 report is reviewed. Proceed with **planning-document creation only**.

Before writing:
1. Preserve my original Phase 1 prompt verbatim in `prompts/raw-prompt-log.md`, to the extent the exact original text is available. Do not reconstruct missing text and label any unavailable portion honestly.
2. Include environment/tool availability checks in Day 1, without installing packages or changing system configuration.
3. Treat atomic stock updates, prevention of negative stock under concurrent requests, cross-user data isolation, and API-level testing as explicit high-risk tasks with planned validation.
4. Reserve testing, evidence capture, integration, defect fixes, README completion, and GitHub submission time across the three days.
5. Keep every assumption and proposed technical decision clearly distinguished from the assignment's explicit requirements.
6. Keep test results unexecuted until they are actually run. Do not claim that evidence exists unless verified.
7. Keep the scope limited to the assignment.

Create only the previously proposed planning artifacts:
- `docs/requirements.md`
- `docs/user-stories.md`
- `docs/sprint-plan.md`
- `docs/architecture.md`
- `docs/test-cases.md`
- `docs/known-issues.md`
- `docs/traceability-matrix.md`
- `docs/api-test-plan.md`
- `docs/evidence-plan.md`
- `prompts/raw-prompt-log.md`
- `CLAUDE.md`
- `README.md` outline only

Do not create application code, initialize frameworks, install packages, start servers, or modify system configuration.

After creating the documents, inspect the changes and report the files created, any gaps or unresolved questions, and the validation performed. Then STOP. Wait for my approval before beginning any implementation task.
```

---

## PROMPT-003 — T-101 approval (environment checks only)

| Field | Value |
|---|---|
| Received | 2026-10-27 |
| Phase | Day 1 implementation — T-101 only |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively on 2026-10-28: the turn's own guardrails prohibited any file edits, so the log could not be updated at receipt time. |

```text
I approve T-101 only.

Before proceeding, verify the current directory is `E:\Inventory-Stock-Management` and read the T-101 instructions in `docs/sprint-plan.md` and the applicable rules in `CLAUDE.md`.

Perform read-only availability and version checks for Node.js/npm, Go, PostgreSQL, and Git. Do not install packages, initialize a repository, change environment variables, modify configuration, or start servers.

Report each tool's availability, version where available, exact commands and results, and any blockers. Do not begin T-102 or any other task. Stop after T-101 and wait for my approval.
```

---

## PROMPT-004 — Read-only environment recheck (T-101 follow-up)

| Field | Value |
|---|---|
| Received | 2026-10-27/28 (session boundary) |
| Phase | T-101 follow-up — recheck only |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively on 2026-10-28 — same guardrail as PROMPT-003 (no project file edits permitted on that turn). |

```text
Continue the Inventory and Stock Management assignment.

Perform a read-only environment recheck for T-101 follow-up only.

1. Confirm the working directory is `E:\Inventory-Stock-Management`.
2. Check Node.js, npm, Go, Git, PostgreSQL client (`psql`), and whether a PostgreSQL server is available.
3. Check common installation locations and PATH without scanning the entire disk.
4. Report the exact commands, versions, and results. Distinguish "command not found" from "server not available" and "not installed."
5. Do not install or uninstall anything, modify environment variables or configuration, initialize Git, start servers, or edit project files.
6. Do not begin T-102 or any implementation task.

Stop and wait for my approval after reporting the results.
```

---

## PROMPT-005 — Environment-readiness documentation reconciliation

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Documentation/status reconciliation only (no implementation) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Message arrived with the project's CLAUDE.md injected as a system reminder; that injection is not user-authored prompt text and is not reproduced here. No secrets present. |

```text
We are working on the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

**Current scope: Environment-readiness documentation update only. Do not begin implementation.**

Before making any changes:

1. Read `CLAUDE.md` completely.
2. Read the relevant environment-setup and Task T-101/T-102 sections of `docs/sprint-plan.md`.
3. Read the relevant environment blocker entries in `docs/known-issues.md`.
4. Inspect the existing evidence and project status relevant to these checks. Do not perform broad directory scans or inspect unrelated projects.
5. Follow all project instructions and the established SDLC workflow.

**Verified environment results from PowerShell:**

- Node.js: `v23.11.0`
- npm: `10.9.2`
- Git: `2.43.0.windows.1`
- Go: `go1.27.3 windows/amd64`
- PostgreSQL client: `psql (PostgreSQL) 17.11`
- PostgreSQL Windows service: `postgresql-x64-17` — Running
- PostgreSQL connection test succeeded using `psql -U postgres -h localhost -p 5432 -d postgres`.
- `SELECT version();` returned PostgreSQL 17.11 on 64-bit Windows.
- The PostgreSQL executable is at `C:\Program Files\PostgreSQL\17\bin\psql.exe`.
- The PostgreSQL `bin` directory is saved in the user PATH. I refreshed PATH in the current terminal and confirmed `psql --version` works.

**Required actions:**

1. Compare these verified results with the existing T-101/T-102 statuses and ENV-001/ENV-002 entries.
2. Identify which documented blockers can legitimately be marked resolved based on the evidence above.
3. Update only the relevant environment-readiness sections in `docs/sprint-plan.md` and `docs/known-issues.md`, if justified.
4. Record the verification accurately. Do not claim that evidence files exist unless you inspect and confirm them.
5. Do not mark database schema creation, application configuration, API tests, frontend tests, or application implementation as complete. These have not been verified.
6. If T-102 includes additional checks beyond the results supplied above, leave those checks pending and explain what remains.
7. Preserve all unrelated planning artifacts and existing project decisions.

**Strict guardrails:**

- This is a documentation/status-reconciliation task only.
- Do not create or modify application source code.
- Do not install packages or tools.
- Do not create a database for the application or change PostgreSQL configuration.
- Do not initialize Git, commit, push, publish, or delete anything.
- Do not modify the original assignment document.
- Do not access or modify unrelated projects or previous assignment folders.
- Do not invent test results or reconstruct missing prompt history.
- Do not proceed to another sprint task automatically.

**Validation and stop condition:**

After the proposed documentation edits, show:
- Exact files changed.
- A concise summary of each change.
- The complete diff.
- Any remaining environment blockers or unverified checks.
- Validation results, if applicable.

Stop and wait for my approval. Do not make any further changes after presenting the diff.
```

---

## PROMPT-006 — T-101 validation-step wording correction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-101 documentation wording correction only (T-102 not started) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Message arrived with the project's CLAUDE.md injected as a system reminder; that injection is not user-authored prompt text and is not reproduced here. No secrets present. |

```text
We are working on the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

**Approved scope: T-101 documentation wording correction only. Do not start T-102.**

Before editing:
1. Verify the working directory is `E:\Inventory-Stock-Management`.
2. Read `CLAUDE.md` and the relevant T-101 section of `docs/sprint-plan.md`.
3. Read the prompt-log conventions in `prompts/raw-prompt-log.md`.

Make only these changes:

1. In `docs/sprint-plan.md`, update the T-101 validation-step wording from `Command outputs captured` to wording that accurately states that command outputs were checked and results were documented, but no separate evidence files were saved. Do not change the T-101 status, task, deliverable, dependencies, or any other row.
2. Append a new prompt-log entry documenting this exact request, following the existing metadata and formatting conventions. Preserve the exact prompt body without summarizing it. Do not rewrite existing entries.

Strict guardrails:
- Modify only `docs/sprint-plan.md` and `prompts/raw-prompt-log.md`.
- Do not create evidence files or claim that any exist.
- Do not change `docs/known-issues.md` or any other file.
- Do not run database setup, change configuration, install anything, or begin T-102.
- Do not initialize Git or run unrelated commands.

After editing:
1. Show the exact files changed.
2. Show the complete diff for both files.
3. Validate that only the approved wording and the new prompt-log entry changed.
4. Confirm that T-101 remains Done and T-102 remains Not started.
5. Stop and wait for my review. Do not make further changes.
```

---

## PROMPT-007 — T-102 read-only inspection instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 Phase A/B/C — read-only inspection and proposal |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28: the prompt's own guardrails forbade any project-file modification, so the log could not be appended at receipt time. Message arrived with the project's CLAUDE.md injected as a system reminder; that injection is not user-authored prompt text and is not reproduced here. No secrets present. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

**Approved scope: T-102 only — PostgreSQL connectivity and application database setup.**

T-101 is complete. Do not repeat T-101 or start any later task.

### Phase A — Read the instructions and inspect the current state

Before doing anything else:

1. Verify the working directory is exactly `E:\Inventory-Stock-Management`.
2. Read `CLAUDE.md` completely.
3. Read the complete T-102 task and its dependencies in `docs/sprint-plan.md`.
4. Read the relevant requirements, architecture, known issues, and database/environment configuration instructions.
5. Check `docs/test-cases.md` and `docs/evidence-plan.md` for the applicable validation and evidence requirements.
6. Read the relevant existing prompt-log conventions in `prompts/raw-prompt-log.md`.
7. Inspect only the files and settings relevant to T-102. Do not scan unrelated directories or access previous assignments.

### Phase B — Read-only PostgreSQL checks

Check the existing state without changing anything:

1. Confirm PostgreSQL is available and the expected service is running.
2. Check whether the application database already exists.
3. Verify database connectivity using the credentials already configured or explicitly supplied by the user. Do not print or log passwords.
4. Inspect the project's current environment/configuration files to determine whether database settings are already defined.
5. Check whether the application database name and development configuration are specified in the existing project documentation. Do not invent undocumented requirements without identifying them as proposals.
6. Distinguish confirmed facts, user-reported results, and unverified assumptions.

Do not expose passwords, tokens, or other secrets in terminal output, logs, evidence, or the repository.

### Phase C — Report and stop for approval

Before making any changes, report:

- The exact commands run and their results, with secrets redacted.
- Whether the PostgreSQL connection succeeded during this run.
- Whether the application database already exists.
- Whether development credentials and environment configuration already exist.
- Any missing information or blockers.
- The proposed database/configuration changes, if any.
- The exact files or database objects that would be created or modified.
- The applicable T-102 acceptance criteria and how they will be validated.
- Any required evidence and where it should be saved.

### Strict guardrails

- This authorization covers T-102 only.
- Do not create or drop databases, roles, or users during this initial inspection.
- Do not change passwords, environment variables, configuration files, or PostgreSQL settings.
- Do not install packages, start application servers, initialize Git, or modify application source code.
- Do not place credentials in source code, documentation, or the raw prompt log.
- Do not claim a test passed unless it was actually executed and its result observed.
- Do not mark T-102 Done merely because PostgreSQL is installed or the default `postgres` database is reachable.
- Do not modify project files during this initial inspection.

Stop after presenting the findings and proposed changes. Wait for my approval before making any database or configuration changes.
```

---

## PROMPT-008 — T-102 decisions and read-only check request

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 — decisions approved; authentication blocker resolution |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28 — the prompt forbade project-file changes at receipt time. No secrets present. |

```text
We are continuing the Inventory and Stock Management System assignment at `E:\Inventory-Stock-Management`.

I reviewed the T-102 read-only report. Proceed with T-102 only, using the following decisions:

- Proposed application database: `inventory_db`
- Proposed dedicated application login role: `inventory_app`
- Keep all credentials out of chat, source files, documentation, logs, and Git.
- Use a local `.env` for actual development settings, a placeholder-only `.env.example`, and `.gitignore` to exclude `.env`.

Before creating anything, resolve the authentication blocker safely:

1. Do not ask me to paste my PostgreSQL password into chat.
2. Tell me the exact command I should run in my own terminal or interactive `psql` session to check whether `inventory_db` and `inventory_app` already exist. I will execute it with my local authentication.
3. Do not create, drop, or alter any database or role until I return the check results and explicitly approve the setup commands.

After I return the results, prepare the minimum safe T-102 setup steps. Keep database creation and role creation separate, and use least-privilege access for the application role. Do not reset existing passwords or overwrite existing objects.

For now, make no project-file changes. Do not start T-103 or any later task. Preserve the raw prompt-log conventions. Once setup is explicitly approved and successfully validated, update only the T-102 documentation and prompt log required by the project instructions.

Stop after giving me the exact read-only check command and explaining where to run it.
```

---

## PROMPT-009 — T-102 setup preparation instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 — setup command preparation (no execution) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28 — the prompt forbade project-file changes at receipt time. No secrets present. |

```text
We are continuing the Inventory and Stock Management System assignment.

Project directory: `E:\Inventory-Stock-Management`

**Scope: T-102 only — PostgreSQL connectivity and application database setup.**

The read-only check was executed successfully in my PowerShell session. It returned `(0 rows)` for both proposed objects:
- Database: `inventory_db`
- Role: `inventory_app`

This confirms that neither object currently exists.

### Phase 1 — Prepare the setup commands only

1. Re-read the relevant T-102 instructions in `CLAUDE.md`, `docs/sprint-plan.md`, `docs/architecture.md`, `docs/requirements.md`, and `docs/known-issues.md` as needed.
2. Treat `inventory_db` and `inventory_app` as approved names for the proposed setup.
3. Prepare the minimum SQL commands to:
   - Create the dedicated `inventory_app` login role with a password entered securely by me, not embedded in the SQL text or shared in chat.
   - Create `inventory_db` as a separate operation, owned by `inventory_app`.
   - Grant only the permissions needed for the intended application development setup.
4. Explain which commands I should run in my own interactive `psql` session and which require reconnecting to the new database.
5. Do not reset passwords, overwrite existing objects, or grant superuser privileges.
6. Do not assume an existing role or database is absent without checking again immediately before creation.
7. Prepare the local configuration plan for `.env`, `.env.example`, and `.gitignore`. Keep real secrets only in the local `.env`; use placeholders in `.env.example` and ensure `.env` is ignored by Git.

### Phase 2 — Validation plan

Explain how we will verify:
- The database and role exist.
- The application role can connect to `inventory_db`.
- The connection identifies the expected database and user.
- The application role has no unnecessary administrative privileges.
- The local `.env` is excluded from Git.

Do not claim any validation has passed until it is actually executed and observed.

### Strict guardrails

- Do not execute SQL or create/modify database objects yet.
- Do not create or edit project files yet.
- Do not install packages, start application servers, initialize Git, or start T-103.
- Never ask me to paste a password into chat.
- Do not put passwords in command-line arguments, source code, documentation, evidence, or the raw prompt log.
- Preserve the exact raw-prompt-log conventions. Do not append the prompt log during this preparation phase.
- Do not mark T-102 Done yet.

### Stop condition

Show me the exact proposed commands, the secure password-entry procedure, the expected results, and the configuration files that would eventually be created or modified.

**Stop and wait for my explicit approval before executing any setup command or modifying any file.**
```

---

## PROMPT-010 — T-102 setup approval with safeguards

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 — setup approved (10 safeguards) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28 — the turn required user-side execution first; file edits were deferred by the prompt's own instructions. No secrets present. |

```text
I approve the proposed T-102 setup plan for `inventory_db` and `inventory_app`, subject to these safeguards.

1. Recheck that the database and role do not exist immediately before attempting creation. If either exists or the state is unexpected, stop and report; do not overwrite anything.
2. Create the `inventory_app` login role and `inventory_db` database in separate operations.
3. Set the role password securely through my own interactive PostgreSQL session. Never put the password in command-line arguments, scripts, project files, prompts, or logs. Do not ask me to paste it into chat.
4. Do not grant superuser, database-creation, role-creation, replication, or other unnecessary administrative privileges.
5. Configure database ownership and the minimum required access for application development. Explain any additional grants before applying them.
6. After the database is successfully created and verified, prepare `.env`, `.env.example`, and `.gitignore` as proposed. `.env` must be ignored by Git and contain local secrets only; `.env.example` must contain placeholders only. Do not initialize Git or start T-103.
7. Validate the database connection as `inventory_app`, verify the effective role privileges, and inspect the configuration files without printing secret values.
8. Do not mark T-102 Done until its applicable acceptance criteria have actually been validated. Record the Git-side ignore check as deferred to T-103 if no repository exists.
9. Update only the T-102-related entries in `docs/sprint-plan.md`, `docs/known-issues.md`, and `prompts/raw-prompt-log.md`, following the existing conventions, after validation. Preserve existing prompt-log entries.
10. Stop if a command fails or any unexpected state is found. Do not start T-103 or any other task.

Proceed with the approved T-102 setup only. Tell me when I must execute a command or enter a password myself. Never expose secrets. Report actual results, not assumptions.
```

---

## PROMPT-011 — Pause: complete SQL command listing request

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 — complete Action 1 command listing for review |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28 — the prompt forbade execution, files, and documentation updates at receipt time. No secrets present. |

```text
Pause execution. You have not provided the complete SQL commands for Action 1.

For T-102 only, show me the exact commands to run in my own interactive PostgreSQL session, in this order:

1. Recheck that `inventory_db` and `inventory_app` do not exist.
2. Create the `inventory_app` login role with a password entered securely through an interactive PostgreSQL mechanism. Do not put the password in SQL text, command-line arguments, scripts, logs, or chat.
3. Create `inventory_db` as a separate operation owned by `inventory_app`.

Explain when I must reconnect to another database and which commands require PostgreSQL administrator privileges.

Do not assume ownership automatically provides every schema permission the application will need; explain the PostgreSQL 17 defaults and any additional schema grants that may be required, without applying unapproved grants.

Do not ask me to create `pgpass.conf` yet. Do not execute commands, create files, change database objects, update documentation, or mark T-102 complete in this response.

Show the commands and stop for my review.
```

---

## PROMPT-012 — T-102 completion instruction (manual setup results reported)

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 — completion, config files, validation, documentation |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present — the user explicitly kept the role password out of the message. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

Scope: finish T-102 only. Do not start T-103 or implementation tasks.

I manually completed and observed these results:

1. Read-only object check returned `(0 rows)` for `inventory_db` and `inventory_app`.
2. `CREATE ROLE inventory_app LOGIN;` returned `CREATE ROLE`.
3. `\password inventory_app` completed interactively without an error. The password was not shared.
4. `CREATE DATABASE inventory_db OWNER inventory_app;` returned `CREATE DATABASE`.
5. I connected successfully using `psql -h localhost -p 5432 -U inventory_app -d inventory_db -W`.
6. `SELECT current_database(), current_user;` returned `inventory_db | inventory_app`.
7. The role privilege query returned `rolsuper=f`, `rolcreatedb=f`, `rolcreaterole=f`, `rolreplication=f`, `rolbypassrls=f`, `rolcanlogin=t`.

Proceed with the remaining T-102 work only:

1. Read the relevant T-102 instructions and existing file conventions before editing.
2. Check the database and role state with read-only queries before relying on it. Do not display secrets.
3. Prepare `.gitignore` so `.env` is excluded, `.env.example` with placeholder values only, and a local `.env` with the required connection setting names. Do not put real credentials in files created or edited by you. Leave a clearly identified local password placeholder for me to replace securely.
4. Ensure the Go configuration approach is consistent with the architecture documents; do not install packages or write application code.
5. Validate the configuration file contents without printing the local secret. If Git has not been initialized, record the Git-side ignore check as deferred to T-103 rather than claiming it passed.
6. Update only the T-102-related entries in `docs/sprint-plan.md`, `docs/known-issues.md`, and `prompts/raw-prompt-log.md`, following existing conventions and preserving previous prompt-log entries. Record the manual PostgreSQL results as user-reported, not Claude-observed.
7. Mark T-102 Done only if every applicable T-102 acceptance criterion has been validated; clearly record any deferred check.

Do not create/drop/alter database objects, change passwords, initialize Git, modify application source, or start T-103. Stop and report if unexpected state or validation failure occurs. At the end, list the exact files changed and the checks performed.
```

---

## PROMPT-013 — T-102 configuration verification instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-102 — post-setup configuration verification (read-only) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28 at T-103: the prompt's no-file-change scope forbade appending at receipt time. The system reported `.env` had changed on disk (the user's local password edit) — treated as current state, never displayed. No secrets present. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

**Scope: verify the completed T-102 configuration only. Do not start T-103.**

I have replaced the `PGPASSWORD` placeholder in `.env` locally and saved the file.

1. Read the relevant configuration conventions in `CLAUDE.md`, `docs/architecture.md`, and the T-102 section of `docs/sprint-plan.md`.
2. Validate `.env` without printing, echoing, logging, or otherwise exposing its password or full contents.
3. Check that the required PostgreSQL connection keys are present and non-empty, and that no placeholder remains. Report only pass/fail results and safe key names.
4. Validate `.env.example` separately and confirm it contains placeholders only, not actual credentials.
5. Inspect `.gitignore` and confirm its rules exclude `.env`. Do not claim Git-side ignore verification passed because Git initialization is deferred to T-103.
6. Do not connect to the database using a command that could expose the password. Do not change files, credentials, database objects, or documentation unless a validation failure requires stopping and reporting it.
7. Report the checks performed and any remaining limitations. Stop after T-102 verification.

Never print the value of `PGPASSWORD` or request that I share it.
```

---

## PROMPT-014 — T-103 Git initialization and initial commit instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-103 — Git repository initialization and initial commit |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present — the prompt itself states the local `.env` contains the PostgreSQL password, which must remain private. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

**PHASE: T-103 ONLY — Git repository initialization and initial commit.**

T-101 and T-102 are reported complete. T-102's Git-side `.env` ignore verification is deferred to this task. The local `.env` contains my PostgreSQL password, which must remain private.

### Before doing anything
1. Read `CLAUDE.md` and the exact T-103 section in `docs/sprint-plan.md`.
2. Read the relevant requirements, architecture, known-issues, evidence-plan, and raw prompt-log conventions.
3. Inspect the available Git-related skills and read their instructions if relevant.
4. Check the current directory and Git status. Determine whether a repository already exists before attempting initialization.
5. Inspect the existing `.gitignore` and confirm `.env` is excluded at the file level.
6. Do not print `.env` contents or expose any credentials.

### Execution rules
- Follow the T-103 acceptance criteria and existing project conventions exactly.
- Do not overwrite, delete, or reset existing work.
- Do not initialize a second repository if one already exists.
- Do not stage or commit `.env`, passwords, credentials, local-only files, or other secrets.
- Verify that `.env` is ignored by Git using an appropriate Git ignore check once the repository exists.
- Review the exact staged file list before committing. Stage only the files permitted by T-103.
- Create the initial commit only if the acceptance criteria are satisfied and the staged contents are safe.
- Do not modify application source code, install packages, create database objects, or begin T-104 or any other task.
- Preserve existing documentation and all previous raw prompt-log entries. Follow the established prompt-log procedure for recording this prompt and any subsequent approved file changes.
- Capture evidence only from actual command results; never fabricate evidence.
- If unexpected Git history, existing changes, secrets in the staged files, or another blocker is discovered, stop and report it before making a risky change.

### Final report
Report:
1. T-103 acceptance criteria and their individual pass/fail/deferred status.
2. Exact files staged and committed.
3. Commit hash and message, if a commit was actually created.
4. The actual result of the Git-side `.env` ignore check.
5. Evidence files created or updated.
6. Exact files changed and any remaining limitations.

Do not claim success without command output confirming it. Stop after T-103.
```

---

## PROMPT-015 — T-104 backend scaffold instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-104 — backend scaffold (first implementation task) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. The local `.env` (holding the real password) was edited only via an in-place `sed` line replacement and an append — never read or displayed. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

**Begin T-104 only: Backend scaffold. Our priority is fast, correct implementation without missing assignment requirements.**

### 1. Read only what's needed
- Read `CLAUDE.md`, the exact T-104 section in `docs/sprint-plan.md`, and the relevant parts of `docs/architecture.md`, `docs/requirements.md`, `docs/known-issues.md`, and `docs/test-cases.md`.
- Inspect the existing Go/backend skills and read the relevant instructions before using them.
- Inspect the current directory, Git status, and existing files before making changes. Preserve the two intentionally uncommitted documentation changes from T-103.
- Do not repeat the completed T-101–T-103 work or perform broad, unrelated scans.

### 2. Implement T-104 only
Follow the assignment and existing architecture. The intended scaffold includes:
- A Go module and clean backend directory structure.
- Environment-based configuration consistent with the architecture.
- PostgreSQL connection-pool setup using the existing local `.env` settings, without exposing credentials.
- A health endpoint that reports application health and database connectivity accurately.
- Appropriate configuration validation, timeouts, resource cleanup, and meaningful errors where required by the architecture.

Use the existing approved stack and conventions. Do not introduce unnecessary frameworks or dependencies. If an essential implementation choice is unspecified, choose the simplest suitable option, document the assumption, and continue unless it risks violating a requirement.

### 3. Protect secrets and existing work
- Never print, log, commit, or include `.env` values or passwords in output.
- Keep `.env` ignored by Git; never stage it.
- Do not discard or overwrite existing changes, reset Git, or force-push.
- Do not modify the assignment requirements to make implementation easier.
- Do not start T-105 or any later task.

### 4. Validate before declaring completion
- Run formatting and relevant Go checks/tests.
- Verify configuration failure behavior without revealing secrets.
- Test the health endpoint and database connectivity if feasible.
- Check the final Git diff and ensure no secrets or unrelated changes are included.
- Create or update evidence only from actual results, following `docs/evidence-plan.md`.
- Append this prompt verbatim to `prompts/raw-prompt-log.md` according to its existing conventions, preserving all previous entries.
- Update only the T-104-related sprint-plan and known-issues entries, and any other artifacts explicitly required by T-104.
- Mark T-104 Done only when its acceptance criteria are satisfied. Record failures or deferred checks honestly.

### 5. Work efficiently
Avoid repeated planning loops and unnecessary approval requests for harmless read-only checks. Make the smallest coherent implementation that satisfies T-104, validate it, and stop.

At the end, report the files changed, dependencies added, commands and tests run, actual results, acceptance criteria status, and any remaining blockers. Give me the next action only if a manual step is genuinely required.

**Stop after T-104. Do not begin T-105.**
```

---

## PROMPT-016 — T-105 database schema and migrations instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-105 — database schema and migrations |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. Database validation ran against `inventory_db` as `inventory_app` with credentials loaded in-process from the local `.env`; none were displayed. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

**Begin T-105 ONLY: Database schema and migrations. Work efficiently and satisfy the assignment requirements without skipping important constraints.**

### 1. Inspect the required context
- Read `CLAUDE.md` and the exact T-105 section in `docs/sprint-plan.md`.
- Read the relevant requirements, user stories, architecture, test cases, traceability matrix, and known issues.
- Read `docs/evidence-plan.md` and the existing raw prompt-log conventions.
- Inspect relevant database/migration skills if available.
- Inspect the current Git status and existing backend implementation before editing.

### 2. Preserve existing work
- Preserve all existing changes, including T-103 documentation changes and the uncommitted T-104 backend scaffold.
- Do not reset, discard, overwrite, or stage unrelated changes.
- Never display, log, commit, or expose `.env` credentials.
- Confirm `.env` remains excluded from Git.
- Do not commit anything unless the T-105 instructions explicitly require it.

### 3. Implement T-105 only
Follow the approved architecture and assignment requirements. Create the database schema and migrations needed for the inventory application, including the required entities and relationships for:
- Users and authentication-related data.
- Per-user product categories.
- Products and SKU uniqueness per user.
- Stock-in/out transaction history with timestamps and remarks.
- Any additional fields or constraints explicitly required by the assignment and approved architecture.

Ensure the design supports:
- Referential integrity and appropriate indexes.
- Nonnegative prices and valid quantity constraints where appropriate.
- User-specific data isolation.
- Preserving transaction history and avoiding hard deletion of products with history.
- Atomic stock updates and transaction recording in later backend operations.
- Safe, repeatable migration execution according to the chosen migration approach.

Do not invent extra business requirements. Resolve unspecified technical choices using the simplest approach consistent with the architecture, and document necessary assumptions.

### 4. Validate
- Run the migration/schema checks required by T-105.
- If applicable, test migrations against the local `inventory_db` using the `inventory_app` role.
- Do not drop or reset existing databases or data.
- Verify the schema and constraints using read-only queries after migration.
- Run relevant Go formatting, build, and tests if T-105 changes backend code.
- Capture evidence only from actual command output.
- Check the final diff and ensure no secrets or unrelated changes were introduced.

### 5. Update project records
- Append this prompt verbatim to `prompts/raw-prompt-log.md` following its conventions.
- Update only the T-105-related entries in the sprint plan and known-issues documentation, plus other artifacts explicitly required by T-105.
- Mark T-105 Done only when its acceptance criteria are actually validated. Clearly record failures and deferred checks.

**Scope restriction:** Do not begin T-106 or any later task. Stop after T-105 and report the exact files changed, migration files created, dependencies added, commands and tests run, actual results, acceptance criteria status, and any manual action genuinely required.
```

---

## PROMPT-017 — T-106 authentication APIs instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-106 — authentication APIs |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Logged retroactively 2026-10-28 at T-106 continuation: the live-test step was interrupted by the user before completion, and this entry was appended together with PROMPT-018 afterwards per the log's conventions. No secrets present. |

```text
We are continuing the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

Begin **T-106 ONLY: Authentication APIs**.

The project has completed T-101 through T-105 according to the previous task reports. T-105 added the database schema and migration runner. Existing changes are uncommitted, so preservation is critical.

## 1. Read and inspect before implementation

1. Read `CLAUDE.md` and the exact T-106 section in `docs/sprint-plan.md`.
2. Read the relevant sections of:
   - `docs/requirements.md`
   - `docs/user-stories.md`
   - `docs/architecture.md`
   - `docs/test-cases.md`
   - `docs/traceability-matrix.md`
   - `docs/known-issues.md`
   - `docs/evidence-plan.md`
3. Read the existing conventions in `prompts/raw-prompt-log.md`.
4. Inspect relevant installed Go authentication, JWT, password-hashing, and testing skills if available.
5. Inspect the current Git status and the existing backend, database schema, migration runner, configuration, and server routing. Reuse existing implementations rather than duplicating them.

Do not start coding until the requirements, approved design, and current implementation have been inspected.

## 2. Preserve existing work and credentials

- Preserve all existing changes from T-103, T-104, and T-105.
- Do not reset, discard, overwrite, or stage unrelated changes.
- Do not commit anything unless T-106 explicitly requires it.
- Never print, log, expose, or commit `.env` credentials, JWT secrets, database passwords, or other secrets.
- Do not read or modify unrelated projects.
- Do not reset, drop, or recreate the database.
- Keep manual command approvals enabled. Request approval for commands when required.

## 3. Implement T-106 only

Follow the exact approved T-106 scope and existing architecture. Implement the required authentication APIs, including the applicable registration, login, and logout behavior.

Requirements:

- Hash passwords securely using a well-established password-hashing implementation. Never store or return plaintext passwords.
- Validate request fields on the backend and return meaningful errors without leaking credentials or internal details.
- Use JWTs securely, with appropriate signing, expiration, and verification. Load secrets from configuration/environment; never hardcode a production secret.
- Reject missing, malformed, invalid, or expired tokens on protected routes as required by the architecture.
- Prevent duplicate accounts according to the approved schema and requirements.
- Use parameterized database operations.
- Keep database connections and server resources managed safely.
- Follow the existing API response and error-handling conventions.
- Do not claim that logout invalidates an already-issued stateless JWT unless the implemented design actually supports revocation. Follow the approved architecture and document any limitation.
- Do not implement category, product, stock transaction, or dashboard features as part of this task.
- Do not add dependencies unless necessary and compatible with the approved design. Explain why before adding any new dependency.

## 4. Validate the implementation

Run appropriate focused tests and the applicable existing Go checks, such as:

- `go test ./...`
- `go build ./...`
- `go vet ./...`
- `gofmt -l .`

Add tests for the authentication behavior and important failure cases required by T-106. Do not mark a test as passed unless it actually ran and passed.

If integration tests require database access, use the existing `inventory_db` and `inventory_app` configuration without exposing credentials or damaging existing data. Do not reset the database.

Inspect the final diff and Git status. Confirm `.env` remains ignored and is not tracked. Do not stage or commit changes.

## 5. Update project records

- Append this complete prompt verbatim to `prompts/raw-prompt-log.md`, following its existing ID and formatting conventions. Preserve all earlier entries.
- Update `docs/sprint-plan.md` and `docs/known-issues.md` based on actual results.
- Mark T-106 Done only if its required implementation and validations genuinely pass. Otherwise record the accurate partial or blocked status and explain why.
- Record evidence only if the evidence was actually captured and saved to the specified location.

## 6. Final report and stop condition

Report:
1. Files created or modified.
2. Authentication endpoints and behavior implemented.
3. Security measures and any limitations.
4. Exact validation commands and their actual results.
5. Any blockers or remaining known issues.
6. Git status and confirmation that prior changes and `.env` were preserved.
7. Evidence files actually created.
8. Updated estimated overall project completion percentage, clearly labeled as an estimate.

**Stop after T-106. Do not begin T-107 or any other task.**
```

---

## PROMPT-018 — T-106 continuation after interrupted live validation

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-106 — resumption after user-interrupted live API validation |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | The interrupted command was rejected before execution (no server started, no tests ran — verified, not assumed). No secrets present; JWT material was masked in all output and redacted from the evidence file. |

```text
Resume the current T-106 Authentication APIs task from the exact current state. Do not restart T-105, rewrite completed authentication code, or begin T-107.

The authentication implementation and Go checks were reported complete, with 55 test cases across three suites. The live API validation TC-01 through TC-09 was interrupted.

1. Inspect the current Git status and relevant diffs before proceeding. Preserve all existing T-103, T-104, T-105, and T-106 changes.
2. Check whether the server process started by the interrupted command is still running and whether port 8081 is occupied. Do not terminate an unrelated process.
3. Inspect the interrupted command and determine which live tests actually ran. Do not assume any test passed unless its result is available.
4. Resume only the remaining safe live authentication tests. Use the existing local database without resetting or dropping it. If a test account already exists, handle the duplicate-account case deliberately. Do not create unnecessary test accounts or delete existing user data.
5. Preserve useful test output as evidence under the project's prescribed evidence location. Never print `.env`, database credentials, `JWT_SECRET`, passwords, or a full JWT token. Mask tokens in output.
6. Ensure any server process started for testing is stopped safely, and remove only temporary files created by this test run.
7. Run any necessary final Go checks, inspect the final diff, and confirm `.env` is ignored and untracked.
8. Append the exact continuation prompt to the raw prompt log according to its existing convention, preserve all prior entries, and update the sprint plan and known issues based on actual evidence.
9. Mark T-106 Done only if the task's required checks genuinely pass. Report completed versus incomplete tests, changed files, evidence created, known limitations, and an estimated overall completion percentage.

Keep manual approvals enabled. Ask for approval when required. Stop after T-106.
```

---

## PROMPT-019 — T-107 category CRUD APIs instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-107 — category CRUD APIs with per-user scoping |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. DB integration tests ran against `inventory_db` as `inventory_app` with credentials loaded in-process from the local `.env`; none were displayed. |

```text
Continue the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

Begin **T-107 ONLY: Category CRUD APIs with per-user scoping**.

Before implementation:
1. Read `CLAUDE.md` and the exact T-107 section in `docs/sprint-plan.md`.
2. Read the relevant requirements, user stories, architecture, test cases, traceability matrix, known issues, evidence plan, and raw prompt-log conventions.
3. Inspect the existing Git status, database schema and migrations, authentication middleware, handlers, and routing. Reuse the existing T-106 implementation.
4. Preserve all uncommitted T-103 through T-106 work. Do not reset, discard, stage or commit changes. Never expose `.env` or secrets.
5. Implement only the category CRUD endpoints specified by the approved design. Every operation must use the authenticated user's ID from the verified JWT context—not a user ID supplied by the request.
6. Enforce per-user category-name uniqueness and prevent deletion when the category has active products, as required. Use parameterized SQL, backend validation, meaningful HTTP errors, and appropriate database constraints.
7. Add focused unit and integration tests for CRUD behavior, invalid input, unauthenticated requests, cross-user access, duplicate names, and deletion restrictions. Do not reset or recreate the database.
8. Run the required Go tests, build, vet and formatting checks. Preserve genuine evidence in the prescribed location without recording secrets.
9. Append this prompt verbatim to `prompts/raw-prompt-log.md`, preserve earlier entries, and update the sprint plan and known issues using actual results. Mark T-107 Done only if its required checks pass.
10. Report files changed, actual test results, evidence created, remaining issues, Git status, `.env` exclusion, and estimated overall completion.

Keep manual command approvals enabled. **Stop after T-107. Do not begin T-108 or any other task.**
```

---

## PROMPT-020 — T-108 frontend scaffold and authentication pages instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-108 — frontend scaffold and authentication pages |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present — the JWT_SECRET was generated server-side into `.env` and the frontend never receives it. |

```text
Continue the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

Begin **T-108 ONLY: Frontend scaffold and authentication pages**, following the exact approved scope in `docs/sprint-plan.md`.

1. Read `CLAUDE.md`, the T-108 sprint-plan section, requirements, user stories, architecture, test cases, traceability matrix, known issues, evidence plan, and raw prompt-log conventions.
2. Inspect current Git status and the existing backend authentication/category APIs, response formats, configuration, and frontend files. Preserve all uncommitted T-103 through T-107 work.
3. Implement only the T-108 scope using React functional components and hooks, consistent with the assignment's approved frontend architecture.
4. Implement registration, login and logout UI, appropriate form validation, loading/error states, and token handling as required by the approved design. Never expose database credentials or `JWT_SECRET` in frontend code or build-time configuration.
5. Integrate with the actual backend endpoints and response formats; do not invent APIs. Handle failed login, invalid input, expired/invalid tokens, and logout behavior appropriately. Do not claim server-side JWT revocation if the backend only discards the client token.
6. Follow the existing project conventions and use only necessary dependencies. Do not implement product, stock, dashboard, or other future tasks.
7. Add appropriate frontend tests where required, run the available lint/build/test checks, and record only genuine results. Do not reset the database, delete unrelated files, stage changes, or commit.
8. Append this prompt verbatim to `prompts/raw-prompt-log.md`, preserve earlier entries, and update the sprint plan and known issues according to actual results.
9. Confirm `.env` remains Git-ignored and untracked. Report files changed, actual checks and results, evidence created, limitations, Git status, and updated estimated project completion.
10. Mark T-108 Done only when its required validations pass.

Keep manual approvals enabled. Stop after T-108; do not begin T-109 or any other task.
```

---

## PROMPT-021 — T-109 Day-1 checkpoint instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-109 — integration smoke + evidence checkpoint |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. |

```text
Continue the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

Begin **T-109 ONLY**, following the exact approved scope in `docs/sprint-plan.md`.

1. Read `CLAUDE.md` and the T-109 sprint-plan section. Read any requirements, architecture, user stories, test cases, known issues, traceability, or evidence-plan sections relevant to T-109.
2. Inspect the current Git status and only the files needed for T-109. Preserve all uncommitted work from T-103 through T-108.
3. Check whether T-108's required validations and documentation updates are actually present. Do not assume the previous summary is proof; inspect the available files and evidence. Do not redo T-108 unnecessarily.
4. Before implementing, explain T-109's exact scope, prerequisites, affected files, and validation plan. If a prerequisite is incomplete or a scope conflict exists, stop and report it rather than silently expanding scope.
5. Implement only T-109. Follow the assignment architecture and project conventions. Do not start any later task.
6. Add or update tests as required, run appropriate checks, and record only genuine results. Do not reset the database, delete unrelated files, stage changes, or commit.
7. Append this exact prompt verbatim to `prompts/raw-prompt-log.md`, preserving all previous entries.
8. Update the sprint plan, known issues, and evidence only according to actual results. Do not mark T-109 Done unless its required validations pass.
9. Verify that `.env` remains Git-ignored and untracked without printing its contents or secrets.
10. Keep manual approvals enabled. Do not switch to auto mode.

Report the files changed, commands and actual results, evidence created, unresolved issues, Git status, and updated estimated overall completion. Then stop. Do not begin T-110 or any other task.
```

---

## PROMPT-022 — T-109 continuation after interrupted smoke attempt

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-109 — resumption after user-interrupted smoke attempt |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | The interrupted command was rejected before execution. An earlier malformed-pipeline attempt produced a 0-byte evidence file (proven to be this attempt's empty output; no earlier T-109 evidence existed). Test-password string appears only inside masked request payloads that are never printed; evidence scans report zero JWT/password/Authorization-header material. |
| Fidelity | The tail of this prompt beyond "confirm the evidence corresponds to" was not captured in the session transcript; the missing text is marked [unavailable] below rather than reconstructed. |

```text
Continue T-109 ONLY in `E:\Inventory-Stock-Management`. The previous smoke-test attempt was interrupted.

Before retrying:

1. Inspect the current state. Determine whether any backend or Vite processes are running and identify their PIDs and executable paths. Do not terminate a process unless you confirm it belongs to this task.
2. Inspect the relevant smoke-test command and evidence file. Determine whether the previous attempt created, truncated, or modified `evidence/T-109_integration-smoke.txt`. Preserve existing evidence unless it is proven to be incomplete output from this attempt; do not overwrite earlier evidence blindly.
3. Do not delete `backend/bin-server.exe` or any other project file. Do not reset the database or delete existing records.
4. Review the approved T-109 scope in `docs/sprint-plan.md`. This is a checkpoint task, not permission to implement new application features.
5. Run the integration smoke test in small, independently verifiable steps. Capture tokens only in shell variables; never print or persist tokens, passwords, Authorization headers, JWT secrets, or database credentials. Mask sensitive data before it reaches logs or evidence.
6. Use fresh, uniquely named test accounts and categories. Do not assume that test data is automatically cleaned up. Report any records left behind; do not delete them unless cleanup is explicitly authorized and can be restricted to records created by this run.
7. Validate expected HTTP status codes and response contents at each step. If a prerequisite fails, stop and report the exact failure instead of continuing and recording misleading evidence.
8. Capture only genuine results in `evidence/T-109_integration-smoke.txt`. Run a final review for secrets and confirm the evidence corresponds to [unavailable — message tail not captured in this transcript]
```

---

## PROMPT-023 — T-201 and T-202 product APIs instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-201 (product CRUD + deactivation) then T-202 (search/filter) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. DB integration tests ran against `inventory_db` as `inventory_app` with in-process `.env` credentials; none displayed. |

```text
Continue the Inventory and Stock Management System assignment in `E:\Inventory-Stock-Management`.

Proceed efficiently through **T-201, then T-202**, in that order. Do not skip validation gates or start any later task.

### 1. Read and verify the required context
1. Read `CLAUDE.md`.
2. Read the exact T-201 and T-202 sections in `docs/sprint-plan.md`.
3. Read the relevant requirements, user stories, architecture, test cases, known issues, traceability matrix, API test plan, and evidence plan.
4. Read the relevant existing backend code, migration/schema, T-107 category implementation, and current test/evidence files.
5. Inspect `git status` before changing anything. Preserve all existing uncommitted work. Do not stage or commit.

Confirm the exact T-201/T-202 scope, acceptance criteria, dependencies, and test cases from the documents. If a prerequisite is not met, stop and report it rather than guessing.

### 2. T-201 ONLY
Implement the sprint-plan-defined product CRUD APIs, backend validation, and product deactivation endpoint.

Requirements:
- Follow the existing architecture and coding conventions.
- Enforce authenticated, user-scoped access and ownership checks.
- Enforce the assignment's product validation and uniqueness rules.
- Respect foreign keys, stock transaction history, and the required deactivation behavior.
- Use parameterized SQL and meaningful HTTP errors.
- Do not implement T-202 search/filter yet.
- Do not implement stock movement, dashboard, or frontend product functionality unless the T-201 sprint-plan explicitly requires it.

Run the exact relevant tests and validations listed for T-201. Inspect the results and fix failures within T-201 scope. Record truthful results in the required evidence and project documents. Mark T-201 Done only if its completion criteria pass.

### 3. T-202 ONLY AFTER T-201 PASSES
Re-read the T-202 sprint-plan section and implement only its specified product search/filter behavior. Reuse the T-201 APIs and follow the documented requirements, query parameters, validation rules, and test cases.

Run the required T-202 tests, inspect the results, fix in-scope failures, and record truthful evidence. Mark T-202 Done only when its completion criteria pass.

### 4. Safety and execution rules
- Keep manual command approvals enabled. Never enable auto-approve or bypass permission prompts.
- Do not print, copy into logs, or expose secrets, tokens, or `.env` contents.
- Do not reset or drop the database, delete existing records, or perform broad cleanup.
- Do not modify unrelated files or previous assignments.
- Do not stage or commit changes.
- Append the exact prompts used to `prompts/raw-prompt-log.md` according to the project's established format. Do not fabricate missing prompt text.
- Follow existing documentation conventions and keep task statuses and evidence accurate.
- If a test fails, a schema assumption is unclear, or an operation has an ambiguous outcome, stop at that boundary and report the exact issue rather than proceeding blindly.

### 5. Final report
Report T-201 and T-202 separately:
- Status: Done, Blocked, or In Progress.
- Files changed.
- Acceptance criteria and test cases executed, with actual pass/fail results.
- Evidence files created or updated.
- Known issues and remaining manual tests.
- Updated overall completion estimate.
- The single next task from the sprint plan.

Do not begin T-203 or any later task.
```

---

## PROMPT-024 — T-201/T-202 documentation updates instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Documentation updates for T-201 and T-202 execution records |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. Preceded by a read-only verification checkpoint that confirmed the evidence files and test results being recorded here. |

```text
Continue the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

**Task: Documentation updates for T-201 and T-202 only. Do not start T-203.**

1. Read `CLAUDE.md` and the relevant sections of `docs/sprint-plan.md`, `docs/test-cases.md`, `docs/traceability-matrix.md`, and `docs/evidence-plan.md`.
2. Verify the existing T-201/T-202 evidence files and actual test results before recording any status.
3. Update `docs/test-cases.md` with only the test cases demonstrably supported by the evidence:
   - T-201: TC-21–TC-29, TC-34, TC-35.
   - T-202: TC-30–TC-32.
   - Record the actual outcome and evidence path for each case. Leave unsupported cases Not Run or appropriately marked.
4. Update `docs/traceability-matrix.md` only for acceptance criteria demonstrably covered by the implementation and tests. Do not mark unrelated or partially covered criteria fully Tested.
5. Do not mark TC-03 passed unless its existing evidence directly proves the required frontend behavior.
6. Follow existing documentation conventions. Do not fabricate results or evidence.
7. Preserve all existing uncommitted work. Keep manual command approvals enabled. Do not stage, commit, reset, delete records, or expose secrets.
8. Inspect the final diff and `git status`. Report the exact files changed and explain each status update.

**Stop after these documentation updates. Do not begin T-203 or any later task.**
```

---

## PROMPT-025 — T-201/T-202 documentation persistence verification

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Verification that the T-201/T-202 documentation updates persisted |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Prompted by an `Error editing file` on the first docs/test-cases.md edit attempt; the immediate corrected retry had succeeded. This verification confirmed the on-disk state directly. No secrets present. |

```text
Continue the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

The previous documentation update reported `Error editing file` for `docs/test-cases.md`. Do not assume the update succeeded.

1. Read `CLAUDE.md` and inspect the current contents and diff of `docs/test-cases.md` and `docs/traceability-matrix.md`.
2. Verify the actual T-201 and T-202 evidence files before making any test-status changes.
3. If the test execution records are missing, update `docs/test-cases.md` using the existing table format:
   - TC-21 through TC-29, TC-34, TC-35: record the outcomes only when supported by `evidence/T-201_product-tests.txt`.
   - TC-30 through TC-32: record outcomes only when supported by `evidence/T-202_product-search-tests.txt`.
   - TC-03 must remain Partial unless specific evidence proves the required frontend validation behavior.
   - Preserve all unrelated rows and existing notes.
4. Inspect the final diff and confirm that the intended changes were actually saved. If an edit fails again, stop and report the exact failure.
5. Do not modify code, database data, or unrelated documentation. Do not stage or commit. Keep manual approvals enabled.

Report the exact files changed and their verified final status. Stop without starting T-203.
```

---

## PROMPT-026 — T-203 and T-204 stock engine and history instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-203 (stock engine) then T-204 (history), with a validation gate between them |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. Concurrency evidence produced by real parallel HTTP requests against `httptest.NewServer` backed by `inventory_db`. |

```text
Continue the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

Work through **T-203 followed by T-204 in the same session**, but preserve the validation gate between them.

1. Read `CLAUDE.md`, the exact T-203 and T-204 sections in `docs/sprint-plan.md`, and all relevant requirements, architecture, test cases, known issues, API test plan, traceability matrix, and evidence plan.
2. Inspect `git status`, existing uncommitted work, the database schema, authentication, product APIs, and current tests/evidence. Preserve all existing work.
3. Implement **T-203 only**: stock-in/out APIs, input validation, user isolation, atomic stock updates, immutable transaction records, and concurrency safety as required by the sprint plan.
4. Run the required T-203 tests, including insufficient-stock, invalid-quantity, atomicity, and concurrency cases where specified. Inspect actual results and capture genuine evidence.
5. **Do not start T-204 unless T-203's required tests pass and its acceptance criteria are satisfied.** If blocked, stop and report the exact issue.
6. After T-203 passes, re-read the T-204 scope and implement only the specified stock transaction history functionality. Do not begin T-205 or later tasks.
7. Test T-204 against the real implementation and update required evidence and documentation accurately.
8. Append the exact prompt to `prompts/raw-prompt-log.md` according to project conventions. Do not fabricate missing prompt text.
9. Keep manual command approvals enabled. Never use auto-approve. Do not expose secrets, reset/drop the database, delete existing user data, stage, or commit. Use safe, scoped test fixtures and preserve existing uncommitted work.
10. Report each task separately: status, files changed, acceptance criteria, actual test results, evidence paths, known issues, pending manual tests, and updated completion estimate.

Stop after T-204. Do not start T-205 or any later task.
```

---

## PROMPT-027 — T-205 dashboard instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-205 — dashboard endpoint and validation |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. Dashboard integration tests ran against `inventory_db` as `inventory_app` with in-process `.env` credentials; none displayed. |

```text
Continue the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

**TASK: T-205 only — Dashboard. Do not start T-206 or any later task.**

### Phase 1: Read and inspect before implementation

1. Read `CLAUDE.md` completely.
2. Read the exact T-205 section in `docs/sprint-plan.md`, along with the relevant sections of:
   - `docs/requirements.md`
   - `docs/user-stories.md`
   - `docs/architecture.md`
   - `docs/test-cases.md`
   - `docs/known-issues.md`
   - `docs/traceability-matrix.md`
   - `docs/api-test-plan.md`
   - `docs/evidence-plan.md`
3. Inspect the current Git status and existing uncommitted changes. Preserve all existing work.
4. Inspect the existing database schema, authentication and user-isolation patterns, product APIs, stock movement APIs, transaction history API, frontend structure, tests, and relevant evidence.
5. Identify T-205's exact acceptance criteria, dependencies, expected dashboard behavior, and required test cases before making changes. Do not assume requirements that are not in the project documents.

### Phase 2: Implement T-205 only

Implement the dashboard strictly according to the T-205 requirements and acceptance criteria.

- Reuse existing APIs and components where appropriate.
- Preserve authentication and ensure every user's dashboard shows only that user's data.
- Follow existing project architecture, naming, validation, and error-handling conventions.
- Do not introduce unrelated features or refactor unrelated code.
- Do not weaken existing security, database constraints, tests, or user-isolation rules.

### Phase 3: Validate

1. Run the relevant automated tests and required project checks.
2. Test every T-205 acceptance criterion using safe, scoped fixtures and the project's approved test approach.
3. Verify dashboard counts, low-stock information, recent transactions, or other dashboard elements only to the extent required by the actual T-205 specification.
4. Check relevant edge cases, empty results, API errors, authentication, and user isolation where applicable.
5. Run the required build, formatting, vetting, and test checks applicable to the changed code.
6. Capture genuine evidence at the paths required by the evidence plan. Verify that each evidence file exists and accurately reflects the test results.
7. If any required test fails or an acceptance criterion remains unverified, investigate only within T-205's scope. If blocked, stop and report the blocker instead of claiming success.

### Phase 4: Documentation and prompt log

- Update the T-205 execution status, test-case results, traceability, known issues, and evidence references as required by the project conventions.
- Mark T-205 Done only when its required acceptance criteria and validation gate pass.
- Append this entire prompt verbatim to `prompts/raw-prompt-log.md`, following its existing entry format. Do not summarize or fabricate prompt text.

### Safety and scope guardrails

- Keep manual command approvals enabled. Never use auto-approve or bypass approval prompts.
- Do not expose passwords, JWTs, `.env` values, or other secrets in output, logs, or evidence.
- Do not reset or drop the database, run destructive migrations, or delete existing user data.
- Do not modify unrelated assignments or files.
- Do not stage or commit anything.
- Preserve all existing uncommitted changes.
- Do not claim that browser testing, API testing, or evidence capture occurred unless it actually did.
- Do not start T-206 or any later task.

### Final report

Report:
1. T-205 status: Done, Blocked, or Incomplete.
2. Acceptance criteria and their actual verification status.
3. Files changed and why.
4. Exact test commands and results, including failures.
5. Evidence files that exist and what they prove.
6. Known issues, unverified criteria, and pending manual tests.
7. Git status, confirming nothing was staged or committed.
8. Updated overall completion estimate, clearly labelled as an estimate.

**Stop after T-205 and wait for my approval before proceeding.**
```

---

## PROMPT-028 — T-206 frontend category/product pages instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-206 — frontend category and product pages |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. Frontend tests run via vitest/jsdom with mocked fetch; the live proxy smoke ran against the real backend. |

```text
Continue the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

**TASK: T-206 only. Do not start T-207 or any later task.**

### Phase 1 — Read and inspect first

1. Read `CLAUDE.md` completely.
2. Read the exact T-206 section in `docs/sprint-plan.md`.
3. Read the relevant requirements, user stories, architecture, test cases, known issues, traceability matrix, API test plan, and evidence plan.
4. Inspect current Git status and preserve every existing uncommitted change.
5. Inspect the existing React/Vite frontend, authentication flow, API proxy/client patterns, category and product APIs, stock movement API, transaction history API, dashboard API, frontend tests, and current evidence.
6. Identify T-206's exact acceptance criteria, dependencies, expected screens, and required tests before implementing anything. Follow the project documents as the source of truth; do not assume requirements.

### Phase 2 — Implement T-206 only

- Implement only the frontend features explicitly assigned to T-206.
- Reuse the existing React functional components, hooks, styles, and API patterns where appropriate.
- Preserve authentication and ensure users can only access their own data through the authenticated backend.
- Implement frontend validation, loading states, error handling, and success feedback as required by the documented acceptance criteria.
- Keep the UI consistent with the existing project and make it usable at the required screen sizes if specified.
- Do not implement T-207's dashboard rendering or any later task unless the T-206 documentation explicitly requires it.
- Do not perform unrelated refactors or overwrite existing work.

### Phase 3 — Test and validate

1. Run the relevant frontend unit/component tests and required build checks.
2. Test every T-206 acceptance criterion against actual results.
3. Add or update scoped tests for required validation, success, error, and loading behavior.
4. Run the relevant backend tests if changes could affect API integration.
5. Capture genuine evidence at the paths required by the evidence plan. Verify that the files exist and accurately reflect what was tested.
6. Clearly distinguish automated tests from manual browser tests. Do not claim browser testing unless it was actually performed.
7. If a required test fails or a criterion cannot be verified, fix it only within T-206's scope or report the blocker. Do not mark T-206 Done prematurely.

### Phase 4 — Documentation and prompt log

- Update the T-206 status, test-case execution records, traceability, known issues, and evidence references as required by project conventions.
- Append this entire prompt verbatim to `prompts/raw-prompt-log.md` in the established format. Do not fabricate or summarize the prompt.
- Mark T-206 Done only when its required acceptance criteria and validation gate pass.

### Safety and scope guardrails

- Keep manual command approvals enabled. Never use auto-approve or bypass permission prompts.
- Do not expose passwords, JWTs, `.env` values, or other secrets.
- Do not reset/drop the database, run destructive migrations, or delete existing user data.
- Preserve all uncommitted work.
- Do not stage or commit anything.
- Do not modify unrelated assignments or files.
- Do not start T-207 or any later task.

### Final report

Report:
1. T-206 status: Done, Blocked, or Incomplete.
2. Each acceptance criterion and its actual verification status.
3. Files changed and why.
4. Exact test commands and results, including failures.
5. Evidence files that exist and what they prove.
6. Known issues, unverified criteria, and pending manual tests.
7. Git status, confirming nothing was staged or committed.
8. Updated overall completion estimate, clearly labelled as an estimate.

**Stop after T-206 and wait for my approval before proceeding.**
```

---

## PROMPT-035 — Concise read-only prompt-log audit request

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Concise read-only audit of the prompt log (after the T-206 run) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Read-only turn; logging deferred (no file edits were authorized). This prompt requested the audit that identified the entry-ordering defect (027, 028, 026, 022 tail) and the missing T-207 run entry. No secrets present. |

```text
STRICT REQUIREMENT: ENGLISH ONLY. All communication, reasoning summaries, and final reports must be in English.

PROJECT: Inventory and Stock Management
WORKING DIRECTORY: E:\Inventory-Stock-Management

TASK: Verify and correct prompts/raw-prompt-log.md.

Before editing:
1. Read CLAUDE.md and the applicable prompt-logging instructions.
2. Inspect prompts/raw-prompt-log.md and verify the original source text for PROMPT-029 through PROMPT-034 from accessible original user messages or other authoritative, directly accessible records.
3. Do not reconstruct, paraphrase, or invent prompt text.
4. Report which prompts are VERIFIED, PARTIAL, or UNVERIFIED. Only complete, verbatim-verified prompts may be added.
5. Preserve existing prompt IDs and wording. Do not renumber prompts to hide historical ordering issues.

APPROVED CORRECTIONS:

A. Reorder existing entries:
- Move PROMPT-022 immediately after PROMPT-021.
- Move PROMPT-026 immediately after PROMPT-025.
- Preserve each entry's complete text and formatting exactly. Do not duplicate either entry.

B. Recover missing entries:
- Append PROMPT-029, PROMPT-030, PROMPT-031, PROMPT-032, PROMPT-033, and PROMPT-034 only when the complete original text of each is verifiably available.
- Preserve the original wording, formatting, lists, code blocks, paths, and constraints.
- If a prompt cannot be verified in full, do not add it. Report its ID and what source verification is missing.
- Do not create substitute or reconstructed prompts.

C. Disclosure:
- Add a concise disclosure line only if its claims are supported by verified facts.
- Clearly distinguish recovered prompt text from any text that remains unavailable.
- Do not claim chronological order is exact unless the original receipt order is verified.
- Do not invent dates, times, or provenance.

SAFETY AND SCOPE:
- Edit only prompts/raw-prompt-log.md.
- Do not modify any other file, including documentation, evidence, application code, or configuration.
- Do not read .env or expose secrets.
- Do not run tests/builds, access the database, make API requests, start/stop servers, install packages, or stage/commit/reset/discard Git changes.
- Preserve all existing uncommitted work.
- Use the smallest necessary changes.

AFTER EDITING:
1. Verify the final prompt IDs and order using read-only checks.
2. Confirm PROMPT-022 follows PROMPT-021 and PROMPT-026 follows PROMPT-025.
3. Report which of PROMPT-029–034 were added and which were omitted due to insufficient source evidence.
4. Report the exact file changed and a concise summary of the edits.
5. Show the relevant before/after ordering and any disclosure line added.
6. Confirm that no other files were modified by this task. Do not claim this unless you verify it.
7. Do not stage or commit changes.

Proceed with the approved, narrowly scoped corrections. If a required source cannot be verified, skip that prompt and explain why rather than guessing.
```

---

## PROMPT-036 — T-208 Day 2 integration checkpoint instruction

| Field | Value |
|---|---|
| Received | 2026-10-29 |
| Phase | T-208 — Day 2 integration checkpoint |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Received in the new session after the manual-verification checkpoint, server-start approval, and T-207 completion. No secrets present. |

```text
We are continuing the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

**CURRENT TASK: T-208 — Day 2 Integration Checkpoint**

GOAL: Complete T-208 according to the existing sprint plan, verify the end-to-end application, perform the two-user data-isolation spot check, and capture truthful evidence.

I am a beginner. Explain important findings in simple English. Work only on T-208 and its directly required evidence/documentation.

## PHASE 1 — READ AND PLAN

Before executing:
1. Read CLAUDE.md.
2. Read the complete T-208 section in docs/sprint-plan.md.
3. Read relevant acceptance criteria from docs/requirements.md and docs/user-stories.md.
4. Read applicable cases in docs/test-cases.md and docs/api-test-plan.md.
5. Read docs/traceability-matrix.md, docs/known-issues.md, docs/evidence-plan.md, and relevant sections of docs/architecture.md.
6. Inspect current Git status and relevant evidence for T-203 through T-207.
7. Inspect server/process/port readiness without sending requests or changing anything.
8. Identify the exact T-208 acceptance criteria, prerequisites, and test sequence.

Preserve all existing changes and untracked files. Never read .env or display secrets.

## PHASE 2 — CHECK SERVER READINESS

- If the required backend and frontend servers are already running, verify their processes and ports, then proceed to the permitted checks.
- If either server is stopped or unavailable, STOP and ask me for explicit approval before starting or restarting it. Do not start servers silently.
- Do not kill processes occupying ports. Remember that port 8080 was previously used by Jenkins and the application backend is configured for port 8081.
- Do not change configuration or install packages without asking.

## PHASE 3 — EXECUTE T-208

Once prerequisites are satisfied, execute the T-208 checks in the order required by the sprint plan.

A. End-to-end application flow
- Verify login and authenticated navigation.
- Verify category and product workflows.
- Verify stock-in and stock-out operations.
- Verify transaction history reflects successful movements.
- Verify dashboard totals, low-stock information, and recent transactions are consistent with the operations performed.
- Verify relevant validation and error handling.

B. Two-user data isolation
- Use two existing test accounts if suitable and safe, or create clearly identified test accounts only if permitted by the assignment and necessary.
- Verify that User A cannot access or modify User B's categories, products, stock transactions, or dashboard data.
- Test both direct API access and the UI where applicable.
- Do not expose passwords or tokens in logs or evidence.
- Do not delete existing user data or reset the database.
- Record exact outcomes and distinguish verified results from failures or untested cases.

C. Existing outstanding tests
- Check TC-03, the blank-field browser case, if it fits the checkpoint.
- TC-46, mixed concurrent stock-in/stock-out, is a separate outstanding test. Do not claim it passed unless it is actually executed and verified. If running it is outside T-208 or risks existing data, leave it pending and report why.
- Do not repeat already-passing tests unnecessarily; inspect their evidence first.

D. Evidence
- Save evidence only in the existing evidence/ directory, following docs/evidence-plan.md.
- Mask credentials, passwords, JWTs, cookies, and other secrets.
- Do not fabricate screenshots, logs, test results, or timestamps.
- Record failed tests and blockers honestly.

## PHASE 4 — UPDATE DOCUMENTATION

After testing:
- Update only documentation/evidence directly required to record T-208 results, using the existing project conventions.
- Update the relevant test-case results, traceability, known issues, evidence plan, and sprint-plan status as appropriate.
- Mark T-208 Done only if its acceptance criteria are met and the required evidence exists. Otherwise mark it Partial or Blocked and list the remaining work.
- Append my full original task prompt verbatim to prompts/raw-prompt-log.md according to CLAUDE.md. Do not rewrite existing prompt entries.
- Do not edit unrelated application code or unrelated documentation unless a verified defect directly blocks T-208 and I approve the fix.

## PHASE 5 — FINAL VERIFICATION

Use read-only checks to verify:
1. The exact files changed and current Git status.
2. Evidence files actually exist.
3. Each reported test result is supported by its output or evidence.
4. No files were staged or committed.
5. No existing data was deleted or reset.

Report:
- T-208 acceptance criteria and status.
- Tests passed, failed, blocked, and not run.
- Two-user isolation results.
- Evidence paths.
- Documentation updated.
- Outstanding TC-03 and TC-46 status.
- Exact remaining steps and estimated completion percentage, clearly labelled as an estimate.

SAFETY:
Never read .env, print secrets, reset/drop the database, delete existing data, kill unrelated processes, or stage/commit/reset/discard Git changes. Ask before starting or restarting servers, installing packages, changing configuration, or making out-of-scope fixes.

Proceed with the read/plan phase now. If servers are already running, continue with T-208. If starting a server is necessary, stop and ask for approval first. All communication must remain in English.
```

---

## PROMPT-029 — T-203/T-204 read-only validation checkpoint

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-203/T-204 — read-only validation checkpoint before T-205 |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Read-only turn; logging deferred (no file edits were authorized). Recovered verbatim on 2026-10-29 during the approved prompt-log repair. No secrets present. |

```text
Perform a read-only validation checkpoint for T-203 and T-204 only in `E:\Inventory-Stock-Management`.

Do not edit files, run migrations, modify database records, stage, or commit.

1. Read `CLAUDE.md`, the T-203/T-204 sprint-plan sections, and current `git status`.
2. Inspect the existing implementation and existing test/evidence files for both tasks.
3. Run the relevant non-destructive tests, using only the project's approved, scoped test fixtures.
4. Report exact results for T-203: valid IN/OUT, invalid quantities, insufficient stock, atomicity/rollback, user isolation, immutable transaction history, and concurrent stock-outs. Distinguish tests that were actually run from requirements not yet verified.
5. Report exact results for T-204: product/type filters, combined filters, ordering, limits, invalid filters, and user isolation.
6. Confirm which evidence files exist and what each actually proves.
7. Check whether `docs/test-cases.md` needs execution records. Do not edit files during this checkpoint; report the exact proposed documentation updates instead.
8. Do not assume that a prior completion report proves success.
```

---

## PROMPT-030 — Manual verification pre-check before T-207

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Read-only preflight before starting servers for manual browser verification |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Read-only turn; logging deferred (no file edits were authorized). Recovered verbatim on 2026-10-29 during the approved prompt-log repair. No secrets present. |

```text
We are doing a MANUAL VERIFICATION CHECKPOINT before T-207.

Do not implement T-207 or modify application code.

1. Read CLAUDE.md, the T-205 and T-206 sections of docs/sprint-plan.md, and the relevant test instructions in docs/test-cases.md and docs/evidence-plan.md.
2. Inspect the current Git status. Preserve all uncommitted changes.
3. Perform a read-only check of whether the backend and frontend servers are already running, which ports they use, and whether the existing health endpoint responds.
4. Identify the documented local URLs and the safe steps to open the existing application in a browser.
5. Do not start duplicate servers, install packages, change configuration, reset the database, create or delete records, or run any tests that modify data.
6. Do not display credentials, tokens, or .env values.
7. Report the exact commands and observed results. If a server is not running, explain the safest next step and wait for my approval before starting it.

Stop after this preflight check. Do not begin T-207.
```

---

## PROMPT-031 — Server-start approval for manual verification

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Approval to start backend and frontend servers for the manual browser session |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Turn only started servers and verified endpoints; no file edits, so logging deferred. Recovered verbatim on 2026-10-29 during the approved prompt-log repair. No secrets present. |

```text
I approve starting the existing backend and frontend servers for manual verification only.

Follow the startup sequence documented in CLAUDE.md and the project architecture. Keep manual command approvals enabled.

1. Start the backend first, using the documented command. Do not display .env contents or credentials.
2. Verify that the backend starts successfully and GET /healthz reports a healthy database connection. Do not create or modify database records.
3. Stop and report if the backend fails; do not change configuration to work around a failure.
4. Only after the backend is healthy, start the frontend in a separate terminal using the documented command.
5. Verify that the frontend responds at http://localhost:5173.
6. Do not run T-207, edit application code, install packages, reset the database, or stage or commit files.

Report the actual results and then stop so I can perform the browser tests manually.
```

---

## PROMPT-032 — T-207 frontend stock/history/dashboard instruction

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | T-207 — frontend stock, history, and dashboard pages |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Received after the servers were started for manual verification; not appended during the working turn (oversight, recovered in the 2026-10-29 prompt-log repair). No secrets present. |

```text
Continue the Inventory and Stock Management System in `E:\Inventory-Stock-Management`.

Proceed with **T-207 only** in `E:\Inventory-Stock-Management`.

Before implementation:
1. Read `CLAUDE.md` and the exact T-207 section in `docs/sprint-plan.md`.
2. Read the relevant requirements, user stories, architecture, test cases, known issues, traceability matrix, API test plan, and evidence plan.
3. Inspect current Git status and preserve all existing uncommitted changes.
4. Inspect the existing implementation and identify T-207's exact scope, acceptance criteria, dependencies, and validation gate.

Implement T-207 only. Do not start T-208 or any later task.

Keep manual command approvals enabled. Do not reset the database, delete existing records, expose secrets, stage or commit files, or make unrelated changes. Do not stop the running servers unless necessary and approved.

Run the required scoped tests and checks, capture genuine evidence, update documentation, and append this entire prompt verbatim to prompts/raw-prompt-log.md.

If a required check fails or a blocker prevents completion, report it honestly. Mark T-207 Done only when its documented acceptance criteria and validation gate pass.

Stop after T-207 and report the exact results. Wait for my approval before T-208.
```

---

## PROMPT-037 — Day 3 regression verification instruction

| Field | Value |
|---|---|
| Received | 2026-10-29 |
| Phase | Day 1 and Day 2 regression verification before Day 3 |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. |

```text
Inventory & Stock Management — Day 1 and Day 2 Regression Verification Before Day 3.

Bro, Day 2 was completed before the frontend UI redesign. Before starting Day 3, verify that all completed Day 1 and Day 2 functionality still works after the redesign.

**English only. This is a verification task first, not a new implementation task. Follow the existing project documents and preserve all existing work.**

## 1. Read the project instructions

Read:
- `CLAUDE.md`
- `docs/requirements.md`
- `docs/user-stories.md`
- `docs/sprint-plan.md`
- `docs/test-cases.md`
- `docs/known-issues.md`
- `docs/traceability-matrix.md`
- `docs/api-test-plan.md`
- `docs/evidence-plan.md`

Identify the exact Day 1 and Day 2 tasks and their acceptance criteria from the existing documents. Do not invent or change task statuses.

## 2. Inspect the current project safely

- Check Git status and inspect the current changes.
- Preserve all modified, untracked, and existing evidence files.
- Confirm the frontend and backend are available on their configured ports before testing.
- Never read `.env` or print secrets, passwords, or JWTs.
- Do not reset the database or delete existing records.
- Do not stage, commit, push, reset, or clean anything.

## 3. Verify completed Day 1 and Day 2 functionality

Use the project's existing tests and test plan. Verify as applicable:

### Environment and application
- Frontend starts and loads.
- Backend health endpoint responds and database connectivity works.
- Frontend API proxy/integration works.

### Authentication and user isolation
- Registration and login work.
- Invalid credentials and required-field validation behave correctly.
- Authenticated API access works.
- One user cannot access another user's categories, products, or transaction history.

### Categories
- Create, list, update, and delete categories according to the acceptance criteria.
- Duplicate category names and deletion restrictions are handled correctly.

### Products
- Create and list products.
- Update and deactivate products as supported.
- Validate SKU uniqueness, category selection, prices, quantities, and required fields.
- Search and filtering work.
- Preserve the existing product and category data.

### Stock and transaction history
- Stock In and Stock Out update stock correctly.
- Stock Out cannot reduce stock below zero.
- Invalid quantities are rejected.
- Transaction history records the correct product, movement type, quantity, remarks, and timestamp.
- Historical transactions remain consistent with stock movements.

### Dashboard and frontend redesign regression
- Dashboard metrics and recent transactions load correctly.
- Products, Categories, and Stock pages render correctly with the new shared styles.
- Existing forms, buttons, navigation, loading/error states, and validation messages remain usable.
- Check desktop and narrow viewport layouts where feasible.
- Verify that styling changes did not break existing behavior.

### Automated and API tests
- Run the complete existing frontend test suite.
- Run the production build.
- Run the existing backend tests.
- Run relevant API-level integration tests and existing user-isolation checks.
- Include the existing high-risk stock tests where feasible.

Do not rerun destructive or data-mutating scenarios against existing user data unless a safe, isolated test setup is confirmed. Use isolated test data and clean up only data created specifically for the test, if safe and authorized.

## 4. Report results honestly

Create a verification report only if supported by the project's existing documentation workflow. Do not modify unrelated documents.

For every applicable Day 1 and Day 2 requirement, report:
- Task or acceptance criterion.
- PASS, FAIL, BLOCKED, or NOT RUN.
- Exact test or verification performed.
- Actual result and evidence location, if available.
- Any remaining defect or limitation.

Distinguish automated test results from manual browser checks. Do not mark a requirement as passed merely because a previous report said it passed.

If a defect is found, explain it and its scope first. Do not begin unrelated fixes or modify backend behavior without authorization.

## 5. Day 3 readiness decision

At the end, give me:
1. A Day 1 and Day 2 verification summary.
2. Failed, blocked, and untested acceptance criteria.
3. Whether Day 3 can safely begin.
4. The highest-priority remaining tasks for Day 3.

**Do not start Day 3 implementation yet. Finish the regression verification and stop for my review.**

English only.
```

---

## PROMPT-038 — Day 3 Finalization instruction

| Field | Value |
|---|---|
| Received | 2026-10-29 |
| Phase | Day 3 — finalization (T-304…T-308, backend startup fix, prompt-log reconciliation) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | No secrets present. Backend found down and restarted per section 4 of this prompt. |

```text
Inventory & Stock Management — Day 3 Finalization

We have very limited time remaining. Finish all remaining required assignment work efficiently, preserve existing work, and prepare the project for submission.

## 1. Mandatory first steps
- Read `CLAUDE.md`, `docs/sprint-plan.md`, `docs/requirements.md`, `docs/test-cases.md`, `docs/api-test-plan.md`, `docs/traceability-matrix.md`, `docs/evidence-plan.md`, `docs/known-issues.md`, and `README.md`.
- Inspect `git status`, the current diffs, and the existing evidence inventory before changing anything.
- Preserve all existing implementation, user data, evidence, and unrelated modifications. Do not reset, clean, discard, or overwrite existing work.
- Do not read `.env` or print secrets.
- Do not stage, commit, or push without showing me the exact proposed actions and receiving my explicit approval.

## 2. Verify reported completion
- Audit T-301, T-302, and T-303 against actual test outputs, integration assertions, and evidence files. Do not rely on sprint-plan status text alone.
- Confirm TC-10, TC-11, TC-12 cross-user isolation, TC-45, TC-46, and TC-47 concurrency/atomicity results from actual evidence.
- Report every result as PASS, FAIL, BLOCKED, NOT RUN, or NOT VERIFIED.
- Keep the known flaky JWT `tampered_signature` test failure visible and documented. Do not fix it without approval or claim the entire backend suite passes if it does not.

## 3. Complete remaining tasks in priority order
### T-304 — Frontend validation and error polish
- Inspect existing UI and test coverage before making changes.
- Verify TC-52 and TC-53 using the safest available method.
- Fix only genuine, in-scope defects. Do not spend time on optional styling polish.
- Run relevant tests and the production build after any changes.

### T-305 — Defect buffer
- Fix only verified, important defects within the remaining time.
- Before any backend code change, explain the exact proposed change and wait for my approval.
- Do not modify the backend as part of this run without approval.
- Record unresolved issues honestly.

### T-306 — Evidence completion
- Inventory evidence already on disk and reuse valid existing results.
- Ensure test records reference evidence that actually supports the reported result.
- Do not fabricate evidence, claim unrun tests passed, or rerun tests unnecessarily.

### T-307 — README
- Replace relevant TODOs with verified setup, database/migration, backend, frontend, testing, and execution instructions.
- Inspect actual scripts, paths, ports, and commands before documenting them.
- Do not expose secrets or invent clean-clone verification results. Mark steps not verified where necessary.

### T-308 — Traceability and submission readiness
- Reconcile acceptance criteria AC-01 through AC-16 with requirements, user stories, tests, evidence, and known issues.
- Update only necessary documentation and sprint statuses supported by actual results.
- Inspect all diffs for formatting, accidental truncation, unrelated changes, and consistency.
- Verify required deliverables: frontend, backend, migrations/schema, README, SDLC documents, and evidence.
- Check Git status and submission readiness. Do not stage, commit, push, publish, or change repository visibility without my explicit approval.

## 4. Backend startup issue
The frontend at `http://localhost:5173` previously returned HTTP 200, but the backend health check at `http://localhost:8081/healthz` failed because nothing was listening on port 8081.

- Diagnose the startup issue by inspecting the project's existing startup instructions and entry point.
- If starting the server is safe and permitted, start it using the verified project command and check `/healthz`.
- Do not kill processes, change ports/configuration, rebuild executables, install packages, or run migrations without first explaining why and obtaining approval where needed.
- Do not alter or delete existing database records.

## 5. Operating rules
- English only. Work efficiently; reuse valid test results.
- Do not begin unrelated work or add features beyond the assignment.
- Do not modify `prompts/raw-prompt-log.md` or repair its ID collision as part of this finalization unless explicitly approved.
- No destructive database operations.
- Do not stage, commit, or push.
- If blocked by an approval requirement, stop that action, explain the blocker, and continue with independent safe tasks.
- Do not claim the project is fully complete or submission-ready until the required checks support that conclusion.

## 6. Final report
When finished, report:
1. Actual completion status for each T-301 through T-308.
2. Exact changed-file list and a concise diff summary.
3. Test/build commands executed and actual results, including failures.
4. AC-01 through AC-16 coverage gaps and unresolved known issues.
5. Whether the frontend and backend are running and the verified health-check result.
6. README and evidence readiness.
7. Git status, whether anything was staged/committed/pushed, and what explicit approval is still needed.
8. Recommended next action, prioritizing submission-critical work.

Start by inspecting the current state, then proceed with the safe, independent tasks. Do not wait for optional polish before completing README, evidence, and traceability.
```

---

## PROMPT-033 — T-207 read-only checkpoint audit

| Field | Value |
|---|---|
| Received | 2026-10-28 |
| Phase | Read-only checkpoint audit of the completed T-207 |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Read-only turn; logging deferred (no file edits were authorized). This audit inspected `backend/internal/api/auth.go` (verified unchanged from the T-106 implementation), confirmed both servers stopped, identified the missing T-207 run-prompt entry, and found the entry-ordering defect corrected in this repair. No secrets present. |

```text
Perform a READ-ONLY checkpoint audit of T-207 in E:\Inventory-Stock-Management.

Do not edit files, run tests, start or stop servers, or change database records.

1. Read CLAUDE.md and the T-207/T-208 sprint-plan sections.
2. Inspect git status and the current diffs without staging or committing anything.
3. Inspect Stock.jsx, Dashboard.jsx, their tests, route/navigation integration, and evidence/T-207_stock-dashboard-ui-tests.txt.
4. Verify that the evidence file exists and that its recorded results support the reported 31/31 frontend tests and build.
5. Identify which changes belong to T-205 and which belong to T-206 or earlier tasks. Do not attribute existing changes to the latest run without evidence.
6. Report whether the dashboard endpoint and its tests exist, which T-205 acceptance criteria have actually been verified, and which remain unverified.
7. Check whether prompts/raw-prompt-log.md contains the complete original T-205 prompt verbatim. Do not reconstruct missing text.
8. Do not run tests that modify data or make any changes during this audit.

Return findings and recommended next steps only. Do not begin T-208.
```

---

## PROMPT-034 — Prompt-log repair-source verification request

| Field | Value |
|---|---|
| Received | 2026-10-29 |
| Phase | Verification of the proposed raw prompt-log repair sources (planning only; no changes) |
| Tool / model | Claude Code CLI, model GLM (glm-5.3-flash) |
| Notes | Read-only verification turn; no changes made. This prompt defined the approved entry numbering (029–034) used in this repair. No secrets present. |

```text
STRICT LANGUAGE REQUIREMENT — ENGLISH ONLY

All responses, progress updates, audit reports, documentation, comments, and generated text must be written entirely in English. Do not use Chinese, Hindi, or any other language. If any tool output appears in another language, explain its relevant content in English. Before submitting your final report, verify that the entire report is in English. If you cannot comply, stop and report the limitation in English without making changes.

We are continuing the Inventory and Stock Management assignment in:

`E:\Inventory-Stock-Management`

CURRENT TASK: Verify the proposed raw prompt-log repair sources ONLY. Do not edit any files yet.

## Mandatory rules

1. Read `CLAUDE.md` and the relevant prompt-logging instructions before proceeding.
2. Inspect `prompts/raw-prompt-log.md` and the proposed source records using read-only operations.
3. Preserve every existing prompt's wording, formatting, code blocks, and metadata. Do not rewrite, summarize, reconstruct, or invent any missing prompt.
4. Do not edit, create, delete, rename, move, or format any files during this verification task.
5. Do not run tests or builds, start or stop servers, make API requests, modify the database, or change configuration.
6. Do not stage, commit, reset, or discard any Git changes. Preserve all existing work.
7. Do not read `.env` or expose passwords, JWTs, tokens, or other secrets.
8. Do not access unrelated assignments or their files.
9. If an original prompt cannot be located in a verifiable source, mark it unavailable. Do not reconstruct it from memory, summaries, headings, or previous assistant responses.
10. Distinguish verified facts from assumptions. Do not claim that any prompt is complete or verbatim unless you have checked its full original text.
11. Do not proceed to the actual repair. This task is verification and planning only.

## Verification tasks

### 1. Inspect the existing prompt log

- Report the exact current order of all entries from `PROMPT-001` through `PROMPT-028`.
- Identify the exact positions of `PROMPT-022` and `PROMPT-026`.
- Confirm whether any existing entries are duplicated, truncated, missing content, or otherwise damaged based on the read-only checks you can perform.
- Do not change existing IDs or prompt content.

### 2. Verify the six proposed missing prompts

The previous audit proposed these missing entries:

- `PROMPT-029`: T-203/T-204 read-only validation checkpoint.
- `PROMPT-030`: Manual verification pre-check before T-207.
- `PROMPT-031`: Approval to start the existing backend and frontend servers for manual verification.
- `PROMPT-032`: T-207 implementation instruction.
- `PROMPT-033`: T-207 read-only checkpoint audit.
- `PROMPT-034`: The prompt-log repair-source verification request from the previous turn.

For EACH entry:

1. Locate the exact original user prompt in a verifiable conversation record or project source.
2. Identify the exact source location.
3. Check the complete prompt, including all paragraphs, numbered lists, code blocks, paths, and explicit constraints.
4. Report one status:
   - **Available verbatim** — the entire original prompt is accessible and can be copied exactly.
   - **Partially available** — only some of the original text can be verified.
   - **Unavailable** — the original text cannot be located.
5. If the full text is unavailable, identify the missing material without attempting to recreate it.
6. Do not append or edit any prompt during this task.

Do not assume that a short description of a prompt proves that its complete original text is available.

### 3. Verify the proposed ordering repair

The previous audit proposed moving existing prompt blocks without changing their contents:

- Move `PROMPT-022` to immediately after `PROMPT-021`.
- Move `PROMPT-026` to immediately after `PROMPT-025`.

Inspect the current file and report the exact source and destination positions for both blocks. Confirm whether these moves alone would restore numeric order for the existing entries from `PROMPT-001` through `PROMPT-028`.

Do not move anything during this task.

### 4. Review the proposed disclosure line

The previous audit proposed adding this disclosure to the prompt-log recording-method section:

"Entries were reordered on 2026-10-29; all prompt text was preserved; entries 029–034 were appended from conversation records, with metadata reflecting their actual receipt dates."

Do not add this line. Instead:

- Locate the relevant recording-method section.
- Report the exact file and section where this line would be added.
- Determine whether the proposed wording accurately describes the planned repair.
- Identify any uncertainty, especially if a prompt's original text or receipt date cannot be independently verified.

### 5. Produce a minimal repair plan

Based only on verified information, propose the smallest safe repair that would:

- Preserve all existing prompt wording and IDs.
- Correct the existing ordering anomalies, if confirmed.
- Add only missing prompts whose complete original text is verifiably available.
- Avoid inventing prompt content, dates, metadata, or historical explanations.
- Document the repair accurately without unrelated code or configuration changes.

Clearly separate required changes from optional documentation changes.

### 6. Identify ambiguities requiring my decision

List any unresolved questions about:

- Missing prompt source text.
- Entry numbering and chronological order.
- Receipt dates or metadata.
- The proposed disclosure line.
- Any changes beyond the prompt log itself.

Do not make assumptions on my behalf.

## Required final report

Write the entire report in English only, using these sections:

1. Audit scope and read-only actions performed.
2. Current prompt IDs and exact order.
3. Ordering anomalies and proposed block positions.
4. Verification results for PROMPT-029 through PROMPT-034.
5. Disclosure-line review.
6. Minimal repair plan.
7. Ambiguities requiring my decision.
8. Recommended next step.

For every conclusion, state the supporting evidence or explain why it remains unverified.

Finish by confirming that no files were modified and that you are waiting for my approval before making any repair.

STOP after presenting the report. Do not perform the repair or begin another assignment task.
```

