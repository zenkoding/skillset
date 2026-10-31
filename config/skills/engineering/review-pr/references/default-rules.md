# Default review rules

`review-pr` applies these when the repo has no `REVIEW.md`. A repo's `REVIEW.md` replaces the section it names and keeps the defaults for the rest.

## What blocks a merge

- Wrong behaviour against the stated intent, or a regression to shipped behaviour callers rely on.
- A security hole: injection, missing authorization, a leaked secret or token, unsafe input handling.
- Data loss or an irreversible change without a way back.
- Behaviour that changes inside a change described as a refactor.
- A call to a function, field, flag or option that does not exist.
- A change that does not belong in the project, per its vision or non-goals.

## Do not report

- Anything CI already enforces: formatting, lint, the test suite, type checks.
- Generated or vendored files, lockfiles.
- Missing tests in files the diff did not touch.
- Style preferences no written rule backs.

## Nits

At most three `Suggestion` findings about naming, structure or duplication per review; fold the rest into one line.

## REVIEW.md format

A repo's `REVIEW.md` sits at the repo root and uses these headings, each optional:

```markdown
# Review rules

## What failure looks like here
<what blocks a merge in this repo, beyond or instead of the defaults; e.g. "for docs: it lies; for config: it misroutes">

## Do not report
<paths to skip, checks CI already runs>

## Extra checks
<repo-specific checks, each tied to one of the four review areas>
```
