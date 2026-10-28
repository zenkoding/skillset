# PR authoring guidance

Follow the repository's PR template and conventions. If none exist, use a concise title describing the change and a body explaining the problem, resulting behavior, and validation. A before/after example helps when behavior is subtle; avoid restating the diff.

- Base the description on the actual changes and tests, not the branch name or diff statistics.
- Preserve the user's supplied title and body. If they conflict with a required repository rule, explain the specific conflict rather than silently rewriting them.
- Distinguish passed checks, failures, skipped checks, and work still pending. Do not invent test evidence or claim a PR is merged or approved without confirmation.
- Include related issue links only when they are relevant and verified. Use no placeholder references.
- Add attribution only when required by the applicable host or project instructions; use the actual authoring host and avoid duplicate footers.

When no other template exists for the PR body, use the following template instead:

```markdown
### What

<What this PR does. If possible, include examples which show the behaviour before & after this change.>

### Why

<Why is this PR necessary.>

### How it was tested

<How did you verify that the PR achieves its stated goal.>
```

An optional `### Review guide` section can follow the template's last section; see [review-guide.md](review-guide.md).

Before publication, check that the title and body describe the final scope, preserve supplied text, and contain no stale claims from earlier attempts. Use a separate temporary file for multiline text and pass it as data through `gh --body-file`; caller input files remain unchanged.
