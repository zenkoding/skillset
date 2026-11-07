---
name: find-bugs
description: Hunt for defects in a code area (broken edge cases, contract mismatches, swallowed errors, races) and file each evidenced one as an issue. Use when asked to find bugs, audit an area for defects, or from a bug-hunting routine; not for diagnosing one known bug (that's fix-bug).
argument-hint: "[path or area] [--budget N]"
compatibility: Requires a readable checkout and the project's test tools; authenticated gh or the tracker CLI for filing. Parallel subagents are optional.
---

# Find bugs

Find defects nobody has reported yet, prove each one, and file it so the quality lane can fix it. A short list of reproduced bugs beats a long list of suspicions.

Arguments: the area to hunt in (default: the whole repo, most-changed code first) and a budget (default 5, per the [find rules](references/find-rules.md#budget)).

## 1. Pick where to look

Read `CONTRIBUTING.md` and the docs it points to for the area, `CONTEXT.md`, and the contracts the code promises: public APIs, CLI flags, config formats, documented behaviour. Rank the area's files by risk: churn (`git log --since=90.days --name-only`), recent bug fixes (a fixed bug often has siblings), size, and how many callers depend on them.

Done when you have a short, ranked list of files or flows to read.

## 2. Hunt

Read each target looking for behaviour that differs from what the code or docs promise:

- **Boundaries**: empty, zero, one, max, unicode, very large inputs; off-by-one in loops and slices.
- **Errors**: swallowed or ignored errors, wrong error returned, cleanup skipped on failure, partial writes.
- **Contracts**: a caller relying on behaviour the callee doesn't guarantee; docs and code that disagree.
- **State**: shared state without synchronisation, ordering assumptions, retries that aren't idempotent.
- **Siblings**: the same pattern as a recently fixed bug, elsewhere.

For a large area, when the host supports subagents, give each one a slice of the ranked list and gather their candidates.

Done when every target has been read under every heading that applies.

## 3. Prove

For each candidate, write the smallest failing check you can: a test, a CLI call, a request. Run it. A candidate that doesn't fail is dropped or kept as "suspected" in the report, never filed. Keep the failing check's code and output for the issue's Evidence; don't commit it.

Done when every filed candidate has a failing check, run.

## 4. File and report

File each proven bug per the [find rules](references/find-rules.md#filing), with the lane label `bug`, the failing check in Evidence and a severity. Then report per the find rules.
