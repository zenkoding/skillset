# PR authoring guidance

Follow the repository's PR template and conventions. When none exist, use the template below. A reviewer should be able to answer three questions from the body alone: what changed, why, and how to check it.

- Base the description on the actual changes and tests, not the branch name or diff statistics.
- Preserve the user's supplied title and body. If they conflict with a required repository rule, explain the specific conflict rather than silently rewriting them.
- Distinguish passed checks, failures, skipped checks, and work still pending. Do not invent test evidence or claim a PR is merged or approved without confirmation.
- Include related issue links only when they are relevant and verified. Use no placeholder references.
- Add attribution only when required by the applicable host or project instructions; use the actual authoring host and avoid duplicate footers.

## Template

```markdown
### What

<the smallest visual that makes the change clear, then one or two sentences; before/after evidence>

### Why

<the problem; link the issue, tracker plan, RFC or ADR that asked for or settled it>

### How to test

1. <steps a reviewer runs on a fresh checkout to see the "after">

### Reviewer notes

**Merge danger**: <one-way | two-way> door: <why>. Blast radius: <what else this can break>.
```

### What

Open with the smallest visual that makes the change clear, then one or two sentences. Pick one:

- logic or an algorithm as pseudocode
- runtime control flow as a call tree
- UI structure as a component tree, with the state and module boundaries that matter
- file responsibilities or a broad refactor as a shallow file tree
- interaction between components as a Mermaid sequence diagram
- a `diff` sketch when the shape already exists and the point is what changes

```diff
 submitForm
   createSession
     persistPrompt
+    expandSkillMention
     launchAgent
```

Keep only the calls, files, props and states the reviewer needs. Then the before/after evidence, best first: a screenshot when the change is visible; command output or the test that failed before and passes now; for a refactor with no visible change, the tests that show behaviour stayed the same.

The menu of visuals comes from Dex Horthy's [`show-me`](https://github.com/humanlayer/skills/blob/main/plugins/show-me/skills/show-me/SKILL.md) skill (HumanLayer), adapted to a diff.

### Why

What the change is for, in a line or two. Link what asked for it: the issue, or the tracker's plan when the repo uses one (`docs/agents/issue-tracker.md` says how to link it; put the tracker identifier in the title or branch when the tracker links PRs that way). Link the RFC or ADR that settled the approach. In a stack, say which part this PR delivers and what is left.

### How to test

Numbered steps a reviewer can run on a fresh checkout: commands, URLs, seed data, flags. Start from the plan's test steps when there is one. For a bug fix, the first steps reproduce the bug on the base branch. Then what you ran and its result, and what was not tested.

### Reviewer notes

Always carries **Merge danger**:

- **Door**: one-way or two-way, as [work sizing](work-sizing.md) defines it.
- **Blast radius**: one line on what else this can break.

Large or risky PRs add the full notes; see [reviewer-notes.md](reviewer-notes.md).

## Before publication

Check that the title and body describe the final scope, preserve supplied text, and contain no stale claims from earlier attempts. Use a separate temporary file for multiline text and pass it as data through `gh --body-file`; caller input files remain unchanged.
