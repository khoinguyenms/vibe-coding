# vide-be

Go backend service using PostgreSQL and sqlc.

## Tech Stack

- **Go 1.23**
- **uber-fx** — dependency injection + lifecycle
- **zap** — structured logging
- **Gin** — HTTP framework
- **pgx/v5** — PostgreSQL driver + connection pool
- **sqlc** — type-safe SQL code generation
- **golang-migrate** — database migrations

## Project Layout

```
vide-be/
├── cmd/api/              # application entrypoint
├── internal/
│   ├── config/           # env config loader
│   ├── database/         # pgxpool wiring (with fx lifecycle)
│   ├── handler/          # HTTP handlers
│   ├── service/          # business logic
│   ├── repository/       # data-access layer (wraps sqlc)
│   ├── middleware/       # custom HTTP middleware
│   ├── router/           # route registration
│   └── module/           # fx modules (composition root)
├── db/
│   ├── migrations/       # *.up.sql / *.down.sql
│   ├── queries/          # sqlc query definitions
│   └── sqlc/             # generated code (do not edit)
├── pkg/
│   ├── logger/
│   └── response/
├── sqlc.yaml
├── Makefile
└── go.mod
```

## Getting Started

1. Copy env file:
   ```bash
   cp .env.example .env
   ```

2. Install deps:
   ```bash
   go mod tidy
   ```

3. Install tools:
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

4. Run migrations:
   ```bash
   make migrate-up
   ```

5. Regenerate sqlc code (after editing `db/queries/*.sql`):
   ```bash
   make sqlc
   ```

6. Run server:
   ```bash
   make run
   ```

## API Endpoints

- `GET  /health`
- `GET  /api/v1/users`
- `POST /api/v1/users`
- `GET  /api/v1/users/{id}`
