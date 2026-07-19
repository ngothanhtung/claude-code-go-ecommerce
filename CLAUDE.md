# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go RESTful API (Gin + GORM/Postgres) with feature-based modular layout. Supports JWT access + refresh auth, RBAC, local file uploads, and Swagger docs. Designed for easy extension — each feature lives in its own package under `internal/features/`.

## Commands

Run from this directory (`go-tutorials/`):

```bash
# Start Postgres + Redis (from repo root)
cd docker/postgresql && docker compose up -d
cd docker/redis     && docker compose up -d

# Create env file (once)
cp configs/.env.example configs/.env   # edit as needed

# Apply schema
make migrate-up
make migrate-down                       # roll back one migration

# Run the server (foreground, http://localhost:8080)
make run

# Build only (no run)
go build -o /tmp/server ./cmd/server

# Lint / static analysis
go vet ./...

# Tests
go test ./...                           # all tests (pkg/jwt)
go test ./pkg/jwt/ -v                   # one package, verbose
go test ./pkg/jwt/ -run TestAccessToken # single test by name

# Swagger (regenerate from code comments)
make swagger

# Go module cleanup
make tidy
```

## Architecture

One-directional dependency flow: **handler → service → repository → db**. Handlers never call GORM directly.

### Directory layout

```text
cmd/server/main.go              # Cobra entrypoint (serve, migrate)
internal/
├── app/app.go                  # App struct: db, redis, logger, jwt, gin engine
├── app/router.go               # Middleware chain + feature route registration
├── config/config.go            # Viper + godotenv loader
├── db/postgres.go              # GORM + pgx connection pool + Ping helper
├── db/redis.go                 # go-redis client + PingRedis helper
├── common/
│   ├── errors/                 # AppError (code + HTTP status + message)
│   ├── response/               # Envelope {code, message, data}
│   ├── paging/                 # ?page=&page_size=&sort= parser
│   ├── middleware/              # requestid, logger, recovery, cors, ratelimit, auth, audit
│   └── storage/                # Storage interface + local filesystem impl
├── features/
│   ├── auth/                   # register/login/refresh/logout/forgot/reset
│   ├── user/                   # CRUD + /me + role management
│   ├── rbac/                   # Role/Permission models + AdminGuard
│   ├── health/                 # /healthz, /readyz (at root, no /api/v1 prefix)
│   ├── uploads/                # multipart file upload to local storage
│   └── seed/                   # dev admin + roles seeder (runs when APP_ENV != production)
pkg/
├── jwt/                        # Token manager (sign, parse, refresh rotation)
├── logger/                     # Zap wrapper (dev=human, prod=JSON)
└── validator/                  # Struct + WrapBind validation helpers
migrations/                     # golang-migrate SQL files
deployments/Dockerfile          # Multi-stage build (golang:1.23 → distroless)
http/api.http                   # VS Code REST Client collection
```

### Feature package convention

Each feature under `internal/features/<name>/` owns:

| File | Role |
| --- | --- |
| `model.go` | GORM entity (maps to DB table) |
| `dto.go` | Request/response structs with `validate` tags |
| `repository.go` | `Repository` interface + GORM implementation |
| `service.go` | Business logic (calls repository, never GORM) |
| `handler.go` | Gin handler (binds DTO, calls service, returns envelope) |
| `routes.go` | `Register()` function mounting routes on `*gin.RouterGroup` |

New features follow this convention. Import cycle: do NOT import one feature from another — shared logic goes in `common/` or `pkg/`.

### Key patterns

- **Response envelope**: Every endpoint returns `{"code": 0, "message": "ok", "data": ...}` or `{"code": 4000, "message": "...", "data": null}`. Use `response.OK()` / `response.Error()`.
- **Errors**: All errors flow through `common/errors.AppError`. Services return `AppError` constructors (`NewNotFound`, `NewConflict`, `NewInvalidCredentials`, etc.); the handler calls `response.Error(c, err)`.
- **Auth middleware**: Validates `Authorization: Bearer <token>`, stores `*jwt.Claims` in context via `c.Set("claims", claims)`. Read it with `middleware.Claims(c)`.
- **RBAC guard**: `middleware.RequireRole("admin")` or `rbac.AdminGuard()` — applied per-route group.
- **Paging**: `?page=1&page_size=20&sort=-created_at&q=keyword` parsed by `paging.Parse(c)`; results wrapped in `paging.Result`.
- **Health endpoints**: Mounted at root (`/healthz`, `/readyz`) — NOT under `/api/v1`. Ping closures must create fresh context per call (the old shared-context pattern causes "context canceled").
- **File uploads**: `POST /api/v1/uploads` multipart, validated by MIME (jpeg/png/webp/pdf) and size limit (`MAX_UPLOAD_MB`), saved as `YYYY/MM/<uuid>.<ext>` under `uploads/`.
- **JWT flow**: Register/Login → returns `access_token` (15 min) + `refresh_token` (7 days). Refresh rotates both, revokes old refresh token (stored in DB, not Redis). Logout revokes the refresh token.

## Infrastructure

- **Postgres 17**: `docker/postgresql/docker-compose.yml` — user `postgres`, password `postgres_password`, database `claude_code_flutter`, port 5432.
- **Redis 7**: `docker/redis/docker-compose.yml` — no auth (dev), port 6379. Used for rate-limiting.
- **Migrate CLI**: Install via `brew install golang-migrate`. Script at `scripts/migrate.sh` loads `.env` and runs `migrate up`.
- **Tools registration**: `tools/tools.go` (build tag `tools`) registers postgres + file drivers so `go run` of the migrate CLI works without a pre-built binary.

## Config

All config lives in `configs/.env` (git-ignored). Copy from `.env.example`:

| Variable | Default | Notes |
| --- | --- | --- |
| `APP_ENV` | `development` | `production` disables seed + debug logging |
| `APP_PORT` | `8080` | |
| `DB_HOST/PORT/USER/PASS/NAME` | `localhost:5432/postgres/postgres_password/claude_code_flutter` | Must match docker compose |
| `REDIS_ADDR` | `localhost:6379` | |
| `REDIS_PASSWORD` | *(empty)* | Must match redis compose (no auth) |
| `JWT_SECRET` | `change-me-...` | Change before production |
| `JWT_ACCESS_TTL_MIN` | `15` | Access token lifetime |
| `JWT_REFRESH_TTL_HOUR` | `168` | Refresh token lifetime (7 days) |
| `UPLOAD_DIR` | `./uploads` | Local storage root |
| `MAX_UPLOAD_MB` | `10` | |
| `RATE_LIMIT_PER_MIN` | `60` | Per-IP, Redis-backed |
| `CORS_ORIGINS` | `http://localhost:3000,http://localhost:5173` | Comma-separated |
| `AUDIT_ENABLED` | `true` | Writes to `audit_logs` table; skips `/healthz`, `/metrics` |

## REST Client

`http/api.http` is a VS Code REST Client collection covering all endpoints. Install the [REST Client extension](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) to use it. Requests chain automatically via variables (`{{login.response.body.data.access_token}}` etc.).

## Testing conventions

- Unit tests live next to the package they test (e.g. `pkg/jwt/jwt_test.go`).
- Repository/service mocks use `go.uber.org/mock` (gomock) — generate with `mockgen`.
- DB-dependent tests use a separate test database; set `DB_NAME` env var accordingly.
- Integration tests (API end-to-end) are deferred — currently only unit tests under `pkg/jwt` exist.

## Demo credentials

The dev seeder creates an admin user on first run (when `APP_ENV != production`):

- Email: `admin@go-tutorials.local`
- Password: `Admin@123456`

These are for local development only — never commit production credentials.
