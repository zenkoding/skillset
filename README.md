# skillset

Build your own agent setup for supported agents from composable **sets**: instructions, skills and settings. Runtime state and credentials stay local.

## Quick start

```sh
go install github.com/zenkoding/skillset/cmd/skillset@latest

skillset init        # pick hosts (Claude, Codex) and sets; required sets are added
skillset plan        # preview what changes in your home directory
skillset apply       # install
skillset update      # move to the newest release tag, then apply
skillset uninstall   # remove what apply installed (--dry-run to preview)
```

`skillset update --ref <tag>` switches back to an older release, and `--ref main` follows main.

For installation you can also clone the repo and use the [Developing](#developing) flow.

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

Paths default to the layout above and must stay inside the catalog. Today the whole setup ships as one `legacy` set (`config/sets/legacy/`) that points at `config/skills/`, `config/instructions/` and `config/claude/`; splitting it into smaller sets is follow-up work.

### What apply does

- Each skill directory is linked as `<name>` into `~/.claude/skills`, and for Codex into `~/.codex/skills` and `~/.agents/skills`. Drafts in `in-progress/`, skills with an empty `SKILL.md` or no `name:`, and skills whose `name:` differs from the directory are skipped.
- Instructions are rendered into the catalog's gitignored `build/` and linked from `~/.claude/CLAUDE.md` and `~/.codex/AGENTS.md`.
- `~/.claude/settings.json` is copied, because Claude Code writes to it. Fragments from several sets merge first: string lists as sets, any other key set by two sets is an error naming both. Later installs merge key by key against the last install: a key only one side changed takes that change, a key both sides changed is a conflict. In a terminal `apply` asks which side wins; otherwise conflicts stop the run. Merge bases live in `~/.local/state/skillset/copies/`.
- Links are recorded in `~/.local/state/skillset/installed-links.json`. Recorded links that are no longer wanted are removed, only when they point into the catalog. Anything else in the way is a clash: nothing changes and every clash is listed. A real `~/.claude/CLAUDE.md` is a clash too: move it aside and rerun.
- `.docs/` is added once to git's global excludes file, so a `<repo>/.docs` link to a private project knowledge store is never committed.
- Codex settings are not installed: merge missing keys from `config/codex/config.toml` into `~/.codex/config.toml` by hand.

## Developing

The Makefile runs the installer from this checkout with `--catalog config`, so links point into the working tree and edits are live.

```sh
make init        # pick hosts and sets from this checkout's catalog
make plan        # preview
make apply       # install from the working tree
make uninstall   # remove it (make uninstall/plan to preview)
make validate    # lint, build, test, doc-rot guard
make help        # everything else
```

Optional extras: `make install/dcg` installs the dcg binary and `make install/plugins` the Claude plugins, marketplaces and MCP servers listed in `config/claude/settings.json`.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow.
