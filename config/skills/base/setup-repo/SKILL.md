---
name: setup-repo
description: Write or migrate a repo's human-first contributing docs (CONTRIBUTING.md, docs/contributing/) and the agent-only overlay (AGENTS.md, docs/agents/issue-tracker.md) that the other skills read their repo facts from, with built-in tracker and docs-store templates. Run as /setup-repo, or /setup-repo e2e or /setup-repo review for one part.
---

# Set up a repo

Write the docs a human contributor reads and the skills read too: `CONTRIBUTING.md` and `docs/contributing/`. `AGENTS.md` and `docs/agents/issue-tracker.md` only add what agents alone need; nothing agent-specific goes in `docs/contributing/`. Every write is shown to the user and confirmed first.

| Run | Writes |
| --- | --- |
| `/setup-repo` | `CONTRIBUTING.md`, `docs/contributing/{local-dev,testing,issues,docs}.md`, `AGENTS.md`, `docs/agents/issue-tracker.md` |
| `/setup-repo e2e` | the **E2E stack** section of `docs/contributing/testing.md` |
| `/setup-repo review` | `REVIEW.md` |

## 1. Look at what's there

Read, without asking:

- `README.md`, `CONTRIBUTING.md`, `docs/contributing/`, `AGENTS.md`, `CLAUDE.md`, `REVIEW.md`, `VISION.md`, `ARCHITECTURE.md`, `CONTEXT.md`;
- an older agent-only setup: `docs/agents/*.md` other than a current `issue-tracker.md` (tracker, labels, docs store, domain and e2e adapters, and its `README.md` catalog), `docs/vision.md`, `docs/adr/`;
- the build files (`Makefile`, `package.json`, `go.mod`, `pyproject.toml`, CI workflows) for the run and check commands;
- the tracker in use: the git remote host, a `.beads/` directory, tracker identifiers in branch names and commit messages;
- `.github/ISSUE_TEMPLATE/` and the repo's labels (`gh label list` on GitHub).

Done when you can say, per target file, whether it exists, what an older setup already says for it, and what you inferred.

## 2. Pick the backends

Ask, one question at a time with your recommendation from step 1 first (grill's rules, or `grill` when there's more to settle):

- **Tracker**: GitHub Issues, Linear or Beads; anything else starts from the blank template.
- **Docs store**: repo Markdown, Notion or Google Docs; anything else starts from the blank template.
- Facts the templates leave as `<...>` that step 1 couldn't settle: the Linear team and project, the Notion database, the Drive folder, who triages.

## 3. Draft the docs

Fill each template with what step 1 found and step 2 settled. The repo owns the copy from then on, so write it as this repo's doc, not as a template: no `<...>` left, no sections that don't apply.

| File | Template |
| --- | --- |
| `CONTRIBUTING.md` | [contributing](templates/contributing.md) |
| `docs/contributing/local-dev.md` | [local-dev](templates/local-dev.md); run commands inferred from the build files, each confirmed with the user |
| `docs/contributing/testing.md` | [testing](templates/testing.md); the check command from CI or the `Makefile` |
| `docs/contributing/issues.md` | [GitHub](templates/issues/github.md), [Linear](templates/issues/linear.md), [Beads](templates/issues/beads.md) or [blank](templates/issues/blank.md) |
| `docs/agents/issue-tracker.md` | the same tracker: [GitHub](templates/agents/issue-tracker/github.md), [Linear](templates/agents/issue-tracker/linear.md), [Beads](templates/agents/issue-tracker/beads.md) or [blank](templates/agents/issue-tracker/blank.md) |
| `docs/contributing/docs.md` | [repo](templates/docs/repo.md), [Notion](templates/docs/notion.md), [Google Docs](templates/docs/google-docs.md) or [blank](templates/docs/blank.md) |
| `AGENTS.md` | [agents](templates/agents.md); keep `CLAUDE.md` a symlink to it, or a one-line pointer where symlinks don't work |

An existing file is edited, not replaced: add what's missing, keep what the repo already says.

## 4. Migrate an older setup

When step 1 found agent-only adapters, or agent-only content in `docs/contributing/`, move their facts into the new docs:

| From | To |
| --- | --- |
| `docs/agents/issue-tracker.md` | stays, rewritten from its template: commands, bare-number resolution and the plan-edit rule; where work lives and status go to `docs/contributing/issues.md` |
| `docs/agents/triage-labels.md` | `docs/contributing/issues.md` |
| agent-only content in `docs/contributing/issues.md`: an **Operations** section, tracker commands, bare-number resolution, the rule on which plan sections an agent may edit | `docs/agents/issue-tracker.md` |
| `docs/agents/product-docs.md`, `doc-storage.md` | `docs/contributing/docs.md` |
| `docs/agents/e2e.md` | the **E2E stack** section of `docs/contributing/testing.md` |
| `docs/agents/domain.md` | nothing: skills read `CONTEXT.md` directly |
| `docs/agents/README.md` | a line in `CONTRIBUTING.md` for each doc a human needs; the rest stays in `docs/agents/` |
| `docs/vision.md` | `VISION.md`, with `git mv` so the history follows |
| `docs/adr/` | stays as frozen history; `docs.md` says so |

- **Keep the hand edits.** Compare each adapter with the stock text it started from (a template, or another tool's default); content the repo added is what matters, and it moves over in the repo's words.
- **Remove what you replaced.** Delete an adapter once every fact in it has a new home (never `docs/agents/issue-tracker.md`, which is current), and update `AGENTS.md` and `CLAUDE.md` lines that pointed at it.
- **Show the diff first**: every move, edit and deletion, before anything is written.

## 5. Check the merge settings

Git is the decision log: a squash merge must use the PR title and body as the commit message. On GitHub:

```sh
gh api repos/<owner>/<repo> --jq '{allow_squash_merge, squash_merge_commit_title, squash_merge_commit_message}'
```

Warn when squash merges are off, or the title isn't `PR_TITLE` or the message isn't `PR_BODY`, and say where to change it (Settings → General → Pull Requests). Don't change it yourself. On another host, say you couldn't check it.

## 6. Write

Show the user each file (new text, or a diff for an edit) and write only what they confirm. Write on a branch and open one PR with `publish` when installed, its **Why** saying which older setup it replaced.

Done when every confirmed file is written, the replaced adapters are gone, nothing in the repo still points at them, and the merge settings were checked or reported.

## Parts on demand

- **`/setup-repo e2e`**: add an **E2E stack** section to `docs/contributing/testing.md` in the [E2E stack format](references/e2e-stack.md), one subsection per way to run the app (host, containers). Infer the commands from the build files and confirm each; a step you can't settle stays a question, never a guess. `test` when installed reads this section.
- **`/setup-repo review`**: write `REVIEW.md` from the [default review rules](references/review-rules.md), then ask which defaults this repo changes (what blocks a merge, what not to report, label strength). Replace the defaults' opening paragraph with one about this repo. `review` when installed reads it in place of the defaults.

Each part follows steps 1, 5 and 6 for its one file.
