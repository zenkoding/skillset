---
name: review-pr
description: Review a pull request, branch, commit range, or uncommitted changes for correctness, security, design, and test gaps, and save a validated review. Use when asked to review changes in depth; posting to GitHub needs a separate request.
compatibility: Requires Git and a readable checkout; authenticated gh for pull requests. Parallel subagents are optional.
---

# Review PR

Produce a review whose every finding is real, located, and worth the author's time. A short review with three confirmed defects beats a long one padded with style notes.

## 1. Scope the change

Resolve the target: a PR number or URL (`gh pr view`, `gh pr diff`), a branch against its merge-base, a commit range, or the working tree. Read the PR description and any linked GitHub issue: that is the **intent** the change is judged against. When the body has a `### Review guide`, start from its Must check items and High files, then verify the rest of the diff rather than trusting the guide. Note the size, the touched areas, and whether security-sensitive code (auth, input parsing, secrets, crypto, permissions) is in the diff.

Done when you can state in two sentences what the change is meant to do and which files carry the risk.

## 2. Gather context

Read each touched file around the changed lines, not only the hunks. Check recent history and callers of changed interfaces. Run the project's build, lint, and tests on the change when that is cheap and safe; record what ran and the result.

## 3. Review along lenses

Apply each lens to the whole change:

- **Correctness**: logic errors, edge cases, error handling, concurrency, behaviour that diverges from the intent.
- **Security**: injection, authorization gaps, secret handling, unsafe input. Apply only where the diff touches such code.
- **Design**: API contracts, coupling, consistency with the repository's existing patterns and documented conventions (`AGENTS.md`, `CLAUDE.md`, contributor docs).
- **Tests**: is the new behaviour covered; do the tests assert the thing that matters.

The lenses are independent. When the host supports subagents, run them in parallel, each with the scope summary and the list of touched files; otherwise run them one after another yourself.

Done when every touched file has been read under every lens that applies to it.

## 4. Validate

Re-check every candidate finding against the actual code: the line exists, the claim holds, the fix is plausible. Drop findings you cannot confirm and merge duplicates. Rank by severity: **blocking** (wrong behaviour, security hole, data loss), **should fix**, **nit**.

## 5. Report

Save the review to `docs/reviews/<date>-<target>.md` unless the user names another place; preserve existing files. For each finding give: severity, `path:line`, the problem, the evidence, and a suggested fix. Open with a verdict (approve / approve with changes / request changes) and a one-paragraph summary; close with what ran and what was not checked.

Return the verdict, the blocking findings, and the report path. Posting comments to the PR is a separate step the user asks for; when asked, post one review with `gh pr review` or inline comments through `gh api`.
