# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

Monorepo with two independent apps:

- `vibe-be/` — Go 1.23 backend (Gin + uber-fx + pgx + sqlc + golang-migrate)
- `vibe-fe/` — Next.js 16 frontend (React 19, Tailwind v4, shadcn/ui, pnpm)

There is no root-level build; commands run inside each app's directory.

## Backend (`vibe-be/`)

### Common commands

```bash
make run               # go run ./cmd/api
make build             # produces bin/$(APP_NAME)
make tidy              # go mod tidy
make sqlc              # regenerate db/sqlc/* after editing db/queries/*.sql
make migrate-up        # apply all up migrations
make migrate-down      # roll back one migration
make migrate-create name=add_foo   # scaffold seq-numbered up/down sql
make test              # go test ./... -v -cover

# Run a single test
go test ./internal/service -run TestUserService_Create -v
```

`migrate-*` targets build the DSN from `DB_USER/DB_PASSWORD/DB_HOST/DB_PORT/DB_NAME/DB_SSLMODE` env vars — export them (or `set -a; source .env; set +a`) before running migrations. The app itself loads `.env` via godotenv at startup.

### Architecture

Composition root is `cmd/api/main.go`, which wires uber-fx modules from `internal/module/`:

```
ConfigModule → LoggerModule → DatabaseModule → RepositoryModule → ServiceModule → HTTPModule
```

Each layer is its own package and is injected via fx — never construct dependencies manually. To add a new feature, follow the existing user pipeline:

1. Add SQL to `db/queries/<entity>.sql` and a migration in `db/migrations/`.
2. Run `make sqlc` to regenerate `db/sqlc/` (treat that directory as build output — do not hand-edit).
3. Add `internal/repository/<entity>_repository.go` wrapping `sqlc.Querier`, plus an interface in `repository/interface.go`. Register in `module/repository.go`.
4. Add `internal/service/<entity>_service.go` and `service/interface.go`. Register in `module/service.go`.
5. Add `internal/handler/<entity>_handler.go` and `handler/interface.go`. Register provider in `module/http.go`.
6. Wire routes in `internal/router/router.go` under the `/api/v1` group.

Database wiring is unusual: `DatabaseModule` provides both `*pgxpool.Pool` and a `sqlc.Querier` (via `sqlc.New` annotated as the interface) so repositories depend on the interface, not the pool. The lifecycle `OnStart`/`OnStop` for the pool lives in `internal/database/`.

HTTP server uses `fx.Hook` in `module/http.go` for graceful start/stop with a separate `net.Listen` (so the server is fully bound before `OnStart` returns).

### Cross-cutting

- `pkg/response` — uniform success/error JSON envelope and sentinel errors (`ErrNotFound`, `ErrInvalidInput`, `ErrInternal`). Handlers must use `response.Success`/`response.Error` rather than `c.JSON` directly.
- `pkg/validation` — wraps gin binding errors; pair with `response.ValidationError`.
- `pkg/logger` — zap wrapper. Use `log.Ctx(ctx).Error(...)` inside handlers/services to propagate request-scoped fields injected by `middleware.Tracing` + `middleware.RequestLogger`.
- Handlers strip `Password` from `sqlc.User` before responding — keep that pattern when adding fields the DB stores but the API should not return.

## Frontend (`vibe-fe/`)

### Common commands

```bash
pnpm dev               # Next dev server on :3000
pnpm build             # next build (uses babel-plugin-react-compiler)
pnpm start             # serve a production build
pnpm lint              # eslint (next core-web-vitals)
```

Use **pnpm** — the lockfile is `pnpm-lock.yaml` and there is a `pnpm-workspace.yaml`.

### Structure

- App Router under `src/app/` (RSC enabled).
- `src/components/ui/` is shadcn/ui output; `components.json` is configured for the `radix-nova` style on the `neutral` base color with `lucide` icons. Add components via `pnpm dlx shadcn@latest add <name>` rather than hand-rolling.
- Path aliases: `@/components`, `@/lib`, `@/components/ui`, `@/hooks` (see `components.json`).
- Tailwind v4 via `@tailwindcss/postcss`; design tokens live in `src/app/globals.css`.
