#!/bin/bash
# SessionStart hook: print the current repo's project knowledge catalogs so
# they land in the session context (config/instructions/shared.md, "Project knowledge").
#
#   private  ~/.local/share/agents/projects/<key>/README.md  outside the repo, one per main checkout
#   public   <worktree>/docs/agents/README.md                committed, per branch
#
# <key> is the origin URL normalized to host/owner/repo and slugified, so all
# clones and worktrees of a repo share one store. Without an origin it is
# "local-" plus the slugified main checkout path.
# The store lives outside the repo because worktree sandboxes block writes
# into the main checkout; ~/.local/share/agents is granted once in settings.
# AGENTS_DATA_HOME overrides that root (tests).
#
# Prints nothing outside a git repo or when neither catalog exists. Always
# exits 0 so a broken repo never blocks a session.

dir=${CLAUDE_PROJECT_DIR:-$PWD}
top=$(git -C "$dir" rev-parse --show-toplevel 2>/dev/null) || exit 0
main=$top
common=$(git -C "$dir" rev-parse --path-format=absolute --git-common-dir 2>/dev/null)
if [ "$(basename "$common")" = ".git" ]; then
	main=$(dirname "$common")
fi

slug() { printf '%s' "$1" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+|-+$//g'; }
norm() { printf '%s' "$1" | sed -E 's#^[a-z+]+://##; s#^[^@/]+@##; s#:([^0-9/])#/\1#; s#:[0-9]+/#/#; s#\.git/?$##; s#/+$##'; }
projects=${AGENTS_DATA_HOME:-$HOME/.local/share/agents}/projects
local=local-$(slug "$main")
if url=$(git -C "$dir" remote get-url origin 2>/dev/null) && [ -n "$url" ]; then
	store=$projects/$(slug "$(norm "$url")")
else
	store=$projects/$local
fi
private=$store/README.md
public=$top/docs/agents/README.md
[ -f "$private" ] || [ -f "$public" ] || exit 0

max=200 # lines per catalog; a catalog is an index, not a doc
print() { # label path
	[ -f "$2" ] || return 0
	printf '\n## %s: %s\n\n' "$1" "$2"
	head -n "$max" "$2"
	n=$(wc -l <"$2")
	[ "$n" -le "$max" ] || printf '\n[cut at %d of %d lines: trim it with the project-knowledge skill]\n' "$max" "$n"
}

echo "# Project knowledge catalogs"
echo
echo "Private store: $store (write private docs and notes here, by this absolute path)."
echo "Open the docs the task needs; paths are relative to each catalog's directory."
echo "Precedence: private > public > user. Write policy: \"Project knowledge\" in your instructions."
if [ "$store" != "$projects/$local" ] && [ -d "$projects/$local" ]; then
	echo "Note: $projects/$local is this checkout's store from before it had an origin; move its docs into the store above."
fi
print "Private (edit freely)" "$private"
print "Public (committed; propose changes)" "$public"
exit 0
