#!/usr/bin/env bash
# check-commit-title.sh: fail unless a commit or PR title is a Conventional
# Commit. Used by .githooks/commit-msg and the PR title workflow.
set -euo pipefail

title="${1:-}"
pattern='^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([a-z0-9][a-z0-9._/-]*\))?!?: [^[:space:]]'

if [[ "$title" =~ $pattern ]]; then
  exit 0
fi

cat >&2 <<EOF
Not a Conventional Commit title: $title

Expected: <type>(<optional scope>): <summary>
Types: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test
Example: feat(skills): add guide to suggest the next skill
EOF
exit 1
