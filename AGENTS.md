# Agent Instructions: skillset

Read [CONTRIBUTING.md](CONTRIBUTING.md) first: it indexes the dev docs, and everything a human contributor needs applies to you too.

## Agent-only rules

- At session start, invoke `i-have-adhd` before the first substantive response. If this host has no skill invocation tool, read the installed `i-have-adhd/SKILL.md` and apply it. Keep it active throughout the session, including after compaction, unless the user says "stop adhd mode" or "normal mode".
- Run `make install` to build the binary from this checkout. Preview with `skillset --catalog "$PWD/config" plan` before installing anything into `$HOME`.
- For issue tracker operations, follow [docs/agents/issue-tracker.md](docs/agents/issue-tracker.md).
- Write commit messages and PR titles as Conventional Commits, in the format `CONTRIBUTING.md` gives. Never use `--no-verify` to get past the `commit-msg` hook.
- Markdown: no hard wraps. One line per paragraph or list item; the viewer wraps.

## Pull requests

Build every PR body from [.github/pull_request_template.md](.github/pull_request_template.md):

- Keep the template's headings, in their order, and replace each HTML comment with real content. Leave no comment or placeholder behind.
- Under **Why**, give the reason for the change and any decision it makes, with the rejected alternative. Add `Closes #<n>` when the PR closes an issue. Link only issues you have verified.
- Under **How to test**, separate the checks you ran and their results from what you skipped or could not run. Never claim a check you did not run.
- Under **Reviewer notes**, always fill in **Merge danger**.
