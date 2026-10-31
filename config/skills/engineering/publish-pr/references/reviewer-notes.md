# Reviewer notes

`### Reviewer notes` is the last section of a PR body. It tells the reviewer where the risk is and what to confirm, so attention goes to the few files and decisions that matter. Every PR carries **Merge danger** ([authoring guidance](pr-authoring-format.md)); large or risky PRs add the full notes below.

## When to offer the full notes

Offer only when creating the PR, and add them only on a yes. Updates refresh existing full notes but never offer them. An explicit request ("publish with full reviewer notes") adds them at any size, on create or update. Offer when the diff is **large** or touches a **risky area**, as [work sizing](work-sizing.md) defines them.

The offer is one line: "This PR is large/risky (<reason>). Add full reviewer notes?"

## Template

Fill every part that has something real to say; drop the rest. Merge danger stays.

```markdown
### Reviewer notes

**Verdict**: ⚠ careful review: <one line on where the risk is>

**Merge danger**: <one-way | two-way> door: <why>. Blast radius: <what else this can break>. Rollback: <how>.

**Main decisions**
- <chose X over Y because Z> ([ADR-NNNN](https://github.com/<owner>/<repo>/blob/<head-branch>/docs/adr/NNNN-....md))

**Contentious points**
- <unusual choice a reviewer will ask "why like this?" about>

**Must check**
- [ ] <breaking contract / migration / security-sensitive change / dep bump>: confirm <what>

**Where to focus**
- High: `path` <why>
- Medium: `path` <why>
- Skim: <the rest, summarised>
```

How the change works belongs in **What**, and the hand-test steps in **How to test**; the notes don't repeat them.

## Verdicts

Pick by [work sizing](work-sizing.md), first match wins:

- **⚠ careful review**: a risky area, a one-way door, a change to behaviour callers rely on, or a gap in the tests.
- **✓ fast track**: trivial and not large; reading the diff once is enough.
- **◐ quick review**: everything else: real logic, contained blast radius, covered by tests.

## Rules

- **Main decisions** quotes the session's `debrief` Decisions bullets when there are any, and links `docs/adr/` entries the branch adds. Never link `.docs/` paths or the private store (`~/.local/share/agents/`): they are private and absent from the repo.
- Every path in Where to focus is in the PR diff.
- Links in the notes are absolute `https://github.com/<owner>/<repo>/blob/<head-branch>/` URLs; relative links in a PR body resolve under `/pull/` and 404.
- Every Must check item names what to confirm.
- When `publish-pr` updates the PR, rewrite the notes against the final diff; the rest of the body follows the usual preservation rules.
