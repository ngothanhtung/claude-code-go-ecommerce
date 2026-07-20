# go-tutorials

Go RESTful service (Gin + GORM/Postgres), feature-based modular layout.

## Features

- Auth: register / login / refresh / logout / forgot / reset (JWT access + refresh)
- Users: CRUD + `/me` + role guard
- RBAC: roles/permissions models + admin guard
- Uploads: local file storage (images/pdf)
- Health: `/healthz`, `/readyz`, `/swagger/*any`

## Setup

1. Copy `configs/.env.example` to `configs/.env` and adjust.
2. Start Postgres + Redis via Docker Compose:

   ```bash
   cd docker/postgresql && docker compose up -d   # Postgres 17 (port 5432)
   cd docker/redis     && docker compose up -d   # Redis 7   (port 6379)
   ```

3. `make migrate-up`
4. `make run`

## Endpoints (prefix /api/v1)

- POST /auth/register
- POST /auth/login
- POST /auth/refresh
- POST /auth/logout
- POST /auth/forgot-password
- POST /auth/reset-password
- GET /users/me
- GET /users
- GET /users/:id
- PUT /users/:id
- DELETE /users/:id
- POST /uploads (multipart, field "file")
- GET /healthz
- GET /readyz
- GET /swagger/\*any

## Project layout

cmd/ entrypoint · internal/ business (features, common, db, config, app) · pkg/ reusable libs · migrations/ SQL.

## Notes

- Password reset link is logged to console in dev (no email sender yet).
- Storage is local-only; `Storage` interface allows adding S3 later.
- Redis is used for rate-limiting and as a future refresh-token store.
- A dev admin is auto-seeded unless APP_ENV=production (email admin@go-tutorials.local / password Admin@123456).

## How to run

### 1. Start Postgres + Redis (once)

```bash
cd docker/postgresql && docker compose up -d    # Postgres 17, port 5432
cd ../redis          && docker compose up -d    # Redis 7,    port 6379
```

Verify both are up:

```bash
docker ps --filter "name=claude_code_"
# claude_code_postgres  Up  ...  0.0.0.0:5432->5432/tcp
# claude_code_redis     Up  ...  0.0.0.0:6379->6379/tcp
```

### 2. Create env file (once)

```bash
cp configs/.env.example configs/.env
```

`configs/.env` is git-ignored — edit it freely. At minimum, change
`JWT_SECRET` before any production deployment.

### 3. Run migrations (apply schema)

```bash
make migrate-up
```

Expected output (first run):

```text
1/u init (13.264333ms)
```

Subsequent runs print `no change` (the command is idempotent).

> `make migrate-up` prefers the `migrate` CLI binary (install via
> `brew install golang-migrate`) and falls back to
> `go run -tags tools ...` if the binary is not installed.

### 4. Start the server

```bash
make run
```

Expected log:

```text
seed: admin user created: admin@go-tutorials.local
INFO  app/app.go:45  server starting  {"addr": ":8080"}
Listening and serving HTTP on :8080
```

A dev admin is auto-seeded (email `admin@go-tutorials.local`,
password `Admin@123456`) — only when `APP_ENV != production`.

### 5. Quick smoke test

```bash
# health
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz

# register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"me@example.com","name":"Me","password":"Password123"}'

# login -> grab access token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"me@example.com","password":"Password123"}' \
  | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["access_token"])')

# call an authenticated endpoint
curl http://localhost:8080/api/v1/users/me -H "Authorization: Bearer $TOKEN"

# upload a file (multipart, field "file")
curl -X POST http://localhost:8080/api/v1/uploads \
  -H "Authorization: Bearer $TOKEN" \
  -F 'file=@./some-image.jpg;type=image/jpeg'

# Swagger UI
open http://localhost:8080/swagger/index.html
```

## Common commands

| Command             | Purpose                                     |
| ------------------- | ------------------------------------------- |
| `make run`          | Run the server (foreground)                 |
| `make migrate-up`   | Apply schema                                |
| `make migrate-down` | Roll back one migration                     |
| `make tidy`         | `go mod tidy`                               |
| `make swagger`      | Regenerate `api/openapi.yaml` from comments |
| `go test ./...`     | Run unit tests (currently under `pkg/jwt`)  |

## Manual API tests

Open the numbered request files in `http/` with the VS Code REST Client or an
IntelliJ-based HTTP Client. Run the login request at the top of each protected
domain file first; response variables automatically feed tokens and created IDs
into the following requests.

- `00-health-auth.http` — health, register, JWT refresh, reset, logout
- `01-users.http` — profile and user CRUD
- `02-catalog.http` — categories, products, search, related items, promos
- `03-cart-wishlist.http` — cart and wishlist lifecycle
- `04-orders-reviews.http` — checkout, status updates, reviews
- `05-notifications-uploads.http` — notification state and multipart upload

## Stopping & cleanup

```bash
# stop the server
# Ctrl+C if running in foreground, otherwise:
pkill -f 'cmd/server'

# stop containers
cd docker/postgresql && docker compose down
cd ../redis          && docker compose down

# also remove data volumes (WARNING: wipes the database)
cd docker/postgresql && docker compose down -v
```
