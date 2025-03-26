# Adapted from https://www.thapaliya.com/en/writings/well-documented-makefiles/

BIN := syncsnipe
FRONTEND_DIR := frontend
FRONTEND_DIST := $(FRONTEND_DIR)/build
STATIC := $(FRONTEND_DIST) sql
GOPATH ?= $(HOME)/go
STUFFBIN ?= $(GOPATH)/bin/stuffbin
GOLANGCI_LINT_VERSION ?= 1.64.6
STUFFBIN_VERSION ?= v1.3.0
LAST_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION ?= $(SYNCSNIPE_VERSION)
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "no-tag")
VERSION ?= $(shell grep -m1 '^v[0-9]+\.[0-9]+\.[0-9]+' VERSION 2>/dev/null || echo "v0.0.0")
BUILDSTR := $(VERSION) (\#$(LAST_COMMIT) $(shell date -u +"%FT%T%Z"))

.DEFAULT_GOAL := build
.PHONY: help install-deps build-frontend run-frontend lint-frontend test-frontend format-frontend build-backend run-backend lint-backend test-backend format-backend generate-sqlc build stuff format lint test audit clean-frontend clean-db clean-binary clean-all push

##@ HELP & UTILS

help: ## Display help message with available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf " \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ DEPENDENCIES

install-deps: ## Install all dependencies
	@echo "Installing dependencies..."
	@go mod download
	@cd $(FRONTEND_DIR) && pnpm install --frozen-lockfile

$(STUFFBIN): ## Install stuffbin if missing
	@echo "Installing stuffbin if missing..."
	@go install github.com/knadh/stuffbin/...

check-golangci-lint:
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint not found."; \
		exit 1; \
	fi; \
	installed_version=$$(golangci-lint version --format=short); \
	if [ "$$installed_version" != "$(GOLANGCI_LINT_VERSION)" ]; then \
		echo "Required golangci-lint version $(GOLANGCI_LINT_VERSION), but found $$installed_version."; \
		echo "Please install golangci-lint version $(GOLANGCI_LINT_VERSION) with the following command:"; \
		echo "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION)"; \
		exit 1; \
	fi

##@ FRONTEND TASKS

build-frontend: install-deps ## Build the frontend for production
	@cd $(FRONTEND_DIR) && VITE_APP_VERSION="$(VERSION)" pnpm build || { echo "Frontend build failed, check logs"; exit 1; }

run-frontend: ## Run the frontend development server
	@cd $(FRONTEND_DIR) && VITE_APP_VERSION="$(VERSION)" pnpm dev

lint-frontend: ## Runs eslint for frontend 
	@cd $(FRONTEND_DIR) && pnpm lint || { echo "Frontend linting failed"; exit 1; }

test-frontend: ## Run frontend tests
	@cd $(FRONTEND_DIR) && pnpm test

format-frontend: ## Format frontend code
	@cd $(FRONTEND_DIR) && pnpm format

##@ BACKEND TASKS

build-backend: $(STUFFBIN) ## Build the backend binary
	@go build -a -ldflags="-s -w" -o $(BIN) main.go

run-backend: ## Run the Go backend server in development mode
	@CGO_ENABLED=0 go run -ldflags="-s -w" main.go

lint-backend: check-golangci-lint ## Runs golangci-lint for backend
	@golangci-lint run ./... || { echo "Backend linting failed"; exit 1; }

test-backend: ## Run backend tests
	@go test -race -v -coverprofile=coverage.out ./...

format-backend: ## Format backend code
	@go fmt ./...

##@ DATABASE MIGRATION & SQLC

generate-sqlc: ## Generate SQLC code
	@test -f sqlc.yaml || { echo "sqlc.yaml missing, you madlad"; exit 1; }
	@docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc generate || { echo "SQLC gen failed, check your SQL"; exit 1; }

##@ BUILD & DEPLOYMENT

build: build-frontend build-backend stuff ## Build both frontend and backend, then bundle static assets
	@echo "→ Build complete, you legend."

stuff: $(STUFFBIN) ## Bundle static assets into binary using stuffbin
	@$(STUFFBIN) -a stuff -in $(BIN) -out $(BIN) $(STATIC)

##@ FORMATTING & LINTING & TESTING & AUDITING

format: format-frontend format-backend ## Format entire workspace
	@echo "→ Formatting complete."

lint: lint-frontend lint-backend ## Run linting for both frontend and backend
	@echo "→ Linting complete."

test: test-frontend test-backend ## Run tests for backend
	@echo "→ Testing complete."

audit: ## Run various code audits for security and best practices
	@echo "→ Running Go module verification..."
	@go mod verify
	@echo "→ Running Go vet..."
	@go vet ./...
	@echo "→ Running Staticcheck..."
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	@echo "→ Running vulnerability scan..."
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "→ Running tests with race detection..."
	@go test -race -buildvcs -vet=off ./...

##@ Clean Up

clean-frontend: ## Remove SvelteKit build and cache files
	@cd frontend && rm -rf build .svelte-kit || true

clean-db: ## Remove SQLite3 database and temporary files
	@rm -f syncsnipe.db syncsnipe.db-shm syncsnipe.db-wal

clean-binary: ## Remove the compiled Go binary
	@rm -f syncsnipe

clean-all: clean-frontend clean-db clean-binary ## Remove all generated artifacts
	@echo "→ Cleanup completed." 

##@ GIT ACTIONS


push: format lint ## Lint, format, and push code to Git
	@git diff --quiet || ( \
		echo "🚀 Formatting & Linting..."; \
		git add .; \
		git commit -m "🚀 Auto-format & lint: $(shell date +'%Y-%m-%d %H:%M:%S')"; \
		git branch --show-current | xargs -I {} git push origin {} \
	)
	@echo "→ Code formatted, linted, committed, and pushed."

