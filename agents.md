# Kvasir — Agent Guide

Guidance for AI coding assistants working on the Kvasir codebase. Read this before making changes.

## Project Identity

- **Name:** Kvasir — "The kernel of your mind"
- **Description:** Beautiful, techy, Nordic-inspired markdown knowledge base. Sync-first, API-first, plugin-extensible.
- **License:** AGPL-3.0
- **Repo:** git@github.com:mojoaar/kvasir.git

## Quick Commands

```sh
# Development (hot reload)
make dev            # Start both frontend (:3000) + backend (:8080) concurrently

# Build
make build          # Frontend build → embed into Go → single binary
make build-frontend # pnpm build (output into backend/internal/embed/dist/)
make build-backend  # go build with -ldflags="-s -w"

# Cross-compile
make build-all      # darwin (amd64/arm64) + windows + linux (amd64/arm64)

# Test & Lint
make lint           # pnpm lint + golangci-lint run
make test           # pnpm test + go test ./...
make check          # lint + test (full CI pass)

# Cleanup
make clean          # Remove dist/, binaries, node_modules/.cache
```

## Tech Stack

| Layer        | Choice                                        |
| ------------ | --------------------------------------------- |
| Frontend     | Next.js 16 App Router + shadcn/ui + Tailwind CSS v4 |
| Backend      | Go 1.26 + Gin                                      |
| Database     | SQLite (modernc.org/sqlite, pure Go, no CGO)       |
| Markdown     | TipTap (react) + Katex + Mermaid                   |
| State        | Zustand + TanStack React Query                     |
| Validation   | Zod                                                |
| CLI          | Cobra + Viper                                      |
| Auth         | JWT (golang-jwt/jwt/v5)                            |
| API Docs     | Swaggo (OpenAPI/Swagger)                           |
| Distribution | Single binary (Go embed frontend)                  |
| Package Mgr  | pnpm (frontend), Go modules (backend)              |

## Architecture

```
┌─────────────────────────────────────────────────┐
│  Desktop App   Web Browser   Android App         │
│             Next.js App Router + shadcn/ui        │
└──────────────────────┬──────────────────────────┘
                       │ REST API (:8080)
┌──────────────────────▼──────────────────────────┐
│                Go Backend (Gin)                  │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
│  │  Auth    │ │  Sync    │ │  Plugin Loader   │ │
│  └──────────┘ └──────────┘ └──────────────────┘ │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
│  │  Notes   │ │  Search  │ │  Themes/Tags     │ │
│  └──────────┘ └──────────┘ └──────────────────┘ │
│  ┌──────────────────────────────────────────────┐ │
│  │  SQLite (modernc.org) + Filesystem Storage   │ │
│  └──────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────┐ │
│  │  Embedded Frontend (//go:embed dist/)        │ │
│  └──────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────┘
```

## Project Structure

```
kvasir/
├── frontend/
│   ├── app/
│   │   ├── (auth)/                # Auth routes
│   │   ├── (dashboard)/           # Main app
│   │   │   ├── notes/
│   │   │   ├── search/
│   │   │   ├── settings/
│   │   │   │   ├── themes/
│   │   │   │   └── plugins/
│   │   │   └── page.tsx
│   │   └── layout.tsx
│   ├── components/
│   │   ├── ui/                    # shadcn/ui components
│   │   ├── editor/                # TipTap markdown editor
│   │   ├── sidebar/               # File tree
│   │   ├── themes/                # Theme provider + toggle
│   │   ├── plugins/               # Plugin loader UI
│   │   └── command-palette/       # Cmd+K
│   ├── lib/
│   │   ├── api/                   # API client
│   │   ├── store/                 # Zustand state
│   │   ├── plugins/               # Plugin system core
│   │   ├── themes/                # Theme definitions
│   │   │   ├── kvasir.ts
│   │   │   ├── dracula.ts
│   │   │   ├── nord.ts
│   │   │   ├── github.ts
│   │   │   └── cyberpunk.ts
│   │   └── utils/
│   ├── package.json
│   └── next.config.ts
│
├── backend/
│   ├── cmd/kvasir/main.go         # Entrypoint
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes.go
│   │   ├── storage/
│   │   │   ├── sqlite.go
│   │   │   ├── fts.go
│   │   │   └── models.go
│   │   ├── sync/
│   │   │   ├── engine.go
│   │   │   ├── conflict.go
│   │   │   └── encryption.go
│   │   ├── auth/
│   │   │   ├── jwt.go
│   │   │   ├── oauth.go
│   │   │   └── mfa.go
│   │   ├── embed/
│   │   │   └── dist/              # Built frontend (generated)
│   │   └── mcp/
│   ├── go.mod
│   └── go.sum
│
├── cli/
│   └── cmd/
│       ├── list/
│       ├── search/
│       ├── export/
│       └── sync/
│
├── themes/                        # Theme JSON definitions
│   ├── kvasir.json
│   ├── dracula.json
│   ├── nord.json
│   ├── github.json
│   └── cyberpunk.json
│
├── plugins/                       # Built-in plugins
│   ├── graph-view/
│   ├── kanban/
│   └── obsidian-importer/
│
├── kvasir                         # Final binary (generated)
├── Makefile
├── plan.md
├── agents.md
└── README.md
```

## Theme System

Themes are CSS custom properties mapped via Tailwind config. 5 built-in themes, each with dark + light mode variants:
**Kvasir** (Nordic, default), **Dracula**, **Nord**, **GitHub**, **Cyberpunk**.
Dark/light mode toggled independently per theme via `data-mode="dark|light"` on `<html>`.

Theme definitions live in:
- `frontend/lib/themes/*.ts` — TypeScript theme objects
- `themes/*.json` — JSON export format
- Tailwind config maps theme tokens to CSS variables

To add a theme:
1. Define its color palette as a TS object in `frontend/lib/themes/`
2. Add JSON export in `themes/`
3. Map CSS variables in Tailwind config
4. Add to the theme selector dropdown

Theme and mode stored in `localStorage` keys `kvasir-theme` and `kvasir-mode` (dark|light). Applied via `data-theme` and `data-mode` attributes on `<html>` before first paint (inline script) to prevent FOUC.

## Code Rules

### Always
- Run linter before marking task complete: `make lint`
- Run tests before marking task complete: `make test`
- Handle all Go error returns explicitly — **never use `_` to ignore errors**
- Validate all user inputs — never trust client data
- Log all errors with context (zerolog)
- Write tests for new features
- Maintain ≥ 80% line coverage for both frontend (vitest) and backend (go test -cover)
- Use `//go:embed` for embedding frontend — never copy files manually

### Never
- Commit secrets, credentials, or `.env` files
- Bypass error returns using blank identifiers (`_`) in Go
- Use deprecated, unmaintained, or archived packages
- Hardcode sensitive data — use environment variables
- Ignore linting warnings — fix them before committing
- Leave TODOs without tracking — create GitHub issues

### Banned Packages

| Package             | Reason                  | Use Instead         |
| ------------------- | ----------------------- | ------------------- |
| gorm                | Unnecessary abstraction | Raw SQL + sqlx      |
| go-sqlite3          | CGO dependency          | modernc.org/sqlite  |
| jwt-go              | Deprecated              | golang-jwt/jwt/v5   |
| gorilla/mux         | Unmaintained            | Gin                 |
| hashicorp/go-plugin | Overkill for MVP        | Custom sandboxed JS |

## Package Audit

Before adding any new dependency:
1. Check last commit date (< 6 months)
2. Check stars (> 1k) or downloads/week (> 10k)
3. Check for unresolved critical CVEs
4. Check issue count relative to stars (< 5%)
5. Verify test coverage exists
6. Run `pnpm audit` (frontend) or `go list -m -u all && go mod verify` (backend)

## Database

SQLite via `modernc.org/sqlite` (pure Go, no CGO). WAL mode enabled. FTS5 for full-text search.

All notes stored in `notes` table with FTS5 virtual table `notes_fts` for instant search.
Attachments stored on filesystem with metadata in `attachments` table.

Connection: Single connection pool. WAL mode handles reads/writes safely.

## Build Pipeline

```
pnpm build (frontend)
    │
    ▼
frontend/out/ ──copy──► backend/internal/embed/dist/
    │
    ▼
go build -ldflags="-s -w" -o kvasir cmd/kvasir/main.go
    │
    ▼
./kvasir                    # Single binary, serves on :8080
```

The Go binary serves the embedded frontend via `//go:embed` and the API on the same port.
In development, the frontend runs separately on `:3000` with hot reload, proxying API calls to `:8080`.

## Testing & Linting

```sh
# Frontend
cd frontend
pnpm lint                  # ESLint + Prettier
pnpm test -- --coverage    # Vitest (with coverage)

# Backend
cd backend
golangci-lint run                       # Go linter
go test -coverprofile=coverage.out ./... # Go tests (with coverage)

# Full CI pass
make check             # Runs lint + test + coverage checks
```

## Release Workflow

- **Normal push** — standard `git add`, `git commit`, `git push origin main`
- **Tagged release** — version bumps and git tags only when explicitly requested
- Never create tags or GitHub releases autonomously

## Gotchas

_To be populated as we discover project-specific issues._

## Environment Variables

```
KVASIR_DB_PATH=~/.kvasir/kvasir.db
KVASIR_VAULT_PATH=~/.kvasir/vaults/
KVASIR_PORT=8080
KVASIR_JWT_SECRET=<32+ char secret>
KVASIR_LOG_LEVEL=debug|info|warn|error
```
