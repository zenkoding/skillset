# Work sizing

Shared criteria for how risky, complex and taste-dependent a piece of work is. Several skills carry this same file, so they size work the same way. Each term below has one definition; each skill decides what to do with it.

## The axes

Judge each axis on the change itself, not on its size. Before any code exists, predict it from the plan; once a diff exists, confirm it on the diff and the PR's **What**. A change can be high on more than one axis.

| Axis | High when the change… |
| --- | --- |
| **Risk** | touches a public contract or API (including CLI flags and config formats), auth, permissions, secrets or crypto, input parsing at a trust boundary, dependencies, a migration, schema or stored data, or code many callers share; or is a one-way door |
| **Complexity** | spans more than one component, adds an abstraction, or its What can't be said in three sentences |
| **Taste** | has more than one reasonable choice in UX, output format, naming or API shape |

## Door

- **One-way**: can't be cleanly reverted. Data migrations, deleted data, published API or schema changes, external side effects (messages sent, money moved, releases published).
- **Two-way**: a revert fixes it.

## Trivial

Work is **trivial** when its scope is bounded and every axis is low. Any high axis makes it **non-trivial**. When in doubt, call the axis high: a short grill costs less than rework.

A **defect** is a separate track: behaviour that differs from what the code or docs promise. Establish the cause first; the fix is then sized like any other work.

## Counted lines

Additions plus deletions in the diff, excluding lockfiles, generated or vendored files, pure renames and test fixtures. **Large** means 500 or more counted lines. Size is not an axis: it only sets review mechanics (splitting a diff, parallel review passes), never how dangerous a change is.

## Substantial

Work is **substantial** when it changed code, opened a PR or made a decision, or when it may need to resume in a later session. Substantial work keeps notes under `<store>/notes/<objective>/` when the user instructions define a private store (see "Project knowledge"); otherwise in a private location the host allows, recorded in the result. Question-and-answer sessions are not substantial.
