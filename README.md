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

## License

[AGPL-3.0](LICENSE)

Copyright 2026 Kvasir Contributors
