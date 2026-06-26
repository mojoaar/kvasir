.PHONY: dev dev-backend dev-frontend build build-frontend build-backend build-all lint lint-frontend lint-backend test test-frontend test-backend check check-coverage-frontend check-coverage-backend clean

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
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend && golangci-lint run; \
	else \
		echo "golangci-lint not found — falling back to go vet"; \
		cd backend && go vet ./...; \
	fi

# ─── Test ──────────────────────────────────────────────────────

test: test-frontend test-backend

test-frontend:
	cd frontend && pnpm exec vitest run --coverage

test-backend:
	cd backend && go test -coverprofile=coverage.out -coverpkg=./internal/storage/...,./internal/api/handlers/... ./...

# ─── Coverage Thresholds (≥ 80%) ───────────────────────────────

COVERAGE_MIN ?= 80

check-coverage-backend:
	@cd backend && \
	coverage=$$(go tool cover -func=coverage.out 2>/dev/null | grep total | awk '{print substr($$3, 1, length($$3)-1)}' | tr -d '\n'); \
	if [ -z "$$coverage" ]; then \
		echo "ERROR: Could not parse Go coverage. Did you run 'make test-backend' first?"; \
		exit 1; \
	fi; \
	if [ "$$(echo "$$coverage < $(COVERAGE_MIN)" | bc)" -eq 1 ]; then \
		echo "FAIL: Go coverage $$coverage% is below $(COVERAGE_MIN)% threshold"; \
		exit 1; \
	fi; \
	echo "OK: Go coverage $$coverage% >= $(COVERAGE_MIN)%"

check-coverage-frontend:
	@cd frontend && \
	coverage=$$(node -e "try{var c=require('./coverage/coverage-summary.json');var p=c.total.lines.pct;process.stdout.write(String(p))}catch(e){process.stderr.write('ERROR\\n');process.exit(1)}" 2>&1); \
	if [ "$$coverage" = "ERROR" ]; then \
		echo "ERROR: Could not parse frontend coverage. Did you run 'make test-frontend' first?"; \
		exit 1; \
	fi; \
	if [ "$$(echo "$$coverage < $(COVERAGE_MIN)" | bc)" -eq 1 ]; then \
		echo "FAIL: Frontend coverage $$coverage% is below $(COVERAGE_MIN)% threshold"; \
		exit 1; \
	fi; \
	echo "OK: Frontend coverage $$coverage% >= $(COVERAGE_MIN)%"

# ─── CI ────────────────────────────────────────────────────────

check: lint test check-coverage-backend check-coverage-frontend
	@echo ""
	@echo "All checks passed: lint, test, coverage ≥ $(COVERAGE_MIN)%"

# ─── Cleanup ───────────────────────────────────────────────────

clean:
	rm -rf frontend/.next frontend/out backend/internal/embed/dist
	rm -f kvasir kvasir-* coverage.out backend/coverage.out
	rm -rf frontend/node_modules/.cache frontend/coverage
