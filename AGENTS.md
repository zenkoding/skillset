# Agent Instructions: skillset

At session start, invoke `i-have-adhd` before the first substantive response.
If this host has no skill invocation tool, read the installed `i-have-adhd/SKILL.md`
and apply it. Keep it active throughout the session, including after compaction,
unless the user says "stop adhd mode" or "normal mode".

## What this repo is

Personal agent setup for Claude Code and Codex: user-level instructions,
settings, and skills, installed into `~/.claude`, `~/.codex` and `~/.agents`.

Layout: code and config are split. `config/` is the catalog the installer
reads: `config/skills/` the single skill source, `config/instructions/` the
shared instructions source, `config/sets/` the sets (`legacy` wraps
everything), `config/shared/` references several skills carry as synced copies
(`make sync/shared`), `config/claude/`, `config/codex/` and `config/dcg/` the host
settings. Code: `cmd/skillset` is the installer for both hosts, packages under
`internal/`, helper scripts under `scripts/`.

## Working here

- Develop with Matt Pocock's skills (`mattpocock-skills` plugin) for now.
  Track work in GitHub issues, not beads.
- Preview with `make plan` before installing anything into `$HOME`.
- Go code: `cmd/` for entry points, `internal/` for packages, table-driven
  tests, standard library first. Run `make test` before committing.
- Commit format: title / what / why.

## Agent skills

### Issue tracker

GitHub Issues on `zenkoding/skillset`, via `gh`. See `docs/agents/issue-tracker.md`.

### Triage labels

Matt's five defaults (`needs-triage`, `needs-info`, `ready-for-agent`,
`ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `CONTEXT.md` and `docs/adr/` at the repo root. See `docs/agents/domain.md`.
