# Memory File Formats

Reference for the two memory files /triage-prs reads and writes.
Canonical schema — do not deviate from these formats.

---

## `github_handle.md`

Written once (Phase 1.3) when the file is absent, or rewritten on `gh` 401. Format:

```markdown
---
name: github-handle
description: Cached GitHub handle for /triage-prs self-vs-other classification.
type: preference
---

login: <handle>
# override: <other-handle>   # uncomment to force a different handle (e.g., service account)
```

Resolution priority (first match wins):
1. `override:` line outside YAML frontmatter → use that value
2. `login:` line → use cached value
3. File absent or `gh` returns 401 → query `gh api user --jq .login`, rewrite file

`MEMORY.md` index entry (append once, do not duplicate):
```
- [GitHub handle](github_handle.md) — cached login for /triage-prs self-vs-other classification
```

---

## `pr-triage.md`

Written only by `skip` and `forget` sub-commands (render is read-only).

**File frontmatter:**
```markdown
---
name: PR triage decisions
description: Per-PR triage decisions from /triage-prs. Re-surfaces when head SHA changes.
type: project
---

# PR Triage Log

One block per PR. Decision honored while head SHA matches `Triaged at SHA`.
```

**Per-PR block schema** (one block per PR, separated by blank lines):
```markdown
## owner/repo#<num>
- **Title**: <title>
- **Author**: <author>
- **Decision**: skip
- **Reason**: <reason>
- **Triaged at SHA**: <head_sha>
- **Triaged on**: <YYYY-MM-DD>
- **Re-surface if**: head SHA changes
```

**Parsing rule:** Each `## owner/repo#N` header begins a new block. Parse the bulleted fields below until the next `## ` header or EOF. Unknown fields are preserved as-is on any rewrite (forward-tolerant).

**SHA-match semantics:**
- Block SHA == live `headRefOid` → exclude PR from dashboard (skip active)
- Block SHA != live `headRefOid` → include PR with `resurfaced = true`
- No block → include PR normally

`MEMORY.md` index entry (append once, do not duplicate):
```
- [PR Triage Log](pr-triage.md) — per-PR skip/defer decisions for /triage-prs
```
