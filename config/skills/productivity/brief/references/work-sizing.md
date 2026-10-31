# Work sizing

Shared criteria for how big, risky and open-ended a piece of work is. Several skills carry this same file, so they size work the same way. Each term below has one definition; a skill decides what to do with it.

## Counted lines

Additions plus deletions in the diff, excluding lockfiles, generated or vendored files, pure renames and test fixtures.

**Large**: 500 or more counted lines.

## Risky area

A change touches a risky area when it changes any of:

- public contracts or APIs, including CLI flags and config formats;
- migrations, schema or stored data;
- auth, permissions, secrets or crypto;
- input parsing at a trust boundary;
- dependency manifests;
- code many callers share.

## Door

- **One-way**: hard to undo. Data migrations, deleted data, published API or schema changes, external side effects (messages sent, money moved, releases published).
- **Two-way**: a revert fixes it.

## Trivial

Work is **trivial** when all of these hold:

- bounded scope;
- one obvious approach;
- no architecture, UX/UI or taste decision;
- a two-way door;
- no risky area.

Anything else is **non-trivial**. When in doubt, call it non-trivial: a short grill costs less than rework.

A **defect** is a separate track: behaviour that differs from what the code or docs promise. Establish the cause first; the fix is then trivial or non-trivial like any other work.

## Substantial

Work is **substantial** when it changed code, opened a PR or made a decision, or when it may need to resume in a later session. Substantial work keeps notes under `<store>/notes/<objective>/` when the user instructions define a private store (see "Project knowledge"); otherwise in a private location the host allows, recorded in the result. Question-and-answer sessions are not substantial.

## How skills use it

| Skill | Uses |
| --- | --- |
| `brief` | Trivial, non-trivial or defect is the call. |
| `publish-pr` | Review verdict: fast track when trivial and not large; careful review on a risky area or a one-way door; otherwise quick review. Full reviewer notes when large or risky. |
| `review-pr` | Parallel finders when large; a separate security pass on a risky area; a one-way door needs a rollback before the verdict can be `merge`. |
| `build-feature` | A large feature becomes a stack of PRs. |
| `debrief`, `fix-bug`, `research-scope` | Notes when the work is substantial. |
