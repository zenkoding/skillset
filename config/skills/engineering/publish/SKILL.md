---
name: publish
description: Publish a Git branch as a GitHub pull request, or update its existing PR, when the user requests publication.
compatibility: Requires Git, authenticated gh, repository access, and project-required checks.
---

# Publish a PR

Publish the intended branch with accurate metadata and verified changes. Infer the branch, base, and useful preparation from the request and repository. It is OK to ask the user to choose between concrete options when a consequential ambiguity remains. Preserve explicitly supplied choices; do not require command-style parameters.

## Prepare

Resolve the head repository, branch and SHA, target base, and existing PR. A branch's upstream can be its published head rather than its PR base. Use an explicit base, an existing PR's base, or the verified remote default; clarify any remaining ambiguity before publication.

Inspect the relevant diff/history, project instructions, and dirty/index state. Run required checks. Keep unrelated edits and the index unchanged when publishing committed work; temporary PR metadata is never part of the staged change. Use an isolated checkout when necessary.

When creating a PR whose diff crosses the [reviewer notes](references/reviewer-notes.md) threshold, offer the full notes; add them when the user says yes or asked for them.

Compose only missing title/body fields using the [authoring guidance](references/pr-authoring-format.md). Preserve supplied text and treat input files as read-only. Prepare multiline text in a separate temporary body file. If approval is still needed, present the concrete result; reuse authorization already given for this publication.

## History and publication

Keep history intact unless cleanup serves the user's request. For authorized cleanup, preserve a backup ref and verify that history-only edits leave the tip tree unchanged. A rejected push is not permission to force: rewriting a published branch requires explicit authorization and an exact expected-SHA lease. Never force-update the default or base branch.

Push to the verified head remote/ref. Reuse a matching existing PR and keep its base; otherwise create the PR as a draft with explicit head/base and a body file. When updating a PR, refresh its Reviewer notes against the final diff. Do not merge, publish a release, or post unrelated comments as part of opening a PR.

## Ready for review

A new PR goes from draft to ready only once it checks out:

1. **Verify.** Re-read the title and body against the final diff and the repo's conventions ([authoring guidance](references/pr-authoring-format.md)), and fix what is off.
2. **Wait for CI**, when the repo has any, on the pushed head. Pending is not green. When it fails, stop here: the PR stays a draft, and the result names the failing check.
3. **Flip and label.** Mark the PR ready for review and add the review label the repo uses: `ready-for-human` when the Merge danger is a one-way door or the user asked for a person's review, otherwise `ready-for-agent`. The repo's triage-label docs win over this default.

An existing PR keeps its draft or ready state unless the user asks to change it.

Verify the resulting PR identity and metadata. Clean up only run-owned temporary files using ordinary filesystem operations. If cleanup is unavailable, leave and report them; never reset, stash, force-checkout, or clean the user's checkout to remove PR metadata. Recheck that unrelated edits and index state are preserved.

Return the PR link, its draft or ready state and label, CI status, changes published, checks actually run, history changes if any, and remaining blockers. A failed or denied operation remains incomplete.
