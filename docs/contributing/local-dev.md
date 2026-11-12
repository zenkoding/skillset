# Local development

skillset is a command-line installer with no server to run. Local development means running a binary built from this checkout against this checkout's `config/`, so that the links it creates point into your working tree.

## Prerequisites

- Go at the version set in `go.mod`.
- `jq`, which the JSON config check uses.
- `curl`, which `make lint` uses to download a pinned golangci-lint into `target/`.

## Set up

```sh
git clone https://github.com/zenkoding/skillset.git && cd skillset
make git-setup                          # use the repo's commit template and commit-msg hook
make install                            # install the binary from this checkout
skillset init "$PWD"                    # pick hosts and sets from this checkout
```

`init` is optional. Without a saved config, every set is installed for both hosts.

## Run

```sh
skillset --catalog "$PWD/config" plan                 # preview what apply would change
skillset --catalog "$PWD/config" apply                # install from the working tree
skillset --catalog "$PWD/config" uninstall --dry-run  # preview what uninstall would remove
skillset --catalog "$PWD/config" uninstall            # remove what apply installed
```

Without `--catalog`, the binary uses the pinned release instead of your checkout. The `catalog:` line in its output shows which one it used.

`apply` changes your real home directory: `~/.claude`, `~/.codex`, `~/.agents` and `~/.config/dcg`. Run `plan` first.

Run `apply` again after you add a skill or change instructions or settings. Instructions are rendered at install time and `settings.json` is copied, so those edits don't appear until you do.

Run `make install` again after you change the Go code. `make help` lists every Makefile target. The README's Quick start covers the optional plugin and dcg installs.
