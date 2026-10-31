---
name: review-pr
description: Review a GitHub pull request for product fit, rationale, approach, security and bugs, show the summary, and post it to the PR as one GitHub review once confirmed. Use when asked to review a PR, or from a PR review routine; not for branches or uncommitted changes without a PR.
compatibility: Requires Git and a readable checkout; authenticated gh for pull requests. Parallel subagents are optional.
---

# Review PR

Produce a review whose every finding is real, located, and worth the author's time. Be pragmatic: perfect is the enemy of good. A short review with three confirmed defects beats a long one padded with style notes, and "no findings" is a valid result.

Arguments: the PR number or URL, then `post` for an unattended run (a routine) that posts without asking. Without `post`, show the review and post it only after the user confirms.

## 1. Scope the change

Resolve the PR (`gh pr view`, `gh pr diff`); a URL resolves to its number first. Diff against the PR's own base: in a stack that is the PR below it, so review only this slice. Record the head SHA and check it out (`gh pr checkout <pr>`, in a worktree when the checkout has other work) so you can read whole files. Without a PR, stop and say so: this skill reviews PRs only.

Stop without reviewing when the PR is a draft, closed or merged, or the diff is empty. Also stop when a review carrying the `<!-- review-pr -->` marker already exists on the head SHA. A review on an older head means this is a re-review: review only the commits since that head, and never repeat a finding already posted.

Size the change with [work sizing](references/work-sizing.md): its counted lines, whether it is large, which risky areas it touches, and its door (one-way or two-way).

Done when you can state in two sentences what the change is meant to do and which files carry the risk.

## 2. Gather context

Read, before judging anything:

- **The PR body.** Its `### Reviewer notes` (written by `publish-pr`): start from the Merge danger and, when present, the Verdict, Must check items and High files, then verify the rest of the diff rather than trusting the notes. Its How to test steps are what a reviewer would run.
- **The intent.** The linked issue, or the tracker's plan when the repo uses one: `docs/agents/issue-tracker.md` says where plans live and how to find one from a branch name, PR title or body. Its acceptance criteria (**Done when** or similar) are what the change is judged against.
- **The rules.** `REVIEW.md` at the repo root, else the [default rules](references/default-rules.md). Then `AGENTS.md` or `CLAUDE.md`, the docs catalog in `docs/agents/README.md`, `CONTEXT.md`, the vision or non-goals doc when the repo has one, and ADRs in the touched area.
- **The code.** Each touched file around the changed lines, not only the hunks; callers of changed interfaces; recent history. Run the project's build, lint and tests when that is cheap and safe; record what ran and the result.

Done when every changed file is matched to the rules that govern it and you hold the intent, or a stated reason there is none.

## 3. Review, in order of importance

1. **Product fit.** Does the change belong in this project? Check it against the vision, non-goals and `CONTEXT.md`.
2. **Rationale.** Does it add a useful feature, make a necessary enhancement, or fix a real issue? Give each acceptance criterion and each Reviewer notes Must check item a verdict: met, partial, missing, or not checkable from the diff. Note scope nobody asked for.
3. **Approach.** Does it reach its goal in a reasonable way? Weigh complexity (less is better), consistency with the codebase, and verbosity (less code is better while it stays easy to read). Check the repo's written rules here.
4. **Security and bugs.** Trace each changed path for a concrete input that breaks it: error paths, boundaries, leaks on failure, concurrency, callers of changed behaviour. Run a separate security pass when the diff touches a risky area. Note missing or weak tests for new behaviour.

An earlier area outranks a later one: a change that doesn't belong needs no bug hunt. Skip what `REVIEW.md` or the defaults say not to report.

For a large diff, when the host supports subagents, run the four areas as parallel finders, each given the scope summary, the intent, the rules and the touched files. Then hand all candidates to one verifier subagent for step 4. Smaller diffs run one pass yourself.

Done when every touched file has been read under every area that applies to it.

## 4. Validate

Reopen the code for every candidate finding: the line exists, the claim holds, the fix is plausible. Drop what you cannot point at with a file and line. Merge findings on the same symbol into one. Tag a defect the PR did not introduce as **pre-existing**, so the author isn't blamed for it.

Label each finding with the [semantic review](https://www.m31coding.com/blog/semantic-reviews.html) label that fits:

| Label | Means | Blocks merge |
| --- | --- | --- |
| `Crucial` | Severe: must change before merge | yes |
| `Important` | A change is expected; discuss if the author disagrees | no |
| `Suggestion` | Worth considering; the author decides | no |
| `Question` | An honest question that needs an answer; never rhetorical | no |
| `Hint` | Useful for the future; no change expected | no |
| `Remark` | An observation; no change or answer expected | no |

Use `Crucial` only for wrong behaviour, a security hole, data loss, a regression to shipped behaviour, or a change that doesn't belong. One inflated label teaches the author to ignore every label. Each `Crucial` finding says what would overturn it, as a question the author can answer.

## 5. Write the review

```markdown
<!-- review-pr -->
**<n> findings**: Crucial <n> · Important <n> · Suggestion <n> · Question <n> · Hint <n> · Remark <n>

**Intent**: <one or two sentences: the problem and how the change solves it>

**Verdict**: <merge | merge after the Important findings | needs changes>. <the findings that decide it>

**Product fit**
<findings, or "Fits.">

**Rationale**
- [x] <criterion or Must check item>: <evidence>
- [ ] <criterion>: <what is missing, or "not checkable from the diff">

**Approach**
<findings, or "No concerns.">

**Security and bugs**
<findings, or "No issues found.">

<what ran; what was not checked>
```

Each finding is `Label: <mechanism, then the ask>` with its `path:line`, in two to four sentences. Use a `suggestion` block for a concrete one-line change. State conclusions, not process. A clean change gets the counts line, the intent and a one-line verdict.

**Verdict** is `merge` only when nothing is `Crucial`, every checkable criterion is met, and a one-way door has a stated rollback or migration plan. It is a recommendation: never merge, approve or request changes yourself.

## 6. Deliver

Publishing follows the [posting guide](references/posting.md): one `COMMENT` review on the head SHA, then remove the `ready-for-agent` label when the PR has it.

**With a person.** Show the summary: the counts, the verdict and every `Crucial` and `Important` finding with its `path:line`. Ask whether to post it. Post on a yes, applying any edits the user asks for first; on a no, post nothing. Nothing is saved to the repo either way.

**`post`** (a routine runs this). Post without asking. Nobody can answer follow-ups, so the run also:

- drops findings it could not verify, instead of marking them unverified;
- caps the review at three `Crucial` findings and eight comments, dropping the lowest-ranked;
- says so plainly when its review could not complete, and posts nothing rather than a clean review that didn't happen.

Done when exactly one `<!-- review-pr -->` review exists on the head SHA, or the user declined to post.
