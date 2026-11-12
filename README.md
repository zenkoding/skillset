# skillset

Build your own agent setup for supported agents from composable **sets**: instructions, skills and settings. Runtime state and credentials stay local.

## Quick start

```sh
go install github.com/zenkoding/skillset/cmd/skillset@latest

skillset init        # pick a source (default zenkoding/skillset@main), hosts and sets
skillset plan        # preview what changes in your home directory
skillset apply       # install
skillset update      # move to the newest release tag, then apply
skillset uninstall   # remove what apply installed (--dry-run to preview)
```

`init` takes a source: a git URL or `owner/repo`, with an optional `@ref` (a branch, a tag or a full commit hash), or a local directory. It writes `~/.config/skillset/config.toml` with the source, hosts and sets, and warns when dcg is missing. `plan`, `apply` and `update` do not read the source yet: they install from the release pin or `--catalog`.

`skillset update --ref <tag>` switches back to an older release, and `--ref main` follows main.

Optional extras: from a clone of this repo, `./scripts/install-plugins.sh` installs the Claude plugins, marketplaces and MCP servers listed in `config/claude/settings.json`. The dcg binary installs separately; `skillset apply` links its config and packs:

```sh
curl -fsSL "https://raw.githubusercontent.com/Dicklesworthstone/destructive_command_guard/main/install.sh" | bash -s -- --no-configure --easy-mode
```

## How it works

```
~/.config/skillset/config.toml  (hosts + sets) ─┐
~/.config/skillset/lock.toml    (pinned tag)  ──┼─> fetch this repo @tag -> ~/.cache/skillset/<tag>/
                                                │
chosen sets in dependency order + ~/.config/skillset/local/ (always last)
  ├─ instructions: shared.md + <host>.md  -> ~/.claude/CLAUDE.md, ~/.codex/AGENTS.md
  ├─ skills, files: symlinks into the cached catalog
  └─ settings.json: merge fragments -> 3-way merge with your copy -> ~/.claude/settings.json
```

- **Catalog**: the `config/` directory of this repo. `apply` and `update` fetch it at a release tag; the previous tag stays cached for rollback.
- **Sets** are directories under `config/sets/<name>/` with a `set.toml`. A set that leaves out the instructions for a host is skipped for that host.
- **Local set**: personal additions go in `~/.config/skillset/local/`, in the same format (`shared.md`, `claude.md`, `codex.md`, `skills/`, optional `set.toml`). It is applied last.
- Without a saved config, every set is installed for both hosts.

### Sets

```
config/sets/<name>/
  set.toml        name, description, hosts, requires, paths
  shared.md       instruction fragment for every host
  claude.md       optional, Claude only
  codex.md        optional, Codex only
  settings.json   Claude settings fragment
  skills/...
```

```toml
# example
name = "dcg"
description = "Destructive-command guard for shell permissions"
hosts = ["claude", "codex"]
requires = ["base"]

[[files]]                       # extra files linked into HOME
home = ".config/dcg/config.toml"
path = "config.toml"
```

Paths default to the layout above and must stay inside the catalog. Each set in the catalog points at its folder under `config/skills/`, with requires that point one way:

```
base  <-  development  <-  product
  ^                    <-  triage
  |                    <-  architecture
  |                    <-  quality
  └── productivity
```

| Set | Holds |
| --- | --- |
| `base` | shared, Claude and Codex instructions, Claude settings, statusline, the project knowledge hook, dcg config; skills brief, grill, debrief, research, project-knowledge, setup-repo, guide |
| `development` | when-to-use instructions for its skills; the build loop: build, tdd, publish, review, test, address, codify |
| `triage` | triage-prs |
| `architecture` | when-to-use instructions for its skills; optimize |
| `quality` | when-to-use instructions for its skills; fix-bug |
| `product` | when-to-use instructions for its skills; envision, scope-milestone, shape-epic, retrospect, build-poc, create-rfc, update-rfc |
| `productivity` | i-have-adhd and the session-start instruction that invokes it; the writing style for docs people read |

A skill names only skills in its own set or a set it requires; a link into any other set says "when installed" right after the name (a test enforces this). The same rule holds for each set's instructions.

Each set's instructions say when to reach for its skills on their own: events and loop order only, since skill descriptions already say what each skill does. A skill a person runs by hand stays out of them.

### What apply does

- Each skill directory is linked as `<name>` into `~/.claude/skills`, and for Codex into `~/.codex/skills` and `~/.agents/skills`. Drafts in `in-progress/`, skills with an empty `SKILL.md` or no `name:`, and skills whose `name:` differs from the directory are skipped.
- Instructions are rendered into the catalog's gitignored `build/` and linked from `~/.claude/CLAUDE.md` and `~/.codex/AGENTS.md`.
- `~/.claude/settings.json` is copied, because Claude Code writes to it. Fragments from several sets merge first: string lists as sets, any other key set by two sets is an error naming both. Later installs merge key by key against the last install: a key only one side changed takes that change, a key both sides changed is a conflict. In a terminal `apply` asks which side wins; otherwise conflicts stop the run. Merge bases live in `~/.local/state/skillset/copies/`.
- Links are recorded in `~/.local/state/skillset/installed-links.json`. Recorded links that are no longer wanted are removed, only when they point into the catalog. Anything else in the way is a clash: nothing changes and every clash is listed. A real `~/.claude/CLAUDE.md` is a clash too: move it aside and rerun.
- `.docs/` is added once to git's global excludes file, so a `<repo>/.docs` link to a private project knowledge store is never committed.
- Codex settings are not installed: merge missing keys from `config/codex/config.toml` into `~/.codex/config.toml` by hand.

## Developing

The Makefile is for local development. `make install` puts a binary built from this checkout on your Go bin path. Pass `--catalog` so it reads this checkout's `config/` instead of the pinned release: links point into the working tree and edits are live.

```sh
make install                                 # install the binary from this checkout
skillset init "$PWD"                         # pick hosts and sets from this checkout
skillset --catalog "$PWD/config" plan        # preview
skillset --catalog "$PWD/config" apply       # install from the working tree
skillset --catalog "$PWD/config" uninstall   # remove it (--dry-run to preview)
make validate                                # lint, build, test, doc-rot guard
make help                                    # everything else
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow.
