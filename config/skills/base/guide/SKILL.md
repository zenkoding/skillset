---
name: guide
description: Explain how the installed skills fit together, or read the current branch, PR, ticket and notes and name the next skill to run, with the exact command. Run as /guide, /guide <question>, or /guide what next?
---

<!-- every skill named here: when installed -->
# Guide

Two jobs in one entry point: explain the workflow, and route to the next step. The guide only reads state; it never runs the skill it suggests.

## Which skills are installed

List the skill directories the host reads (`~/.claude/skills/` for Claude Code; `~/.agents/skills/` or `~/.codex/skills/` for Codex). Every answer names only those skills; a loop with none installed is left out, or mentioned once as a set the user could add.

## Explain

`/guide`, or `/guide <question>` ("how does review work?"): answer from the [workflow map](references/workflow-map.md), filtered to what's installed. Lead with the loop the question is about, as the map's diagram, then the one or two skills it asks about in a few lines each. Point at a skill's own `SKILL.md` for detail rather than restating it.

## Route

`/guide what next?`: read the state, then name one next step and the exact command, with one line on why.

Read, without asking:

- whether `CONTRIBUTING.md` and `docs/contributing/` exist;
- the current branch, and whether it is the default branch;
- the PR for this branch (`gh pr view --json number,state,isDraft,reviewDecision,labels,statusCheckRollup`), its reviews and unresolved review threads;
- the ticket in scope: the identifier in the branch name or PR body, read as `docs/contributing/issues.md` says;
- the objective's notes under `<store>/notes/<objective>/`: a companion, a grill log, a debrief.

Then take the first row that matches:

| State | Next |
| --- | --- |
| no `CONTRIBUTING.md` or `docs/contributing/` | `/setup-repo` |
| a defect report, cause unknown | `/fix-bug <issue>` |
| no ticket and no plan (default branch, or a fresh branch) | `/brief <problem>` |
| a ticket with no plan, or a plan not yet signed off | `/brief <ticket>` |
| a signed-off plan, no PR | `/build <ticket>` |
| a draft PR, or a PR whose CI is red | `/address <n>` |
| an open PR with no review on its head | `/review <n>` |
| review findings or change requests not yet addressed | `/address <n>` |
| reviewed, and the plan has How to test steps not yet run | `/test <n>` |
| approved and green | `/triage-prs <n>` merges it when it passes the merge gate, else a person does; then `/debrief` |
| merged, and it finished the last ticket of an epic or milestone | `/retrospect` |

When the skill a row names isn't installed, say what it would do and suggest the closest installed step, or doing it by hand.

Done when the user has the explanation they asked for, or one next step with its command and why.
