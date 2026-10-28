# claude/settings.json: why each rule exists

Precedence is deny > ask > allow. `defaultMode: auto` handles everything
not listed, so `allow` only lists commands that should never need a check.
Change a rule here and in `settings.json` together.

## allow

Read-only shell basics, never worth a prompt:

- `cat`, `head`, `tail`, `grep`, `find`, `ls`, `tree`, `wc`: read files and directories.
- `command -v`, `which`: check whether a tool is installed.
- `date`, `uuidgen`: harmless value generators.
- `mkdir`: creates directories only; nothing is overwritten.

Docker (local containers, nothing leaves the machine):

- `docker build`, `images`, `inspect`, `logs`, `ps`, `exec`: build and inspect local containers.
- `docker-compose up`, `down`, `exec`, `logs`, `ps`: run local dev stacks (`down -v` is denied below).

GitHub reads:

- `gh auth status`, `gh repo view`: check auth and repo context.
- `gh issue list/view`, `gh pr list/view/diff/checks`, `gh release list/view`: read issues, PRs, CI and releases.

Git (local work plus non-force pushes; destructive forms are denied below):

- `git status`, `log`, `diff`, `show`, `describe`, `rev-parse`, `remote -v`: read history and state.
- `git add`, `commit`, `checkout`, `branch`, `merge`, `rebase`, `stash`, `tag`: everyday local work, recoverable via reflog.
- `git fetch`, `pull`: sync from remotes.
- `git push origin`, `push -u origin`: publish a branch; force forms are denied.
- `git push --force-with-lease`: safe rewrite of your own branch after a rebase.
- `git worktree add/list/remove`: isolated worktrees for parallel agents.

Go toolchain:

- `go build`, `run`, `test`, `vet`, `fmt`, `list`, `mod`, `gofmt`, `goimports`: build, test and format Go code.

Make targets:

- `make build`, `check`, `clean`, `coverage`, `fmt`, `help`, `lint`, `test`, `vet`: common local targets.
- `make plan`: this repo's installer dry run.

Kubernetes and Tanka reads:

- `kubectl get`, `describe`, `logs`, `top`, `explain`, `api-resources`, `version`, `auth can-i`, `config current-context`, `config view`: inspect clusters without changing them.
- `tk diff`, `show`, `env current`, `env list`, `tool`: render and diff Tanka environments without applying.

Files:

- `Read(//tmp/**)`: read temp files outside the project (logs, scratch output).

## ask

Visible to other people, but routine and easy to undo, so confirm each one:

- `gh pr create`, `edit`, `comment`, `ready`, `close`, `merge`, `review`: PR actions others see.
- `gh issue create`, `comment`, `edit`, `close`: issue actions others see.
- `gh release create`: publishes a release; can be deleted, but people may see it first.
- `gh api`: every call, reads included. Ask beats allow, so no allow rule for explicit `-X GET` could win over this prefix; and a method-based ask would miss implicit POSTs via `-f`/`-F`. Asking for all of them is the only safe option.

## deny

Destructive or hard to undo; run by hand if really needed:

- `git push --force`, `-f`: overwrite remote history (`--force-with-lease` is allowed).
- `git push --delete`, `push origin --delete`: delete remote branches.
- `git reset --hard`, `git clean -f`: throw away uncommitted work.
- `git branch -d/-D/--delete`: delete local branches.
- `git filter-branch`: rewrites all history.
- `kubectl apply`, `create`, `delete`, `patch`, `replace`, `edit`, `scale`, `set`, `rollout restart`: change cluster state.
- `kubectl exec`, `cp`, `port-forward`: reach into running workloads.
- `kubectl drain`, `cordon`, `uncordon`: change node scheduling.
- `tk apply`, `prune`, `delete`: change cluster state through Tanka.
- `gh repo delete`, `gh repo archive`: remove or freeze a repository.
- `gh release delete`: removes a published release and its assets.
- `gh issue delete`, `gh issue transfer`: lose or move issue history.
- `gh secret`: reads or changes repository secrets.
- `gh api ... DELETE`: API deletes, in every spelling (`-X DELETE`, `-XDELETE`, `--method DELETE`, `--method=DELETE`, and their lowercase forms).
- `docker push`: publishes images to a registry.
- `docker rmi`, `image rm`, `system prune`, `volume prune`, `docker-compose down -v`: delete images or data volumes.
- `rm -rf`, `rm -r`, `rm -fr`: recursive deletes.
- `sudo`: no privilege escalation from an agent.
- `chmod 777`: world-writable files.

## Other keys

- `hooks.SessionStart` runs `~/.claude/hooks/project-knowledge.sh` (linked from `config/claude/hooks/`) on every start, resume, clear and compaction. It prints the repo's project knowledge catalogs, the private store's `~/.local/share/agents/projects/<key>/README.md` and `docs/agents/README.md`, plus the store's absolute path, into context; see "Project knowledge" in `config/instructions/shared.md`. Outside a repo, without catalogs, or before `make install` links the script, it prints nothing and exits 0.
- `hooks.PreToolUse` runs `dcg` (destructive command guard) from PATH before Bash, PowerShell and Monitor. Without `dcg` installed the hook exits 0 and does nothing. `exec` keeps stdin, so dcg still gets the hook JSON, and keeps dcg's exit code.
- `sandbox.filesystem.allowWrite`: `~/.local/share/agents` holds each repo's private project knowledge store, written from any worktree.
- `permissions.additionalDirectories`: `~/.local/share/agents`, so Edit and Write reach the private store without a prompt. Both grants must be user-level: project settings ignore them, and worktree sessions block writes into the main checkout.
- `statusLine` runs `~/.claude/statusline.sh`.
- `enabledPlugins`: `mattpocock-skills` is the current development skill set; `cc-plugin-you-should-know` is built in.
- `env`: disables the auto-updater, enables agent teams, and runs subagents on Sonnet.
- The rest (`modelSettings`, `outputStyle`, `theme`, `editorMode`, `tui` and so on) are personal preferences.
