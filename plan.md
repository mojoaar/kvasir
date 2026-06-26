# Kvasir — Implementation Plan

**Tagline:** "The kernel of your mind"
Beautiful, techy, Nordic-inspired markdown knowledge base. Sync-first, API-first, plugin-extensible.

**License:** AGPL-3.0 | **Repo:** git@github.com:mojoaar/kvasir.git

---

## Tech Stack

| Layer        | Choice                                              |
| ------------ | --------------------------------------------------- |
| Frontend     | Next.js 16 App Router + shadcn/ui + Tailwind CSS v4 |
| Backend      | Go 1.26 + Gin                                       |
| Database     | SQLite (modernc.org/sqlite, pure Go)                |
| State        | Zustand + TanStack React Query                      |
| Validation   | Zod                                                 |
| Markdown     | TipTap (react) + Katex + Mermaid                    |
| CLI          | Cobra + Viper                                       |
| Auth         | JWT (golang-jwt/jwt/v5)                             |
| API Docs     | Swaggo (OpenAPI/Swagger)                            |
| Testing      | Vitest (frontend), Go stdlib + testify (backend)    |
| Distribution | Single binary (Go embed)                            |

---

## Phase 0: MVP Core (~4-6 weeks)

Objective: Offline-first markdown knowledge base with 5 themes, full-text search,
REST API, CLI, and single-binary distribution.

### 0.1 — Project Scaffold [1-2 days]
- [ ] Root `Makefile` with targets: dev, build, build-frontend, build-backend, build-all, lint, test, check, clean
- [ ] `frontend/` scaffold: `pnpm create next-app` with App Router, TypeScript, Tailwind v4
- [ ] `backend/` scaffold: `go mod init`, Gin skeleton, `cmd/kvasir/main.go`
- [ ] `cli/` scaffold: Cobra root command
- [ ] `.gitignore`, `.env.example`, `.golangci.yml`, `eslint.config.mjs`
- [ ] Verify: `make dev` starts both frontend (:3000) and backend (:8080)

### 0.2 — Backend Core [2-3 days]
- [ ] SQLite schema: `notes`, `notes_fts` (FTS5), `tags`, `note_tags`, `vaults`, `attachments`, `versions`, `themes`, `plugins`, `plugin_permissions`
- [ ] Storage layer: `sqlite.go` (open, migrate, WAL mode), `fts.go` (search), `models.go`
- [ ] Health endpoint: `GET /api/v1/health`
- [ ] Gin router with middleware chain (logging, CORS, recovery)
- [ ] Verify: `go test ./...` passes

### 0.3 — Frontend Scaffold [2-3 days]
- [ ] shadcn/ui init + component library setup
- [ ] Root layout: `app/layout.tsx` with font loading, theme script, providers
- [ ] Dashboard layout: sidebar shell, main content area
- [ ] API client: `lib/api/client.ts` with TanStack React Query
- [ ] Zustand stores: notes, theme, UI state
- [ ] Verify: frontend renders, proxies API calls to backend

### 0.4 — Theme System [2-3 days]
- [ ] CSS custom properties for all 5 themes (dark + light variants each) in Tailwind config
- [ ] Theme definitions: `frontend/lib/themes/kvasir.ts`, `dracula.ts`, `nord.ts`, `github.ts`, `cyberpunk.ts` — each exports `dark` and `light` palette objects
- [ ] Theme provider component: reads/writes `localStorage` keys `kvasir-theme` and `kvasir-mode`
- [ ] Inline script to apply theme + mode before first paint (prevents FOUC)
- [ ] Dark/Light mode toggle (per‑theme, persists via `kvasir-mode` localStorage key)
- [ ] Theme selector in Settings
- [ ] JSON theme exports in `themes/*.json` (each with dark + light variants)
- [ ] Verify: all 5 themes render correctly, toggle works, mode persists on reload

### 0.5 — Markdown Editor (TipTap) [3-4 days]
- [ ] TipTap editor setup with extensions: starter-kit, table, footnote, math (Katex), code-block, mermaid
- [ ] Editor component with toolbar (bold, italic, headings, lists, tables, code, math, mermaid)
- [ ] Split view: editor + live preview
- [ ] Auto-save (debounced, 2s) + manual save (Cmd+S)
- [ ] Keyboard shortcuts within editor
- [ ] Verify: type markdown, tables render, math renders, mermaid diagrams render

### 0.6 — Sidebar with File Tree [2-3 days]
- [ ] Note list: create, rename, delete notes from sidebar
- [ ] Folder tree: nested folder structure
- [ ] Drag-and-drop to reorganize
- [ ] Note count per folder
- [ ] Collapsible sections
- [ ] Verify: CRUD notes from sidebar, drag to reorganize

### 0.7 — Notes API [2-3 days]
- [ ] `GET /api/v1/notes` — list notes (pagination, vault filter)
- [ ] `POST /api/v1/notes` — create note
- [ ] `GET /api/v1/notes/:id` — get single note
- [ ] `PUT /api/v1/notes/:id` — update note
- [ ] `DELETE /api/v1/notes/:id` — soft delete note
- [ ] Zod validation schemas for all request/response bodies
- [ ] Verify: CRUD works via curl, OpenAPI docs reflect endpoints

### 0.8 — Full-Text Search [1-2 days]
- [ ] `GET /api/v1/search?q=query` — FTS5 search endpoint
- [ ] Search input in sidebar header
- [ ] Results panel with highlighted matches
- [ ] Search by tag: `GET /api/v1/search/tags?q=`
- [ ] Verify: search returns results < 100ms with 10k notes

### 0.9 — Tags API + UI [1-2 days]
- [ ] `GET/POST /api/v1/tags` — list/create tags
- [ ] `GET/PUT/DELETE /api/v1/tags/:id` — single tag operations
- [ ] `POST/DELETE /api/v1/notes/:id/tags` — tag a note
- [ ] Tag chips in sidebar + note editor
- [ ] Tag color palette (cyclic 8-color rotation)
- [ ] Verify: add/remove tags from notes, filter by tag

### 0.10 — Command Palette [2-3 days]
- [ ] Cmd+K / Ctrl+K global shortcut
- [ ] Command palette overlay (cmdk)
- [ ] Commands: search notes, create note, switch theme, settings, keyboard shortcuts ref
- [ ] Keyboard-first navigation
- [ ] Verify: Cmd+K opens palette, type to search, Enter to execute

### 0.11 — REST API Documentation [1-2 days]
- [ ] Swaggo annotations on all Go handlers
- [ ] `GET /api/v1/docs` — Swagger UI page
- [ ] OpenAPI JSON at `/api/v1/swagger.json`
- [ ] Verify: all endpoints documented, Swagger UI functional

### 0.12 — CLI [2-3 days]
- [ ] `kvasir list` — list notes
- [ ] `kvasir search <query>` — full-text search
- [ ] `kvasir export <id>` — export note as markdown
- [ ] `kvasir create <title>` — create note from CLI
- [ ] Cobra subcommands with flags
- [ ] Verify: CLI works against running backend

### 0.13 — Single Binary Build [1-2 days]
- [ ] `make build` pipeline: frontend build → copy to `backend/internal/embed/dist/` → `go build`
- [ ] Go `//go:embed` for serving frontend static files
- [ ] `make build-all` cross-compile for 5 targets (darwin/amd64, darwin/arm64, windows/amd64, linux/amd64, linux/arm64)
- [ ] Binary size check (< 100MB target)
- [ ] Verify: single binary serves frontend + API on `:8080`

### 0.14 — Welcome Note + Tutorial [1 day]
- [ ] First-run: create welcome note with tutorial content
- [ ] Welcome note covers: editor basics, themes, search, keyboard shortcuts
- [ ] Verify: fresh install shows welcome note

### 0.15 — Testing + Coverage + Polish [3-4 days]
- [ ] Frontend: vitest tests for theme system, API client, stores, command palette
- [ ] Backend: Go table-driven tests for all handlers, storage layer, search, auth
- [ ] Coverage collection wired into `make test`:
  - `go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out`
  - `pnpm test -- --coverage`
- [ ] Coverage thresholds enforced in `make check`:
  - Backend: ≥ 80% line coverage (fail if below)
  - Frontend: ≥ 80% line coverage (fail if below)
- [ ] `make check` passes (lint + test + coverage)
- [ ] Performance check: startup < 2s, search < 100ms, memory < 500MB
- [ ] Cross-platform smoke test (macOS binary)

### MVP Success Checklist
- [ ] Works offline with local SQLite storage
- [ ] Create, edit, delete notes with full markdown
- [ ] 5 built-in themes toggle correctly
- [ ] Full-text search works instantly
- [ ] REST API fully documented with OpenAPI
- [ ] Single binary runs on macOS (ARM + Intel)
- [ ] CLI can list/search/export notes
- [ ] All linting and tests pass (`make check` green)
- [ ] Command palette functional

---

## Phase 1: Sync Engine (~4-6 weeks)

_Depends on: Phase 0 complete_

- [ ] **1.1** Version vector sync engine: push/pull, conflict detection
- [ ] **1.2** `POST /api/v1/sync/push` + `GET /api/v1/sync/pull` + `GET /api/v1/sync/status`
- [ ] **1.3** Self-hostable sync server — single source of truth for all clients (desktop, Android, web)
- [ ] **1.4** pCloud + Google Drive OAuth integration
- [ ] **1.5** Offline-first with background sync queue
- [ ] **1.6** End-to-end encryption for sync payloads
- [ ] **1.7** Multi-vault support + vault-level passwords
- [ ] **1.8** `POST /api/v1/vaults/:id/lock` + `/unlock`
- [ ] **1.9** Sync conflict resolution UI

---

## Phase 2: Teams & Auth (~4-6 weeks)

_Depends on: Phase 1 complete_

- [ ] **2.1** JWT auth: `POST /api/v1/auth/login`, `/register`, `/refresh`, `/logout`
- [ ] **2.2** OAuth2 flows: Google + GitHub (`/api/v1/auth/oauth/:provider`)
- [ ] **2.3** API key generation and validation
- [ ] **2.4** MFA setup + verify (`/api/v1/auth/mfa/setup`, `/verify`)
- [ ] **2.5** User roles: Admin, Editor, Viewer
- [ ] **2.6** User management: `GET/PUT /api/v1/users/me`, `GET /api/v1/users/:id`
- [ ] **2.7** Admin endpoints: list users, update role, delete user
- [ ] **2.8** Shared vaults with role-based access control
- [ ] **2.9** Audit logging for all auth events

---

## Phase 3: Advanced Features (~4-6 weeks)

_Depends on: Phase 2 complete_

- [ ] **3.1** Plugin system core: sandboxed JS runtime, permission model, manifest loader
- [ ] **3.2** Plugin API: registerTheme, registerCommand, registerView, registerExporter, registerImporter
- [ ] **3.3** Plugin lifecycle hooks: onNoteCreate/Update/Delete
- [ ] **3.4** Built-in plugins: Graph View, Kanban Board, Obsidian Importer, Notion Importer
- [ ] **3.5** Version history: `GET /api/v1/notes/:id/versions`, `POST /api/v1/notes/:id/restore`
- [ ] **3.6** Import: Obsidian vault, Notion export, Markdown folder
- [ ] **3.7** Export: Markdown, PDF, HTML, JSON
- [ ] **3.8** Attachments: images, PDFs, videos (filesystem + metadata in DB)
- [ ] **3.9** MCP server (optional, separate Go component)

---

## Phase 4: Mobile + Web (~4-6 weeks)

_Depends on: Phase 3 complete_

- [ ] **4.1** Android app (Kotlin): full feature parity with desktop, shares auth + sync
- [ ] **4.2** Web app deployment: deploy the binary to a server, users log in from any browser. Same auth, same sync server, same notes.
- [ ] **4.3** iOS app (future)

---

## Verification Commands

```sh
make check          # Full CI: lint + test for both frontend and backend
make build          # Verify single binary builds
./kvasir            # Smoke test the binary
golangci-lint run   # Go lint (backend)
go test ./...       # Go tests (backend)
pnpm lint           # ESLint (frontend)
pnpm test           # Vitest (frontend)
```

## Performance Targets

| Metric         | Target              |
| -------------- | ------------------- |
| Startup time   | < 2 seconds         |
| Search latency | < 100ms (10k notes) |
| Binary size    | < 100MB             |
| Memory usage   | < 500MB (10k notes) |
| Sync latency   | < 500ms             |
| Coverage       | ≥ 80% line (front+back) |

## Risks & Mitigations

| Risk                       | Mitigation                                   |
| -------------------------- | -------------------------------------------- |
| Sync complexity            | Start simple — version vectors + push/pull   |
| Single binary size         | Use `//go:embed`, strip symbols, UPX         |
| Performance with 10k notes | SQLite indexes, pagination, lazy loading     |
| Frontend bundle size       | Next.js code splitting, tree-shaking         |
| Self-hosting complexity    | Provide Docker container + setup guides      |
| Mobile later               | Keep APIs clean, sync logic separable        |

---

## Versioning

Agent-managed. After each commit, the agent appends a row to the changelog (here and in README.md)
and creates an annotated git tag. Patch (0.1.x) and minor (0.x.0) bumps are automated.
Major version bumps (x.0.0) are user-controlled — only do them when explicitly asked.

## Changelog

| Version | Date       | Changes                                                          |
| ------- | ---------- | ---------------------------------------------------------------- |
| 0.1.0   | 2026-06-26 | Initial project scaffold: agents.md, plan.md, README.md, LICENSE |
