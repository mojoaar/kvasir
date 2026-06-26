.PHONY: dev dev-backend dev-frontend build build-frontend build-backend build-all lint lint-frontend lint-backend test test-frontend test-backend check clean

# ─── Development ───────────────────────────────────────────────

dev: dev-backend dev-frontend

dev-backend:
	cd backend && go run ./cmd/kvasir

dev-frontend:
	cd frontend && pnpm dev

# ─── Build ─────────────────────────────────────────────────────

build: build-frontend
	cd backend && go build -ldflags="-s -w" -o ../kvasir ./cmd/kvasir

build-frontend:
	cd frontend && pnpm build && rm -rf ../backend/internal/embed/dist && mkdir -p ../backend/internal/embed/dist && cp -r out/* ../backend/internal/embed/dist/

build-backend:
	cd backend && go build -ldflags="-s -w" -o ../kvasir ./cmd/kvasir

build-all: build-frontend
	cd backend && GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w" -o ../kvasir-darwin-amd64  ./cmd/kvasir
	cd backend && GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o ../kvasir-darwin-arm64  ./cmd/kvasir
	cd backend && GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o ../kvasir-linux-amd64   ./cmd/kvasir
	cd backend && GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w" -o ../kvasir-linux-arm64   ./cmd/kvasir
	cd backend && GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../kvasir-windows-amd64.exe ./cmd/kvasir

# ─── Lint ──────────────────────────────────────────────────────

lint: lint-frontend lint-backend

lint-frontend:
	cd frontend && pnpm lint

lint-backend:
	cd backend && golangci-lint run

# ─── Test ──────────────────────────────────────────────────────

test: test-frontend test-backend

test-frontend:
	cd frontend && pnpm test -- --coverage

test-backend:
	cd backend && go test -coverprofile=coverage.out ./...

# ─── CI ────────────────────────────────────────────────────────

check: lint test
	@echo "All checks passed."

# ─── Cleanup ───────────────────────────────────────────────────

clean:
	rm -rf frontend/.next frontend/out backend/internal/embed/dist
	rm -f kvasir kvasir-* coverage.out backend/coverage.out
	rm -rf frontend/node_modules/.cache
