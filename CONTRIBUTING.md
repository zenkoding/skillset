# Contributing to skillset

skillset builds an agent setup for Claude Code and Codex from composable sets of instructions, skills and settings. A single maintainer develops it, and work is tracked in public GitHub issues.

## Start here

```bash
git clone https://github.com/zenkoding/skillset.git && cd skillset
make git-setup                          # use the repo's commit template and commit-msg hook
make install                            # install the binary from this checkout
skillset init "$PWD"                    # pick hosts and sets from this checkout (optional: default is everything)
skillset --catalog "$PWD/config" plan   # preview
skillset --catalog "$PWD/config" apply  # install from the working tree
make validate                           # lint, build, test, doc-rot guard
```

See [README.md](README.md) for sets and what apply installs.

## Dev docs

| Doc | Read it when |
| --- | --- |
| [Local development](docs/contributing/local-dev.md) | running the installer from your checkout |
| [Testing](docs/contributing/testing.md) | running or writing tests |
| [Issues](docs/contributing/issues.md) | filing, triaging or picking up work |
| [Docs](docs/contributing/docs.md) | writing or finding a plan or other doc |

## Layout

Config lives under `config/`, the catalog. Code lives everywhere else.

| Path | What it holds |
|------|---------------|
| `config/sets/<name>/` | Sets (`set.toml` manifests). |
| `config/skills/<set>/<name>/` | Skills, one folder per set. `config/skills/in-progress/` holds drafts that are not installed. |
| `config/instructions/` | `shared.md` plus `claude.md` / `codex.md`, and one file per set that has its own (`engineering.md`, `product.md`, `productivity.md`). |
| `config/shared/` | References that several skills link. The installer copies each one into the skills that link it. |
| `config/claude/` | Claude settings: `settings.json` (reasons in `settings.md`), `statusline.sh`, `hooks/`, output styles. |
| `config/codex/` | Codex defaults, merged by hand. |
| `config/dcg/` | destructive_command_guard config and packs. |
| `cmd/skillset`, `internal/` | The installer: entry points in `cmd/`, packages in `internal/`. |
| `scripts/` | Plugin install, doc-rot check and commit title check. |

## Making a change

1. Find work in [GitHub issues](docs/contributing/issues.md).
2. Branch from `main`, implement, and run `make validate`.
3. Open one PR per issue; the PR closes its issue.

- Commit messages and PR titles follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/): `<type>(<optional scope>): <summary>`, with `!` before the colon for a breaking change. The types are `build`, `chore`, `ci`, `docs`, `feat`, `fix`, `perf`, `refactor`, `revert`, `style` and `test`.
- A commit body has a paragraph on what changed and a paragraph on why.
- `make git-setup` configures your clone once: `commit.template` points at `.gitmessage` in your main checkout, and `core.hooksPath` points at `.githooks` there, whose `commit-msg` hook rejects a title in the wrong format. The paths are absolute so that every worktree and older branch uses the same files; run it again if you move the clone. Setting `core.hooksPath` means Git no longer runs any hooks you keep in `.git/hooks` for this clone.
- CI checks the PR title in the same format.
- PRs merge by squash. The PR title and body become the commit message, so the body's **Why** records any decision the change makes: "Chose X over Y because Z".
- Fill in every section of the [PR template](.github/pull_request_template.md).
- Go code uses the standard library first.

## Adding a skill

1. Create `config/skills/<set>/<name>/SKILL.md` with `name: <name>` in its frontmatter (drafts go in `config/skills/in-progress/`).
2. To name a skill from a set yours doesn't require, write "`<skill>` when installed"; `make test` checks this. A file that only names skills conditionally (like guide's workflow map) can say so once with `<!-- every skill named here: when installed -->`.
3. If the model should reach for it on its own at a moment its description can't catch (an event, a step in a loop), add one line to its set's instructions file. Skip skills only a person runs.
4. Run `skillset --catalog "$PWD/config" plan`, then the same command with `apply`.
