.DEFAULT_GOAL := help

BINARY_NAME := lcli
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE        := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GO          := go
LDFLAGS     := -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: build run test test-cover lint vet fmt check clean install help

build: ## Build the binary
	$(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/$(BINARY_NAME)

install: build ## Install to GOPATH/bin
	cp bin/$(BINARY_NAME) $(shell $(GO) env GOPATH)/bin/$(BINARY_NAME)

test: ## Run tests with race detector
	$(GO) test -race ./...

test-cover: ## Run tests with coverage report
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run golangci-lint
	golangci-lint run ./...

vet: ## Run go vet
	$(GO) vet ./...

fmt: ## Format code
	gofmt -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w . || true

check: fmt vet test ## Run all quality checks

clean: ## Remove build artifacts
	rm -rf bin/ coverage.out coverage.html

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
