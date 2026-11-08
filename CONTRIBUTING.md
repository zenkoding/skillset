# Contributing to skillset

## Setup

```bash
git clone <repo-url> && cd skillset
make install                            # install the binary from this checkout
skillset --catalog "$PWD/config" init   # pick hosts and sets (optional: default is everything)
skillset --catalog "$PWD/config" plan   # preview
skillset --catalog "$PWD/config" apply  # install from the working tree
make validate                           # lint, build, test, doc-rot guard
```

See [README.md](README.md) for sets and what apply installs.

## Layout

Config lives under `config/`, the catalog. Code lives everywhere else.

| Path | What it holds |
|------|---------------|
| `config/sets/<name>/` | Sets (`set.toml` manifests). |
| `config/skills/<set>/<name>/` | Skills, one folder per set. `config/skills/in-progress/` holds drafts that are not installed. |
| `config/instructions/` | `shared.md` plus `claude.md` / `codex.md`, and one file per set that has its own (`engineering.md`, `product.md`, `productivity.md`). |
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

1. Create `config/skills/<set>/<name>/SKILL.md` with `name: <name>` in its frontmatter (drafts go in `config/skills/in-progress/`).
2. To name a skill from a set yours doesn't require, write "`<skill>` when installed"; `make test` checks this. A file that only names skills conditionally (like guide's workflow map) can say so once with `<!-- every skill named here: when installed -->`.
3. If the model should reach for it on its own at a moment its description can't catch (an event, a step in a loop), add one line to its set's instructions file. Skip skills only a person runs.
4. Run `skillset --catalog "$PWD/config" plan`, then the same command with `apply`.
