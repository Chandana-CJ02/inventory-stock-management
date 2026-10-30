# Inventory and Stock Management System

A web-based Inventory and Stock Management System for a small business: a store owner manages product categories, products, stock movements (Stock In / Stock Out), and identifies low-stock products — with each user's inventory fully isolated.

Built with React 19 (functional components and hooks) on the frontend, Go (REST APIs) on the backend, and PostgreSQL for persistence. JWT-based authentication with bcrypt password hashing.

## Features

- Registration, login, and logout with JWT authentication and bcrypt password hashing
- Per-user independent inventory with strict access isolation (cross-user access → 404)
- Category management (create, rename, delete; per-user unique names; deletion guard when products exist)
- Product management (name, SKU, category, unit price, opening stock, minimum stock level)
- Product search by name or SKU; filter by category; current available stock display
- Stock In / Stock Out with automatic stock updates, atomic PostgreSQL transactions, and concurrency safety
- Immutable stock transaction history with product and type filters
- Dashboard: active product count, category count, low-stock list, five most recent transactions

## Tech stack

| Area | Technology |
|---|---|
| Frontend | React 19 (functional components, hooks), Vite 8 |
| Backend | Go 1.27 (REST APIs, `net/http`, `pgx/v5`) |
| Database | PostgreSQL 17 (relationships and constraints) |
| Authentication | JWT (HS256) + bcrypt password hashing |

## Project structure

```
├── backend/
│   ├── cmd/
│   │   ├── migrate/       # migration runner
│   │   └── server/        # API server entry point
│   ├── internal/
│   │   ├── api/           # HTTP handlers, middleware, integration tests
│   │   ├── auth/          # password hashing + JWT tokens
│   │   ├── config/        # environment config loader
│   │   └── database/      # PostgreSQL connection pool
│   └── migrations/        # embedded SQL migration files
├── frontend/
│   ├── src/
│   │   ├── pages/         # Login, Register, Dashboard, Products, Categories, Stock
│   │   ├── App.jsx        # authenticated shell + navigation
│   │   ├── api.js         # fetch wrapper
│   │   ├── auth.jsx       # session context
│   │   └── styles.css     # design tokens + shared styles
│   └── ...
├── docs/                  # SDLC documentation (requirements, architecture, test cases, etc.)
├── evidence/              # test execution evidence
└── prompts/               # raw prompt log
```

## Prerequisites

Verified with the following versions:

| Tool | Version |
|---|---|
| [Go](https://go.dev/) | 1.27.3 |
| [Node.js](https://nodejs.org/) | 23.11.0 |
| npm | 10.9.2 |
| [PostgreSQL](https://www.postgresql.org/) | 17.11 |
| [Git](https://git-scm.com/) | 2.43.0 |

## Environment configuration

1. Copy `.env.example` to `.env` in the project root (the file must sit in the **project root**, not inside `backend/` or `frontend/` — the server checks both the root and `backend/`).
2. Fill in the values:

| Variable | Example | Description |
|---|---|---|
| `PGHOST` | `localhost` | PostgreSQL host |
| `PGPORT` | `5432` | PostgreSQL port |
| `PGDATABASE` | `inventory_db` | Database name (created in step 3 below) |
| `PGUSER` | `inventory_app` | Application role (least privilege) |
| `PGPASSWORD` | *(your password)* | Application role password |
| `PGSSLMODE` | `disable` | SSL mode (local development) |
| `APP_PORT` | `8081` | Backend listen port |
| `JWT_SECRET` | *(long random string ≥ 32 chars)* | JWT signing secret |

> **Never commit `.env`** — it is already in `.gitignore`.

## Database setup

1. Ensure PostgreSQL 17 is running:
   ```powershell
   Get-Service postgresql-x64-17
   ```
2. Create the application role (replace `your_password` with a strong password):
   ```sql
   CREATE ROLE inventory_app LOGIN PASSWORD 'your_password';
   ```
3. Create the database owned by the application role:
   ```sql
   CREATE DATABASE inventory_db OWNER inventory_app;
   ```
4. Apply the embedded migrations:
   ```powershell
   cd backend
   go run ./cmd/migrate
   ```
   The migration runner is idempotent — it is safe to re-run.

## Installation and execution

### Backend

```powershell
cd backend
go build -o bin-server.exe ./cmd/server
.\bin-server.exe
```

The backend listens on the port set by `APP_PORT` (default 8081) and requires the `.env` file to be present in the project root or the `backend/` directory.

Health check:
```
GET http://localhost:8081/healthz
```

Expected response: `{"status":"ok","database":"up"}`

### Frontend

```powershell
cd frontend
npm install
npm run dev
```

The Vite dev server starts on `http://localhost:5173` and proxies `/api` requests to the backend at `http://localhost:8081`.

## Testing

### Frontend (vitest + jsdom)

```powershell
cd frontend
npm test
```

Expected: **33 tests passing across 5 files**.

### Backend (Go)

```powershell
cd backend
go test ./... -count=1
```

Expected: all packages `ok`.

### Production build

```powershell
cd frontend
npm run build
```

Output: `dist/` directory.

## API overview

All API endpoints are under `/api/v1/` and require a `Bearer` JWT token (except registration and login). The backend serves a health check at `/healthz`.

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/auth/register` | Create a new account |
| POST | `/api/v1/auth/login` | Log in and receive a JWT |
| POST | `/api/v1/auth/logout` | Log out (client-side token discard) |
| GET | `/api/v1/categories` | List categories |
| POST | `/api/v1/categories` | Create category |
| GET | `/api/v1/categories/:id` | Get category |
| PUT | `/api/v1/categories/:id` | Update category |
| DELETE | `/api/v1/categories/:id` | Delete category (blocked if products exist) |
| GET | `/api/v1/products` | List products (supports `?search=` and `?category_id=`) |
| POST | `/api/v1/products` | Create product |
| GET | `/api/v1/products/:id` | Get product |
| PUT | `/api/v1/products/:id` | Update product |
| POST | `/api/v1/products/:id/deactivate` | Deactivate product |
| POST | `/api/v1/transactions` | Record stock IN/OUT |
| GET | `/api/v1/transactions` | Transaction history (supports `?product_id=`, `?type=`, `?limit=`) |
| GET | `/api/v1/dashboard/summary` | Dashboard totals, low-stock list, recent transactions |

## Documentation

- [Requirements](docs/requirements.md) · [User stories](docs/user-stories.md) · [Sprint plan](docs/sprint-plan.md)
- [Architecture](docs/architecture.md) · [Test cases](docs/test-cases.md) · [API test plan](docs/api-test-plan.md)
- [Traceability matrix](docs/traceability-matrix.md) · [Evidence plan](docs/evidence-plan.md) · [Known issues](docs/known-issues.md)
- [Raw prompt log](prompts/raw-prompt-log.md)

## Known issues

See [docs/known-issues.md](docs/known-issues.md) for the current issue tracker.
