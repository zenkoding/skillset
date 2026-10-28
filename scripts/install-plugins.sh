#!/bin/bash

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo_info() {
  echo -e "${BLUE}==>${NC} $1"
}

echo_success() {
  echo -e "${GREEN}✓${NC} $1"
}

echo_warning() {
  echo -e "${YELLOW}!${NC} $1"
}

echo_error() {
  echo -e "${RED}✗${NC} $1"
}

# Skip if running inside Claude Code session
if [[ -n "$CLAUDECODE" ]]; then
  echo_warning "Running inside Claude Code session, skipping plugin installation"
  echo_info "Run this script manually outside of Claude Code to install plugins"
  exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$SCRIPT_DIR/../config"

echo_info "Claude Code Plugin & MCP Installation"
echo ""

# Check if claude CLI is available
if ! command -v claude &> /dev/null; then
  echo_error "Claude Code CLI is not available"
  echo_info "Install Claude Code from https://docs.claude.com/en/docs/claude-code/installation"
  exit 1
fi

echo_success "Claude Code CLI is available"

# Check if jq is available
if ! command -v jq &> /dev/null; then
  echo_error "jq is not installed (required to parse settings.json)"
  echo_info "Install jq: brew install jq (macOS) or apt-get install jq (Linux)"
  exit 1
fi

echo_success "jq is available"

# Check for node/npx (needed for cclsp and playwright)
if ! command -v node &> /dev/null || ! command -v npx &> /dev/null; then
  echo_warning "Node.js/npx is not available"
  echo_info "Some plugins (cclsp, playwright) require Node.js"
  echo_info "Install from https://nodejs.org/"
  echo_warning "Continuing without Node.js — cclsp MCP registration will be skipped if npx is unavailable"
else
  echo_success "Node.js/npx is available"
fi

echo ""

# --- Marketplaces ---

SETTINGS_JSON="$CONFIG_DIR/claude/settings.json"

echo_info "Checking plugin marketplaces..."
MARKETPLACE_LIST=$(claude plugin marketplace list 2>/dev/null || echo "")
MARKETPLACES_ADDED=0
MARKETPLACES_SKIPPED=0
ERRORS=0

while IFS= read -r name && IFS= read -r repo; do
  if echo "$MARKETPLACE_LIST" | grep -q "$repo"; then
    echo_success "$name marketplace already registered"
    ((MARKETPLACES_SKIPPED++))
  else
    echo_info "Adding $name marketplace ($repo)..."
    if claude plugin marketplace add "$repo"; then
      echo_success "$name marketplace registered"
      ((MARKETPLACES_ADDED++))
    else
      echo_error "Failed to add $name marketplace ($repo)"
      ((ERRORS++))
    fi
  fi
done < <(jq -r '.extraKnownMarketplaces | to_entries[] | (.key, .value.source.repo)' "$SETTINGS_JSON")

echo ""

# --- Plugins ---

echo_info "Checking plugins..."
PLUGIN_LIST=$(claude plugin list 2>/dev/null || echo "")
PLUGINS_INSTALLED=0
PLUGINS_SKIPPED=0

plugin_installed() {
  local plugin_name="${1%%@*}"  # extract name before @
  echo "$PLUGIN_LIST" | grep -q "^${plugin_name}@.*user"
}

while IFS= read -r plugin; do
  if plugin_installed "$plugin"; then
    echo_success "$plugin already installed"
    ((PLUGINS_SKIPPED++))
  else
    echo_info "Installing $plugin..."
    if claude plugin install "$plugin" -s user; then
      echo_success "$plugin installed"
      ((PLUGINS_INSTALLED++))
    else
      echo_error "Failed to install $plugin"
      ((ERRORS++))
    fi
  fi
done < <(jq -r '.enabledPlugins | keys[]' "$SETTINGS_JSON")

echo ""

# --- MCP Servers ---

echo_info "Checking MCP servers..."

# Parse MCP list (format varies, but contains server names)
MCP_LIST=$(claude mcp list 2>/dev/null || echo "")

MCP_ADDED=0
MCP_SKIPPED=0

# Check if cclsp is already configured
if echo "$MCP_LIST" | grep -q "cclsp"; then
  echo_success "cclsp MCP server already configured"
  ((MCP_SKIPPED++))
else
  echo_info "Adding cclsp MCP server..."

  # Set CCLSP_CONFIG_PATH to user config location
  CCLSP_CONFIG_PATH="$HOME/.config/claude/cclsp.json"

  if claude mcp add -s user -e "CCLSP_CONFIG_PATH=$CCLSP_CONFIG_PATH" cclsp -- npx cclsp@latest; then
    echo_success "cclsp MCP server configured"
    ((MCP_ADDED++))
  else
    echo_error "Failed to add cclsp MCP server"
  fi
fi

echo ""

# --- Agent Skills ---
#
# Skill installers that follow the Agent Skills layout place skills under ~/.agents/skills/.
# Claude Code loads skills from ~/.claude/skills/, so we symlink each agent skill
# directory into place. Re-running is safe: existing entries are left untouched.

echo_info "Syncing agent skills (~/.agents/skills → ~/.claude/skills)..."

AGENTS_SKILLS="$HOME/.agents/skills"
CLAUDE_SKILLS="$HOME/.claude/skills"
SKILLS_LINKED=0
SKILLS_SKIPPED=0
SKILLS_CONFLICT=0

if [[ ! -d "$AGENTS_SKILLS" ]]; then
  echo_warning "~/.agents/skills not found — skipping"
else
  mkdir -p "$CLAUDE_SKILLS"
  for skill_dir in "$AGENTS_SKILLS"/*/; do
    [[ -d "$skill_dir" ]] || continue
    name=$(basename "$skill_dir")
    # Repo-owned Codex workflows stay independent of Claude's skill catalog.
    if [[ -d "$CONFIG_DIR/codex/skills/$name" ]]; then
      SKILLS_SKIPPED=$((SKILLS_SKIPPED + 1))
      continue
    fi
    target="$CLAUDE_SKILLS/$name"
    if [[ -L "$target" ]]; then
      ((SKILLS_SKIPPED++))
    elif [[ -e "$target" ]]; then
      echo_warning "Conflict: ~/.claude/skills/$name already exists (not a symlink) — skipping"
      ((SKILLS_CONFLICT++))
    else
      ln -s "$skill_dir" "$target"
      echo_success "Linked agent skill: $name"
      ((SKILLS_LINKED++))
    fi
  done

  # Remove dangling symlinks (skills removed from ~/.agents/skills)
  for link in "$CLAUDE_SKILLS"/*/; do
    [[ -L "${link%/}" && ! -e "${link%/}" ]] || continue
    name=$(basename "$link")
    rm "${link%/}"
    echo_warning "Removed stale link: $name (no longer in ~/.agents/skills)"
  done

  if [[ $SKILLS_LINKED -eq 0 && $SKILLS_SKIPPED -gt 0 ]]; then
    echo_success "Agent skills already up to date ($SKILLS_SKIPPED linked)"
  fi
fi

echo ""

# --- Summary ---

echo ""
echo_info "Installation Summary:"
echo ""

if [[ $MARKETPLACES_ADDED -gt 0 ]]; then
  echo_success "Marketplaces added: $MARKETPLACES_ADDED"
fi
if [[ $MARKETPLACES_SKIPPED -gt 0 ]]; then
  echo "  Marketplaces already registered: $MARKETPLACES_SKIPPED"
fi

if [[ $PLUGINS_INSTALLED -gt 0 ]]; then
  echo_success "Plugins installed: $PLUGINS_INSTALLED"
fi
if [[ $PLUGINS_SKIPPED -gt 0 ]]; then
  echo "  Plugins already installed: $PLUGINS_SKIPPED"
fi

if [[ $MCP_ADDED -gt 0 ]]; then
  echo_success "MCP servers added: $MCP_ADDED"
fi
if [[ $MCP_SKIPPED -gt 0 ]]; then
  echo "  MCP servers already configured: $MCP_SKIPPED"
fi

if [[ $SKILLS_LINKED -gt 0 ]]; then
  echo_success "Agent skills linked: $SKILLS_LINKED"
fi
if [[ $SKILLS_CONFLICT -gt 0 ]]; then
  echo_warning "Agent skill conflicts (manual resolution needed): $SKILLS_CONFLICT"
fi

echo ""
echo_success "Plugin installation complete!"
echo ""
echo "Verification commands:"
echo "  claude plugin marketplace list  # View registered marketplaces"
echo "  claude plugin list              # View installed plugins"
echo "  claude mcp list                 # View MCP servers"
echo ""
