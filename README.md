# Go Clean Template

A deliberately small, production-minded Go REST API starter using feature-first
clean architecture. It is the server-side counterpart of the Flutter BLoC
template: one composition root, typed errors, an interactive generator, and
agent instructions that keep generated projects consistent.

## Included

- Feature-first clean architecture (`domain`, `data`, `http`)
- Gin router with request ID, recovery, access logs, CORS, and an auth middleware hook
- Environment-first configuration via `.env` and process env vars
- Typed `apperr` failures (validation with field errors, unauthorized, forbidden, not found, conflict, data, internal)
- `slog` logging with redaction for tokens, passwords, secrets, and cookies
- Single composition root in `internal/app`
- Sample `welcome` feature proving the full vertical slice
- Optional Postgres wiring (`pgx`) without requiring a database to run
- Makefile, multi-stage Dockerfile, and docker compose
- `AGENTS.md` for AI coding agents
- Interactive `scripts/go-clean-init.sh` project generator

The existing thinner Gin skeleton in `go-rest-template` remains available. This
template intentionally adds architecture, tests, and agent rules rather than
forcing unused SDKs into every service.

## Generate a project

From this repository, run:

```sh
./scripts/go-clean-init.sh
```

The initializer prompts for:

1. Project folder name
2. Go module path
3. App display name
4. HTTP port
5. CORS allowed origins
6. Database choice: `none` (local welcome repository) or `postgres`

It copies the template, rewrites the module path, stamps `AGENTS.md` and
`.env`, aligns the `go` directive with your installed toolchain, then runs
`go mod tidy`, `gofmt`, and `go test`.

Interactive terminals get a color UI, review screen, and Esc/Ctrl+C exit.
Non-interactive or `NO_COLOR` environments use plain prompts. You can also pass:

```sh
./scripts/go-clean-init.sh github.com/you/my-service
# or
NEW_MODULE=github.com/you/my-service ./scripts/go-clean-init.sh
```

## Run this template directly

```sh
cp .env.example .env
make tidy
make run
# http://localhost:8080/health
# http://localhost:8080/api/v1/welcome
```

Configuration lives in `.env`:

```dotenv
APP_NAME=Go Clean Template
APP_ENV=local
HTTP_PORT=8080
HTTP_ALLOWED_ORIGINS=*
```

- `APP_ENV` accepts `local`, `dev`/`development`, or `prod`/`production`.
- Leave `DATABASE_URL` empty to use the in-process welcome repository.
- Set `DATABASE_URL` to a Postgres DSN to switch the welcome repository to the
  Postgres implementation (expects an optional `welcome_messages` table).

## Structure

```text
cmd/server/                 # Process entry point
internal/
├── app/                    # Composition root
├── config/                 # Environment-backed configuration
├── apperr/                 # Typed failures + HTTP status mapping
├── logging/                # slog setup and redaction
├── httpserver/             # Server, middleware, JSON helpers
└── features/
    ├── health/http/        # Liveness endpoint
    └── welcome/
        ├── domain/         # Entity, repository contract, use case
        ├── data/           # Local and optional Postgres repositories
        └── http/           # Handlers and route registration
scripts/go-clean-init.sh
api/openapi.yaml
AGENTS.md
```

Dependencies point inward: `http` uses `domain`; `data` implements `domain`
contracts; `domain` does not import Gin, SQL drivers, or env loaders.
`internal/app` is the only place that assembles implementations.

## Add a feature

1. Create `internal/features/<name>/domain` first: entities, repository
   interfaces, and focused use cases.
2. Add implementations under `data`, keeping JSON/SQL details out of domain.
3. Add handlers under `http` and register routes from `internal/app`.
4. Map failures with `apperr` and `httpserver.WriteError` — do not invent status
   codes in handlers.
5. Add tests for use cases, repositories, and handlers (`httptest`).

## Verify

```sh
make fmt
make test
make lint   # requires golangci-lint
```

## AI agent instructions

[`AGENTS.md`](AGENTS.md) holds the architecture, API, and testing rules for AI
coding agents. Include real endpoint documentation and acceptance criteria with
each feature request so the agent follows the API contract instead of inventing
request or response fields.

The initializer stamps the generated `AGENTS.md` with the reviewed project name,
module path, app name, HTTP port, database choice, and Go version.
