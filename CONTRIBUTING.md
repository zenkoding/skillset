# Contributing to skillset

## Setup

```bash
git clone <repo-url> && cd skillset
make init       # pick hosts and sets (optional: default is everything)
make plan       # preview
make apply      # install from the working tree
make validate   # lint, build, test, doc-rot guard
```

See [README.md](README.md) for sets and what apply installs.

## Layout

Config lives under `config/`, the catalog. Code lives everywhere else.

| Path | What it holds |
|------|---------------|
| `config/sets/<name>/` | Sets (`set.toml` manifests). |
| `config/skills/<category>/<name>/` | Skills. `config/skills/in-progress/` holds drafts that are not installed. |
| `config/instructions/` | `shared.md` plus `claude.md` / `codex.md`. |
| `config/claude/` | Claude settings: `settings.json` (reasons in `settings.md`), `statusline.sh`, `hooks/`, output styles. |
| `config/codex/` | Codex defaults, merged by hand. |
| `config/dcg/` | destructive_command_guard config and packs. |
| `cmd/skillset`, `internal/` | The installer. |
| `scripts/` | Plugin install and doc-rot check. |

## Workflow

1. Find work in GitHub issues (`gh issue list`).
2. Branch, implement following [AGENTS.md](AGENTS.md), and run `make validate`.
3. Open one PR per issue; the PR closes its issue.

Commit format: title / what / why.

## Adding a skill

1. Create `config/skills/<category>/<name>/SKILL.md` with `name: <name>` in its frontmatter (drafts go in `config/skills/in-progress/`).
2. Run `make plan`, then `make apply`.
