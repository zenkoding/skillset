# Development tasks for skillset: run the installer from this checkout, and
# build, lint and test the code.

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Colors for output. Real escape bytes, because macOS bash's builtin echo
# prints \033 literally.
BLUE := $(shell printf '\033[0;34m')
GREEN := $(shell printf '\033[0;32m')
YELLOW := $(shell printf '\033[1;33m')
NC := $(shell printf '\033[0m')

.PHONY: help init plan apply install uninstall uninstall/plan bin install/plugins install/tools install/dcg check-docs build lint test test/go test/json validate

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
	@echo "Common workflows (catalog: this checkout's config/):"
	@echo "  $(YELLOW)First time:$(NC)      make init, make plan, then make apply"
	@echo "  $(YELLOW)After edits:$(NC)     make apply"
	@echo "  $(YELLOW)Remove:$(NC)          make uninstall/plan, then make uninstall"
	@echo ""

# The installer, run from this checkout against config/ so links point into
# the working tree.
SKILLSET := go run ./cmd/skillset --catalog "$(CURDIR)/config"

init: ## Pick hosts and sets from this checkout's catalog
	@$(SKILLSET) init

plan: ## Preview what apply would change in ~/.claude, ~/.codex, ~/.agents and ~/.config/dcg
	@$(SKILLSET) plan

apply: ## Install the chosen sets from the working tree
	@$(SKILLSET) apply

install: apply ## Alias for apply

uninstall/plan: ## Preview what uninstall would remove
	@$(SKILLSET) uninstall --dry-run

uninstall: ## Remove the links and unedited settings copies apply installed
	@$(SKILLSET) uninstall

bin: ## Install the skillset binary into GOPATH/bin
	@go install ./cmd/skillset
	@echo "$(GREEN)✓$(NC) Installed $(shell go env GOPATH)/bin/skillset"

install/plugins: ## Install Claude Code plugins, marketplaces, and MCP servers
	@./scripts/install-plugins.sh

# Tool install targets

install/tools: install/dcg ## Install optional tools

install/dcg: ## Install the destructive_command_guard (dcg) binary
	@if command -v dcg >/dev/null 2>&1; then \
		echo "$(GREEN)✓$(NC) dcg already installed ($$(dcg --version 2>&1 | head -1))"; \
	else \
		echo "$(BLUE)Installing dcg...$(NC)"; \
		curl -fsSL "https://raw.githubusercontent.com/Dicklesworthstone/destructive_command_guard/main/install.sh" | bash -s -- --no-configure --easy-mode; \
	fi
	@echo "dcg config and packs are linked by make apply"

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
