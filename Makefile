# Try to get the commit hash from 1) git 2) the VERSION file 3) fallback.
BIN := syncsnipe
FRONTEND_DIR := frontend
FRONTEND_DIST := ${FRONTEND_DIR}/build
STATIC := ${FRONTEND_DIST}
GOPATH ?= $(HOME)/go
STUFFBIN ?= $(GOPATH)/bin/stuffbin

VERSION := $(SYNCSNIPE_VERSION)
VERSION ?= $(shell git describe --tags --abbrev=0 2> /dev/null)
VERSION ?= $(shell grep -Eo 'tag: v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?' VERSION | cut -d' ' -f2)
VERSION ?= v0.0.0


BUILDSTR := ${VERSION} (\#${LAST_COMMIT} $(shell date -u +"%Y-%m-%dT%H:%M:%S%z"))

# The default target to run when `make` is executed.
.DEFAULT_GOAL := build

.PHONY: help install-deps build-frontend build-backend run-frontend run-backend build stuff

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

GOLANGCI_LINT_VERSION := 1.64.6

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


# Install stuffbin if it doesn't exist.
$(STUFFBIN):
	@echo "→ Installing stuffbin..."
	@go install github.com/knadh/stuffbin/...

install-deps: $(STUFFBIN) # Install dependencies for both backend and frontend.
	@echo "→ Installing frontend dependencies..."
	@cd ${FRONTEND_DIR} && pnpm install

frontend-build: install-deps # Build the frontend for production.
	@echo "→ Building frontend for production..."
	@export VITE_APP_VERSION="${VERSION}" && cd ${FRONTEND_DIR} && pnpm build

run-backend: # Run the Go backend server in development mode.
	@echo "→ Running backend..."
	CGO_ENABLED=0 go run -ldflags="-s -w" main.go

run-frontend: # Run the JS frontend server in development mode.
	@echo "→ Installing frontend dependencies (if not already installed)..."
	@cd ${FRONTEND_DIR} && pnpm install
	@echo "→ Running frontend..."
	@export VITE_APP_VERSION="${VERSION}" && cd ${FRONTEND_DIR} && pnpm dev

build-backend: $(STUFFBIN) # Build the backend binary.
	@echo "→ Building backend..."
	@CGO_ENABLED=0 go build -a \
		-ldflags="-s -w" \
		-o ${BIN} main.go

frontend-lint: # Runs eslint for frontend 
	@pnpm lint 

backend-lint: check-golangci-lint # Runs golangci-lint
	@golangci-lint run ./...

build: frontend-build build-backend stuff # Main build target: builds both frontend and backend, then stuffs static assets into the binary.
	@echo "→ Build successful."

stuff: $(STUFFBIN) # Stuff static assets into the binary using stuffbin.
	@echo "→ Stuffing static assets into binary..."
	@$(STUFFBIN) -a stuff -in ${BIN} -out ${BIN} ${STATIC}

