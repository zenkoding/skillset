---
name: triage-prs
description: >
  Use when triaging open pull requests across recurring review sessions —
  deciding which PRs need attention, which to skip, which to land, or
  which already have a reviewer. Triggers on "triage PRs",
  "what should I review", "check PR queue", "review queue",
  "which PRs need attention".
argument-hint: "[owner/repo] | skip <num> \"<reason>\" | forget <num> | --include-mine | --include-others-drafts | --include-others-reviewed | --all"
user-invocable: true
---

> **Requires execution mode.** If plan mode is active, exit first.

## CRITICAL: Run Commands, Don't Describe Them

- Run every `gh` / `git rev-parse` step as a real shell command. Do not paraphrase command output.
- Read and write `pr-triage.md` and `github_handle.md` directly. Do not `mkdir` — the memory directory already exists.

The dashboard render path never asks the user anything. The dashboard is printed; the user reads it and runs the printed `Next:` command themselves.

This skill **MUST NOT** run `/review-pr` or `/land-pr` on the dashboard render path. The `Next:` line emits the command as text; the user dispatches the follow-up.

---

# Triage PRs — Printed Dashboard

You render a printed, sorted dashboard of open pull requests for one repository, classified by a deterministic six-state ladder and filtered by default rules that hide non-actionable noise. Sub-command verbs `skip` and `forget` write to per-PR memory and exit without rendering. Every invocation:

1. Parse args (sub-command verbs short-circuit before fetch).
2. Resolve the target repo and the user's GitHub handle (with auto-memory cache).
3. Fetch open PRs in a single batched `gh pr list`, then per-PR detail only where the state ladder needs it.
4. Load prior skip decisions from `pr-triage.md`.
5. Classify each PR against the state ladder (first-match-wins).
6. Apply default filters (with per-flag opt-out via override flags).
7. Sort by the priority ladder, then by `updatedAt` desc.
8. Render the per-PR blocks plus the summary footer.

**Core principle:** the render path is read-only. Writes happen only via the `skip` and `forget` sub-commands.

---

## Phase 0: Parse Args

The argument string `$ARGUMENTS` resolves to one of these forms. Detect the form **before** any GitHub call.

| Form | Behavior |
|---|---|
| (empty) | Dashboard for the current repo with default filters. |
| `owner/repo` | Dashboard for the given repo with default filters. |
| `--include-mine`, `--include-others-drafts`, `--include-others-reviewed`, `--all` (alone or combined with `owner/repo`) | Dashboard with the named default filter disabled. `--all` disables all three. Flags compose. |
| `skip <num> "<reason>"` | Append-or-update a skip block in `pr-triage.md` for PR `<num>`, then exit. No fetch, no render. |
| `forget <num>` | Remove the matching block from `pr-triage.md`, then exit. No fetch, no render. If no block matches, emit a single-line notice and exit successfully. |

Parsing rules:

1. If the first whitespace-separated token is the literal `skip` → treat the next token as `<num>` (integer) and everything between the next pair of double quotes as `<reason>`. If no double quotes are present, treat the entire remainder after `<num>` as the reason string. Route to **Sub-command: skip** below; do not run Phases 1–6.
2. Else if the first token is the literal `forget` → treat the next token as `<num>` (integer). Route to **Sub-command: forget** below; do not run Phases 1–6.
3. Else split the remaining tokens: a single token containing a `/` is `owner/repo`; tokens starting with `--` are override flags. Anything else is invalid → emit a one-line usage notice and exit.

Store override-flag state as three booleans (`--all` is folded into each rather than a fourth variable):

```
include_mine             = "--include-mine"            in flags OR "--all" in flags
include_others_drafts    = "--include-others-drafts"   in flags OR "--all" in flags
include_others_reviewed  = "--include-others-reviewed" in flags OR "--all" in flags
```

---

## Phase 1: Resolve Repo and Handle

### 1.1 Compute the git root and encoded memory namespace

```bash
GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$GIT_ROOT" ]; then
  echo "Not in a git repo — /triage-prs needs a git repo to anchor its auto-memory. Either cd into a repo or pass owner/repo (and run from a repo for memory to persist)." >&2
  exit 1
fi
ENCODED=$(echo "$GIT_ROOT" | sed 's|/|-|g')
MEMORY_DIR="$HOME/.claude/projects/$ENCODED/memory"
MEMORY_FILE="$MEMORY_DIR/pr-triage.md"
HANDLE_FILE="$MEMORY_DIR/github_handle.md"
MEMORY_INDEX="$MEMORY_DIR/MEMORY.md"
```

**Memory paths are anchored at `git rev-parse --show-toplevel`, not CWD.** Sub-directory invocations resolve to the same namespace as repo-root invocations. Worktrees resolve to their main repo via `--show-toplevel`. Do not derive these paths from `$PWD`.

### 1.2 Resolve the target repository

```bash
if [ -n "$EXPLICIT_OWNER_REPO" ]; then
  TARGET_REPO="$EXPLICIT_OWNER_REPO"
else
  TARGET_REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null)
  if [ -z "$TARGET_REPO" ]; then
    echo "Could not resolve target repo. Run inside a GitHub repo clone or pass owner/repo." >&2
    exit 1
  fi
fi
```

### 1.3 Resolve the user's GitHub handle

Handle resolution order (first match wins):

1. **`override:` line in `github_handle.md`** — if `github_handle.md` exists and contains a line `override: <handle>` outside the YAML frontmatter, use that value.
2. **`login:` line in `github_handle.md`** — if no override is present, use the cached `login` field from the file body.
3. **`gh api user --jq .login`** — if the file is absent OR `gh` later returns an authentication error during PR fetch, query the live API. Write the result to `github_handle.md` and add an index entry to `MEMORY.md`.

Read with the `Read` tool. If the file is absent, query the API and write it:

```bash
HANDLE=$(gh api user --jq .login 2>/dev/null)
if [ -z "$HANDLE" ]; then
  echo "gh is not authenticated — run 'gh auth login' before /triage-prs." >&2
  exit 1
fi
```

Then **Write** `github_handle.md` using the schema in [memory-format.md](references/memory-format.md) (`github_handle.md` section), and **Edit** (or **Write** if absent) `MEMORY.md` to add the index entry defined there.

`MEMORY.md` has no frontmatter per the Claude Code auto-memory spec — append the index line at the bottom; do not reorder existing content.

The resolved value is stored as `ME`. **`ME` is used in every "author == me" / "author != me" check in Phases 4 and 5.**

---

## Phase 2: Fetch Open PRs

### 2.1 Batched metadata fetch

```bash
gh pr list \
  --repo "$TARGET_REPO" \
  --state open \
  --limit 200 \
  --json number,title,author,isDraft,mergeable,reviewDecision,headRefOid,headRefName,baseRefName,createdAt,updatedAt,statusCheckRollup,reviewRequests,url
```

This single call carries every field needed by the state ladder and the default filters. Do not issue a per-PR query yet.

### 2.2 Per-PR detail fetch (only when needed)

For each non-draft PR that the state ladder may classify as `needs_re_review`, `has_pending_review`, or `fresh` — i.e., any PR that is not eliminated by the `draft` rule and not yet classified as `ready_to_merge` / `needs_landing` purely from rollup + `reviewDecision` — fetch reviews, issue comments, review comments, and commits:

```bash
gh pr view <num> --repo "$TARGET_REPO" --json reviews,comments,commits
```

Drafts and PRs already classifiable purely from `reviewDecision` + `statusCheckRollup` (i.e., `ready_to_merge` or `needs_landing`) are skipped — this bounds the per-PR call count to `1 + N_potential_review_candidates`.

### 2.3 Sandbox / TLS fallback

If `gh` exits with a TLS cert-verify or connection error (`x509: certificate signed by unknown authority`, `connection refused`, `Could not resolve host`), retry the same call with `dangerouslyDisableSandbox: true`. This applies to any `gh` call in Phases 1–6 and the sub-commands — not only the batched `pr list`. This is a known sandbox symptom, not a `gh` problem.

If `gh` exits with `HTTP 401` or any auth error AND `github_handle.md` already existed → invalidate the cached handle and rerun Phase 1.3 with `gh api user --jq .login`. Re-fetch.

---

## Phase 3: Load Skip Memory

Read `$MEMORY_FILE` with the `Read` tool if it exists. If missing, treat as empty (no prior decisions).

Parse into a map: `{ "owner/repo#num": { decision, reason, sha, date } }`.

### Memory file format

See [memory-format.md](references/memory-format.md) — `pr-triage.md` section. Covers the file frontmatter, per-PR block schema, parsing rule (each `## owner/repo#N` header begins a block; unknown fields preserved for forward compatibility), and SHA-match semantics.

---

## Phase 4: Classify Each Open PR

For each fetched PR, compute `ci_status`, evaluate the state ladder, and apply SHA-matched skip filtering.

### 4.1 Derive `ci_status` (from `statusCheckRollup` + bot reviews folded in)

`statusCheckRollup` is an array of check entries, each with `state` / `conclusion`. Evaluate **after** folding bot reviews into the rollup:

- **Bot detection.** A review is considered a bot review when the author login matches the allowlist regex `(.*-bot$|dependabot\[bot\]|renovate\[bot\])`. (`reviews[].author` from `gh pr view --json reviews` does not expose an `is_bot` field — regex is the sole detection mechanism.) Bot reviews **do not** count toward `human_review_exists`. They count toward `ci_status` instead:
  - bot review state `APPROVED` → equivalent to a passing check
  - bot review state `CHANGES_REQUESTED` → equivalent to a failing check
  - bot review state `COMMENTED` or `PENDING` → ignored
- **Rollup interpretation (after bot folding):**

  ```
  GREEN   → every entry is SUCCESS or NEUTRAL (and no bot CHANGES_REQUESTED)
  RED     → any entry is FAILURE, ERROR, or TIMED_OUT (or any bot CHANGES_REQUESTED)
  PENDING → any entry is PENDING, IN_PROGRESS, or QUEUED, and none are RED
  NONE    → statusCheckRollup is empty and no bot reviews exist
  ```

  Precedence when multiple apply: RED > PENDING > GREEN > NONE.

### 4.2 Compute `human_review_exists` and `author_activity_after_last_human_review`

- **`human_review_exists`** = there is at least one review whose author is **not** a bot (per the bot detection rule above).
- **`author_activity_after_last_human_review`** = let `T` = `submitted_at` of the most recent non-bot review. Returns true when, after `T`:
  - **(a)** a commit appears on the PR head with `committer.date > T` AND that commit has exactly one parent (i.e., not a merge commit with the base), **OR**
  - **(b)** a PR-issue comment, review comment, or inline review comment authored by `PR.author` has `createdAt > T`.

  If no non-bot review exists, this evaluates to false.

### 4.3 Evaluate the state ladder (first match wins — STRICT ORDER)

Evaluate the seven rules below in this exact order. The first rule that matches assigns the state; later rules are not considered.

```
1. draft               → PR.isDraft == true
2. needs_landing       → mergeable == "CONFLICTING"
                         OR (reviewDecision == "APPROVED" AND ci_status == "RED")
3. ready_to_merge      → reviewDecision == "APPROVED"
                         AND mergeable != "CONFLICTING"
                         AND ci_status in ("GREEN", "PENDING")
4. needs_re_review     → human_review_exists
                         AND author_activity_after_last_human_review
5. has_pending_review  → human_review_exists
6. fresh               → human_review_exists == false AND ci_status == "GREEN"
7. other               → fallthrough — anything that did not match above
```

**Approval-with-new-commits rule (FR-005 + key decision):** when `reviewDecision == "APPROVED"` and a new single-parent commit landed after the approval, the PR stays in `ready_to_merge`. Approval implies trust; the skill **MUST NOT** downgrade approved-with-new-commits to `needs_re_review`. Rule 4 only matches when the PR was not approved (i.e., the prior rule fell through).

**Approval-with-pending-CI rule:** rule 3 explicitly accepts `ci_status == "PENDING"`. Pending CI is not a hard block once a human approval exists.

### 4.4 Apply SHA-matched skip filter

For each PR:

- If `pr-triage.md` has a block for `<owner/repo>#<num>` AND the block's `Triaged at SHA` equals the PR's live `headRefOid` AND the block's `Decision` is `skip` or `defer` → **exclude this PR from the dashboard entirely**. Skip-memory-filtered PRs do **not** count toward the footer's "hidden by default filters" tally.
- If a block exists AND the SHAs differ → **include the PR** and record `resurfaced = true` (used in Phase 6 to print a banner).
- If no block exists → include the PR normally.

---

## Phase 5: Apply Default Filters

Apply these three filters in order. Each is individually disable-able via an override flag.

```
F1 (my non-actionable PRs):
    Hide if (PR.author == ME)
            AND state ∉ {needs_landing, ready_to_merge, needs_re_review}
            AND NOT (state == draft AND ci_status == GREEN)
    Disabled by: --include-mine
    (Notably: my drafts with GREEN CI survive — they appear under "my drafts"
     in the priority ladder so I can mark them ready.)

F2 (drafts from other authors):
    Hide if PR.author != ME AND state == draft
    Disabled by: --include-others-drafts

F3 (other-authored PRs already being reviewed by someone other than me):
    Hide if PR.author != ME
            AND PR has at least one non-bot review whose author != ME
            AND ME not in reviewRequests[].login
            AND none of my teams appear in reviewRequests[].slug
    Disabled by: --include-others-reviewed
    (Bot reviews do NOT trigger F3 — bots fold into ci_status, not "reviewed by others".
     A bot-only-reviewed PR from another author correctly stays visible as `fresh`.)
```

`--all` disables F1, F2, and F3 simultaneously.

Count each PR removed by F1, F2, or F3 toward `HIDDEN_BY_DEFAULT`. **Do not** count PRs already removed by SHA-matched skip-memory (Phase 4.4) — those were an explicit user decision, not a default-filter action.

`reviewRequests` (array of `{__typename, login?, slug?}` entries) is read from the Phase 2 batched fetch. User entries carry `login`; team entries carry `slug` (format `org/team-name`).

---

## Phase 6: Sort and Render

### 6.1 Sort

Sort all surviving PRs by this priority ladder, then by `updatedAt` descending within each bucket:

```
1. ready_to_merge
2. needs_landing
3. needs_re_review
4. fresh
5. my drafts (only those with ci_status == GREEN that survived F1)
6. other
```

`has_pending_review` PRs (state 5 in the ladder) fall into bucket 6 (`other`) for display purposes — they are not actionable for the user (someone else is reviewing) and do not need their own bucket.

### 6.2–6.4 Render — see [render-format.md](references/render-format.md)

Emit the dashboard as GitHub-flavored Markdown. For the full format specification — section name mapping, per-bucket table layouts, field formats (AGE, CI glyph, Next command mapping), re-surfaced PR annotation, Skipped / Deferred table, and summary footer — consult [render-format.md](references/render-format.md).

Key constraints that must hold regardless:
- Omit empty buckets entirely; do not emit a section header with zero rows.
- `Next` cell must be a verbatim copy-paste-safe command (code span), never a paraphrase.
- Skipped / Deferred section is always rendered when skip-memory PRs are still open.
- Footer counts exclude Skipped / Deferred (those are explicit user decisions, not default-filter noise).

---

## Sub-command: `skip <num> "<reason>"`

Route here from Phase 0 when the first argument token is `skip`.

1. Resolve `TARGET_REPO`, `GIT_ROOT`, `MEMORY_FILE` per Phase 1.1, 1.2.
2. Fetch the live head SHA for PR `<num>`:
   ```bash
   HEAD_SHA=$(gh pr view <num> --repo "$TARGET_REPO" --json headRefOid -q .headRefOid)
   ```
3. Fetch title and author for the block body:
   ```bash
   gh pr view <num> --repo "$TARGET_REPO" --json title,author -q '"\(.title)\t\(.author.login)"'
   ```
4. Edit `pr-triage.md`: if a `## <owner>/<repo>#<num>` block exists, replace it; otherwise append a new block at the end. If `pr-triage.md` does not exist, **Write** it using the file frontmatter and block schema from [memory-format.md](references/memory-format.md) (`pr-triage.md` section).

5. Update `MEMORY.md` index if absent — append the `[PR Triage Log](pr-triage.md)` index entry from [memory-format.md](references/memory-format.md); do not reorder existing entries.
6. Print a single confirmation line: `Recorded skip for <owner>/<repo>#<num> @ <head_sha[:7]>.`
7. Exit. **Do not** render the dashboard.

---

## Sub-command: `forget <num>`

Route here from Phase 0 when the first argument token is `forget`.

1. Resolve `TARGET_REPO`, `MEMORY_FILE` per Phase 1.1, 1.2.
2. Read `pr-triage.md`. If absent → print `No memory file at <path>; nothing to forget.` and exit successfully.
3. Locate the `## <owner>/<repo>#<num>` block. If absent → print `No memory block for <owner>/<repo>#<num>; nothing to forget.` and exit successfully.
4. **Edit** `pr-triage.md`: delete that block (header through the last bullet, plus the trailing blank line). **Do not touch any other block.**
5. Print: `Removed memory block for <owner>/<repo>#<num>.`
6. Exit. **Do not** render the dashboard.

**Constraint:** `forget` deletes exactly one block — the one named in the invocation. The skill **MUST NOT** delete or rewrite memory blocks for other PRs.

---

## Negative Constraints (binding — see spec.md "Negative Constraints")

The skill **MUST NOT**:

1. Ask the user anything on the dashboard render path — the dashboard is printed and the user acts on it.
2. Execute any GitHub **write** operation (merge, comment, request-review, label) during dashboard rendering — the `Next:` line emits the command; the user runs it.
3. Mutate `pr-triage.md` during dashboard rendering — writes happen only on `skip` and `forget` sub-commands.
4. Count bot reviews toward `human_review_exists` (bot reviews fold into `ci_status`).
5. Downgrade an approved PR with new commits from `ready_to_merge` to `needs_re_review` — approval implies trust per the chosen model.
6. Count PRs filtered by SHA-matched skip-memory toward the footer's "hidden by default filters" tally — those PRs are an explicit user decision, not a default filter.
7. Resolve memory paths from `$PWD` — `$GIT_ROOT = git rev-parse --show-toplevel`, then encode.
8. Delete or rewrite memory blocks for any PR other than the one named in a `forget <num>` invocation.

---

## Migration from the menu-loop version

Users coming from the prior interactive version of this skill:

1. **First invocation auto-creates `github_handle.md`** — no user action required. `gh api user --jq .login` runs once, writes the cache, and adds the `MEMORY.md` index entry.
2. **Existing `pr-triage.md` blocks continue to apply.** The block schema is unchanged. PRs skipped at a remembered SHA stay hidden on first run of the new version. On the first run where a remembered PR's head SHA differs, the PR re-surfaces with a banner (instead of silently appearing).
3. **Skip is now an explicit sub-command.** Where you previously chose `Skip` from a menu, you now run `/triage-prs skip <num> "<reason>"` as a separate invocation.
4. **`/review-pr` and `/land-pr` are still the canonical follow-ups.** The dashboard prints the command verbatim under `Next:`; you copy and run it.

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Asking the user a question during render | Forbidden. The dashboard is printed; the user reads it. |
| Invoking `/review-pr` or `/land-pr` automatically | Forbidden on the render path. The `Next:` cell emits the command; user dispatches. |
| Counting bot reviews toward `human_review_exists` | Forbidden. Bots fold into `ci_status` (APPROVED = pass; CHANGES_REQUESTED = fail). |
| Treating bot reviews as "reviewed by others" in F3 | Forbidden. F3 checks for **non-bot** reviews only. A bot-only-reviewed PR must stay visible as `fresh`. |
| Downgrading approved + new-commits to `needs_re_review` | Forbidden. Approval implies trust. The state stays `ready_to_merge`. |
| Counting skip-memory-filtered PRs in the footer `hidden by default filters` | Forbidden. Footer counts only PRs removed by F1/F2/F3. |
| Resolving memory paths from `$PWD` | Use `git rev-parse --show-toplevel`, then `sed 's|/|-|g'`. Worktrees resolve to the main repo. |
| Mutating `pr-triage.md` during render | Render is read-only. Writes only via `skip` and `forget`. |
| `mkdir -p` on memory dir | It already exists. Write directly. |
| Querying per-PR detail for drafts | Drafts are classified before per-PR fetch. Skip the detail call for them. |
| Forgetting to invalidate cached handle on `gh` 401 | Re-run `gh api user --jq .login` on auth error; refresh the cache. |
| Using internal state names as section headings | Use the human-readable mapping: `needs_landing` → "Needs Landing", `fresh` → "Fresh — Needs First Review", etc. |
| Omitting the Skipped / Deferred section | Always render it when skip-memory PRs are still open — it lets the user see what's parked and why. |
| Using `requestedReviewers` or `requestedTeams` in `gh pr list --json` | These fields do not exist. Use `reviewRequests` (array of `{__typename, login?, slug?}`). |

---

## Quick Reference

```bash
# Memory paths (anchored at git root, never CWD)
GIT_ROOT=$(git rev-parse --show-toplevel) && \
  ENCODED=$(echo "$GIT_ROOT" | sed 's|/|-|g') && \
  MEM="$HOME/.claude/projects/$ENCODED/memory/pr-triage.md" && \
  HND="$HOME/.claude/projects/$ENCODED/memory/github_handle.md"

# Resolve target repo
TARGET_REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner)

# Resolve handle (first time)
gh api user --jq .login   # → write to $HND with type: preference frontmatter

# Batched PR metadata fetch
gh pr list --repo "$TARGET_REPO" --state open --limit 200 --json \
  number,title,author,isDraft,mergeable,reviewDecision,headRefOid,\
headRefName,baseRefName,createdAt,updatedAt,statusCheckRollup,\
reviewRequests,url

# Per-PR detail (only when state ladder needs it)
gh pr view <num> --repo "$TARGET_REPO" --json reviews,comments,commits

# Next commands by bucket → see references/render-format.md (Next command mapping)
```
