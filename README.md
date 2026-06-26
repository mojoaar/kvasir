# Kvasir

> *"The kernel of your mind"*

A beautiful, techy, Nordic-inspired markdown knowledge base. Sync-first, API-first, plugin-extensible, and built for the thinking mind.

[![License](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)](https://go.dev)
[![Node Version](https://img.shields.io/badge/node-22+-339933?logo=node.js)](https://nodejs.org)

## Features

- **Markdown Editor** — TipTap with tables, footnotes, math (Katex), and Mermaid diagrams
- **5 Themes** — Kvasir (Nordic), Dracula, Nord, GitHub, Cyberpunk (each with dark + light mode)
- **Full-Text Search** — Instant search via SQLite FTS5
- **Command Palette** — Cmd+K / Ctrl+K for everything
- **REST API** — Full OpenAPI/Swagger documentation
- **Single Binary** — Go embeds the frontend, runs everywhere
- **Cross-Platform** — Desktop, Android, and web. One account, all your notes on every device
- **CLI** — List, search, export notes from the terminal
- **Plugin System** — Extensible with sandboxed plugins
- **Offline-First** — Local SQLite storage, sync when connected

## Quick Start

### Download Binary

Download the latest release from [releases](https://github.com/mojoaar/kvasir/releases) and run:

```sh
./kvasir
# Open http://localhost:8080
```

### Build from Source

```sh
git clone git@github.com:mojoaar/kvasir.git
cd kvasir
make build
./kvasir
```

## Development

```sh
make dev     # Start frontend (:3000) + backend (:8080)
make build   # Build single binary
make check   # Lint + test
```

## Tech Stack

| Layer    | Tech                                  |
| -------- | ------------------------------------- |
| Frontend | Next.js 16, shadcn/ui, Tailwind CSS   |
| Backend  | Go 1.26, Gin                          |
| Database | SQLite (modernc.org, pure Go, no CGO) |
| Editor   | TipTap + Katex + Mermaid              |
| Auth     | JWT + OAuth2                          |
| Sync     | Custom version vectors                |

## Themes

| Theme     | Dark Palette     | Light Palette  |
| --------- | ---------------- | -------------- |
| Kvasir    | Nordic dark      | Nordic snow    |
| Dracula   | Purple dark      | Warm light     |
| Nord      | Polar night      | Snow storm     |
| GitHub    | Dark dimmed      | Default light  |
| Cyberpunk | Neon black       | Neon white     |

Toggle dark/light per theme via the header switcher or `Cmd+Shift+T`.

## Changelog

| Version | Date       | Changes                                                          |
| ------- | ---------- | ---------------------------------------------------------------- |
| 0.1.0   | 2026-06-26 | Initial project scaffold: agents.md, plan.md, README.md, LICENSE |
| 0.1.1   | 2026-06-26 | Phase 0.1: Makefile, backend (Gin), frontend (Next.js 16), CLI (Cobra), configs |
| 0.1.2   | 2026-06-26 | Phase 0.2: SQLite schema, storage layer (87.5% coverage), FTS5 search, health endpoint, SQLTime type |
| 0.1.3   | 2026-06-26 | Phase 0.3: shadcn/ui, TanStack React Query, Zustand stores, Providers, Sidebar component |
| 0.1.4   | 2026-06-26 | Phase 0.4: 5-theme system (CSS variables), dark/light toggle, theme selector, Settings page |
| 0.1.5   | 2026-06-26 | Phase 0.5: TipTap editor, toolbar, split view, Katex math, Mermaid diagrams, auto-save |
| 0.1.6   | 2026-06-26 | Phase 0.6: Sidebar file tree, drag-and-drop, inline rename, folder CRUD |
| 0.1.7   | 2026-06-26 | Phase 0.7: Notes REST API (CRUD), Zod schemas, handler+storage tests (83 total) |
| 0.1.8   | 2026-06-26 | Phase 0.8: Full-text search endpoint, tag search, sidebar input, results page |
| 0.1.9   | 2026-06-26 | Phase 0.9: Tags API + UI, tag chips, color palette, note tagging |
| 0.1.10  | 2026-06-26 | Phase 0.10: Command palette (cmdk), Cmd+K shortcut |
| 0.1.11  | 2026-06-26 | Phase 0.11: Swaggo annotations, Swagger UI, OpenAPI spec |
| 0.1.12  | 2026-06-26 | Phase 0.12: CLI — list, search, export, create commands |
| 0.1.13  | 2026-06-26 | Phase 0.13: Single binary build, //go:embed frontend, 35MB binary |
| 0.1.14  | 2026-06-26 | Phase 0.14: Welcome note + tutorial, first-run seed |
| 0.1.15  | 2026-06-26 | Phase 0.15: Testing (35 frontend + backend), 80%+ coverage, MVP complete |

## License

[AGPL-3.0](LICENSE)

Copyright 2026 Kvasir Contributors
