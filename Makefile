# Local development tasks for skillset: install the binary from this checkout,
# and build, lint and test the code. The binary itself installs and manages
# sets; end users start from the README.

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Colors for output. Real escape bytes, because macOS bash's builtin echo
# prints \033 literally.
BLUE := $(shell printf '\033[0;34m')
GREEN := $(shell printf '\033[0;32m')
YELLOW := $(shell printf '\033[1;33m')
NC := $(shell printf '\033[0m')

.PHONY: help install check-docs build lint test test/go test/json validate

# Sandboxed tool binaries (version-pinned, never pollutes system PATH)
BIN_DIR := target
GOLANGCI_LINT_VERSION := 2.15.0
GOLANGCI_LINT_BINARY  := $(BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)

$(GOLANGCI_LINT_BINARY):
	@mkdir -p $(BIN_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v$(GOLANGCI_LINT_VERSION)/install.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)
	@mv $(BIN_DIR)/golangci-lint $@

help: ## Show this help message
	@echo "$(BLUE)skillset - Available Targets$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z/_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "Local development:"
	@echo "  $(YELLOW)Install binary:$(NC)  make install"
	@echo "  $(YELLOW)Try the catalog:$(NC) skillset --catalog \"\$$PWD/config\" plan, then apply"
	@echo "  $(YELLOW)Before a PR:$(NC)     make validate"
	@echo ""

install: ## Install the skillset binary from this checkout
	@go install ./cmd/skillset
	@gobin="$$(go env GOBIN)"; echo "$(GREEN)✓$(NC) Installed $${gobin:-$$(go env GOPATH)/bin}/skillset"

build: ## Build all Go packages
	@echo "$(BLUE)Building Go packages...$(NC)"
	@go build ./...
	@echo "$(GREEN)✓$(NC) Go build passed"

lint: $(GOLANGCI_LINT_BINARY) ## Run golangci-lint (sandboxed install)
	@echo "$(BLUE)Running golangci-lint...$(NC)"
	@$(GOLANGCI_LINT_BINARY) run ./...
	@echo "$(GREEN)✓$(NC) Lint passed"

test: test/go test/json ## Run all tests

test/go: ## Run all Go tests
	@echo "$(BLUE)Running Go tests...$(NC)"
	@go test ./...
	@echo "$(GREEN)✓$(NC) Go tests passed"

test/json: ## Validate all JSON config files
	@echo "$(BLUE)Validating JSON configs...$(NC)"
	@failed=0; \
	for f in $$(find config/ -name '*.json' -type f 2>/dev/null); do \
		if ! jq empty "$$f" 2>/dev/null; then \
			echo "Error: Invalid JSON: $$f"; \
			failed=1; \
		fi; \
	done; \
	if [ "$$failed" -eq 1 ]; then exit 1; fi
	@echo "$(GREEN)✓$(NC) JSON validation passed"

check-docs: ## Verify harness docs are free of drift-prone counts/stamps
	@./scripts/check-doc-rot.sh

validate: lint build test check-docs ## Run full verification pipeline (lint + build + test + doc-rot guard)
	@echo "$(GREEN)✓ All checks passed$(NC)"
