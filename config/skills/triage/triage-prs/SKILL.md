---
name: triage-prs
description: Go through open pull requests, merge the ones that pass the merge gate, label the ones that need a person, and name the next step for the rest. Use when asked to triage PRs, check the PR queue or land what's ready, or from a PR triage routine.
argument-hint: "[owner/repo] [#n ...] [--dry-run]"
compatibility: Requires authenticated gh with permission to label and merge PRs, and a readable checkout of the repo.
---

# Triage PRs

Keep the PR queue moving without a person where it is safe: merge what passes the merge gate, hand every other reviewed PR to a person with the reason, and name the skill that unblocks the rest.

Arguments: the repo (default: the current clone), PR numbers to limit the run, and `--dry-run` to print the verdicts without changing anything. Invoking this skill authorizes the labels, comments and merges the [triage rules](references/triage-rules.md#pull-requests) allow, for this run.

## 1. Load the rules and the queue

Read `docs/contributing/issues.md` for the label strings and `REVIEW.md` or `CONTRIBUTING.md` for anything the repo adds to the merge rules: a required approver, protected paths, a merge freeze. The repo's rules can only tighten the gate, never loosen it.

List open PRs with what the state table needs, in one call:

```bash
gh pr list --state open --limit 200 \
  --json number,title,author,isDraft,baseRefName,headRefOid,mergeable,reviewDecision,statusCheckRollup,labels,url
```

Find the default branch with `gh repo view --json defaultBranchRef`.

Done when you hold the rules and the list to triage.

## 2. Judge each PR

Classify each PR by the [state table](references/triage-rules.md#state-first-match-wins). Fetch per-PR detail only for PRs that reach the review rows: `gh pr view <n> --json reviews,comments,commits` and the unresolved review threads.

A review counts for the head commit only when it was submitted on that commit. Read its findings from its body: `review` posts one review carrying the `<!-- review-pr -->` marker, with each finding labelled `Crucial`, `Important`, `Suggestion`, `Question`, `Hint` or `Remark`.

For every **Reviewed and green** PR, run the [merge gate](references/triage-rules.md#merge-gate): size the diff with [work sizing](references/work-sizing.md) (`gh pr diff <n>` and the PR's What), check its door, then each remaining line. Record which lines held and which failed.

Done when every PR has a state and every reviewed, green PR has a gate result.

## 3. Act

With `--dry-run`, skip to step 4. Otherwise:

- **Gate passed**: comment with the gate lines, then `gh pr merge <n> --squash --auto` (or `--squash` when the repo doesn't allow auto-merge). Remove `ready-for-human` if it was set.
- **Gate failed**: add `ready-for-human` and comment with the line that failed and why, in a sentence.
- **Unreviewed**: add `ready-for-agent` when the repo uses it to start reviews.
- **Draft, stacked, broken, unaddressed, waiting**: change nothing. They appear in the report with their next step.

Comment once per PR per head commit: when a triage comment for the current head is already there, don't add another. Never push to a PR, approve it, or close it. When a merge fails, report the error and leave the PR as it was.

Done when every passed PR is merged or queued for auto-merge, and every failed one is labelled with its reason.

## 4. Report

Print one line per PR, grouped: merged, needs a person, needs work, waiting. Each line has the number, title, state and, for anything not merged, the next step as a command (`/address 42`, `/review 57`) or the failed gate line. End with the counts per group.
