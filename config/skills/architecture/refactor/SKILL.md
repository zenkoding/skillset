---
name: refactor
description: Take one architecture issue (usually from find-drift) and bring the code back in line with the documented architecture without changing behaviour; open a PR when the move is trivial, hand it to brief when it isn't. Use when asked to fix a drift or architecture issue.
argument-hint: "<issue>"
compatibility: Requires a writable checkout, the project's build and test tools, and authenticated gh or the tracker CLI.
---

# Refactor

Fix one architecture deviation so the code matches the rule it broke, with behaviour unchanged. The tests prove behaviour; the quoted rule proves the direction.

Arguments: the issue. Invoking `refactor` on an issue authorizes the branch and PR for a trivial move.

## 1. Read the issue and its rule

Read the issue with its comments through `docs/contributing/issues.md`. Re-check the quoted rule against the current architecture docs and history: if the rule changed since the issue was filed, stop and comment that the issue is stale instead of refactoring.

Confirm the deviation still exists at the current commit.

Done when you can name the rule, the code that breaks it, and the target shape in a sentence each.

## 2. Size it

Size the move with [work sizing](references/work-sizing.md).

- **Trivial** (every axis low, two-way door, no public contract touched): go on.
- **Not trivial**: hand the issue to `brief` with the rule and the target shape, and stop. A move across components, a changed public API or a data migration always plans first.

## 3. Move, behaviour unchanged

1. Run the tests that cover the code you'll touch. When coverage is thin, add characterization tests first and commit them on their own.
2. Make the move in small steps, running the tests after each: rename, move, invert a dependency, extract.
3. Leave behaviour alone. A bug found on the way becomes its own issue, not part of this change.

Done when the deviation is gone, the tests pass, and nothing outside the move changed.

## 4. Publish

Open the PR with `publish`, linking the issue. Its Why quotes the rule; its How to test names the tests that prove behaviour held. Add a `Codify:` line to the PR's Reviewer notes when a check could stop this drift from coming back.
