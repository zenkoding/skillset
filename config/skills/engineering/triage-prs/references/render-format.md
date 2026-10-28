# Dashboard Render Format

Reference for Phase 6.2–6.4. Use this file when rendering the PR dashboard.

---

## Section name mapping

Use these exact heading strings — never the internal state names:

| Internal state | Section heading |
|---|---|
| `ready_to_merge` | Ready to Merge |
| `needs_landing` | Needs Landing |
| `needs_re_review` | Needs Re-review |
| `fresh` | Fresh — Needs First Review |
| my drafts bucket (draft + author==ME + ci==GREEN) | My Drafts |
| `other` / `has_pending_review` | Other |

Emit `### Section Name (N)` as the bucket header, followed by a Markdown table. Omit empty buckets entirely.

---

## Table columns by bucket

*Ready to Merge, Needs Landing, Needs Re-review, Fresh — Needs First Review:*
```
| PR | Author | Age | CI | Title | Next |
|----|--------|-----|----|-------|------|
| [#NUM](url) | login | AGE | GLYPH | title | `command` |
```

*My Drafts* (author column omitted — always the current user):
```
| PR | Age | CI | Title | Next |
|----|-----|----|-------|------|
| [#NUM](url) | AGE | GLYPH | title | `gh pr ready NUM --repo owner/repo` |
```

*Other* (no Next column; add Notes for brief context):
```
| PR | Author | Age | CI | Title | Notes |
|----|--------|-----|----|-------|-------|
| [#NUM](url) | login | AGE | GLYPH | title | CI failing / CI not started yet / review in progress |
```

---

## Field formats

| Field | Format |
|---|---|
| `PR` cell | `[#NUM](url)` — clickable link |
| `AGE` | relative age of `updatedAt`: `today`, `3d`, `2w`, `1mo` |
| `CI` glyph | `✓` GREEN · `✗` RED · `⏳` PENDING · `–` NONE |
| `Title` | full PR title, not truncated |
| `Next` cell | verbatim command in a code span — copy-paste-safe |

---

## Next command mapping

| Bucket | Command in `Next` cell |
|---|---|
| Ready to Merge | `` `gh pr merge NUM --repo owner/repo --squash --auto` `` |
| Needs Landing | `` `/land-pr NUM` `` |
| Needs Re-review | `` `/review-pr https://github.com/owner/repo/pull/NUM` `` |
| Fresh — Needs First Review | `` `/review-pr https://github.com/owner/repo/pull/NUM` `` |
| My Drafts | `` `gh pr ready NUM --repo owner/repo` `` |

---

## Re-surfaced PRs

When a PR has `resurfaced = true` (Phase 4.4 SHA mismatch), append to the title cell:

```
*(re-surfaced — was: decision @ sha7)*
```

`sha7` is the first 7 chars of the old `Triaged at SHA` from memory, not the live SHA.

---

## Skipped / Deferred section

Always render this section after all active-PR buckets when skip-memory PRs are still open (appear in the Phase 2 fetch). Omit if no such PRs exist.

```
### Skipped / Deferred (N)

| PR | Author | Age | Title | Reason |
|----|--------|-----|-------|--------|
| [#NUM](url) | login | AGE | title | reason from memory block |
```

`Reason` is the `**Reason**` field from the matching `pr-triage.md` block. Sort by `updatedAt` descending. No `Next` column — these PRs are intentionally parked.

---

## Summary footer

End the dashboard with exactly this single line:

```
Shown: **N PRs** · X hidden by default filters (`--all` to see them)
```

- `N` = count of PRs in the active buckets (excludes Skipped / Deferred).
- `X` = `HIDDEN_BY_DEFAULT` from Phase 5. Does not include Skipped / Deferred (those are explicit user decisions, not default-filter noise).
