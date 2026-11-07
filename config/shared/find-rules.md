# Find rules

How the lane skills that go looking for work (`find-bugs`, `find-vulns`, `find-drift`, `find-hotspots`) report what they find. Each skill decides what counts as a finding; this file holds what they share.

## A finding

A finding is worth filing only when you can point at it and show it matters:

- **Located**: files and lines, or the endpoint, query or benchmark, at a recorded commit.
- **Evidenced**: a failing input, a trace, a quoted rule it breaks, or measured numbers. "This looks fragile" is not a finding.
- **Bounded**: one problem, fixable on its own. Split a finding that needs two unrelated fixes.

Drop what you can't locate or evidence. Tag a finding that only a deterministic check could catch everywhere as a **Codify candidate** (see [work sizing](work-sizing.md)).

## Budget

File at most **5** findings per run unless the user names another number, highest severity first. List the rest in the report with one line each, so the next run or a person can pick them up.

## Before filing

Search the tracker for open and recently closed issues about the same location or symptom. A match gets a comment with the new evidence instead of a new issue. A match closed as `wontfix` is not filed again.

## Filing

File each finding through the operations in `docs/contributing/issues.md`, labelled with the needs-triage role and the skill's lane label (the tracker doc's strings win). The body:

```markdown
## What
<the problem in two or three sentences, for someone new to the area>

## Evidence
<location at <commit>, and the failing input, trace, quoted rule or numbers>

## Severity
<critical | high | medium | low>: <why, in a sentence>

## Suggested fix
<the smallest change you'd try first, or "unknown">
```

`triage-issues` (when installed) judges what you filed; don't label findings ready yourself.

## Report

End the run with the filed issues (number, title, severity), the findings that matched an existing issue, the ones over budget, and anything you looked at and ruled out, in a line each.
