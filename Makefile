# Adapted from https://www.thapaliya.com/en/writings/well-documented-makefiles/

BIN := syncsnipe
FRONTEND_DIR := frontend
FRONTEND_DIST := $(FRONTEND_DIR)/build
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
.PHONY: help install-deps build-frontend build-backend run-frontend run-backend build stuff format lint check-golangci-lint generate-sqlc push frontend-lint backend-lint format-frontend format-backend

##@ HELP & UTILS

help: ## Display help message with available targets
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf " \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ DEPENDENCIES

install-deps: $(STUFFBIN) ## Install dependencies for backend and frontend 
	@cd $(FRONTEND_DIR) && pnpm install --frozen-lockfile

$(STUFFBIN): ## Install stuffbin if missing
	@go install github.com/knadh/stuffbin@$(STUFFBIN_VERSION)

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

frontend-lint: ## Runs eslint for frontend 
	@cd $(FRONTEND_DIR) && pnpm lint || { echo "Frontend linting failed"; exit 1; }

format-frontend: ## Format frontend code
	@cd $(FRONTEND_DIR) && pnpm format

##@ BACKEND TASKS

build-backend: $(STUFFBIN) ## Build the backend binary
	@CGO_ENABLED=0 go build -a -ldflags="-s -w -X main.BuildString='$(BUILDSTR)'" -o $(BIN) cmd/syncsnipe/main.go

run-backend: ## Run the Go backend server in development mode
	@CGO_ENABLED=0 go run -ldflags="-s -w -X main.BuildString='$(BUILDSTR)'" cmd/syncsnipe/main.go

backend-lint: check-golangci-lint ## Runs golangci-lint for backend
	@golangci-lint run ./... || { echo "Backend linting failed"; exit 1; }

format-backend: ## Format backend code
	@go fmt ./...

##@ DATABASE MIGRATION & SQLC

generate-sqlc: ## Generate SQLC code
	@test -f sqlc.yaml || { echo "sqlc.yaml missing, you madlad"; exit 1; }
	docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc generate || { echo "SQLC gen failed, check your SQL"; exit 1; }

##@ BUILD & DEPLOYMENT

build: build-frontend build-backend stuff ## Build both frontend and backend, then bundle static assets
	@echo "→ Build complete, you legend."

stuff: $(STUFFBIN) ## Bundle static assets into binary using stuffbin
	@$(STUFFBIN) -a stuff -in $(BIN) -out $(BIN) $(FRONTEND_DIST)

##@ FORMATTING & LINTING & AUDITING

format: format-frontend format-backend ## Format entire workspace
	@echo "→ Formatting complete."

lint: frontend-lint backend-lint ## Run linting for both frontend and backend
	@echo "→ Linting complete."

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


##@ GIT ACTIONS

push: lint format ## Lint, format, and push code to Git
	@git push origin main
