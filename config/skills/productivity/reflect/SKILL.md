---
name: reflect
description: Look back on finished work (a session, an epic, or recent PRs) and turn its lessons into edits to skills, instructions, docs, ADRs, REVIEW.md and checks, landed as one PR. Use when the user asks for a retrospective, or from a scheduled routine.
---

# Reflect

Find **lessons**: places where the agent's environment (skills, instructions, docs, checks, tooling) made finished work slower or worse than it had to be. A lesson ends as an edit to a file, never as a note. `debrief` records friction per objective; `reflect` reads across many objectives and fixes what keeps recurring.

## 1. Pick the scope

The user names a session, an epic or a time range; with none, use the current session. A routine run with no scope takes the PRs merged since the last `reflect` PR.

Done when you can name the work in scope and where it was tracked.

## 2. Read the sources

Read every source in scope that you can reach, not a sample. List the ones you could not reach; they go in the report.

1. **Tracker.** The epic or issues and their comments, as `docs/agents/issue-tracker.md` says. Open every doc they link through the matching adapter (`docs/agents/product-docs.md` for a docs store): specs, briefs, research, handovers. Read handovers and results first.
2. **Committed debrief output.** What `debrief` filed in the repo: `docs/agents/*.md` and `docs/adr/` entries changed by the work in scope.
3. **PRs.** Each PR in scope with its review comments and reviews (`gh pr view <n> --comments`, `gh api repos/{owner}/{repo}/pulls/<n>/reviews`), its CI failures, and reverts or fixup commits in git history.
4. **Local only, when present.** Private debrief notes under `<store>/notes/` (the store path is in "Project knowledge" in your user instructions); session transcripts: Claude Code under `~/.claude/projects/<project>/*.jsonl` (`<project>` is the repo path with `/` replaced by `-`), Codex under `~/.codex/sessions/`.

Done when every reachable source is read and the unreachable ones are listed.

## 3. Collect friction

Walk the sources in order and note each moment the work went sideways: a long search for a file, a mistake a person or reviewer caught, a step a skill got wrong or left out, a rule someone restated by hand, a repeated review finding, an expensive tool call, information the agent could not reach. Quote the moment and say where it happened. Friction seen in two or more objectives ranks first.

## 4. Turn friction into lessons

Pick the cheapest target that stops each friction next time ([where a lesson goes](#where-a-lesson-goes)). Read the target file before writing the edit, so the edit fits what is there and does not repeat it. Drop friction no edit would have prevented. Write each edit in the `writing-for-agents` style when that skill is installed.

For each lesson give: the target path, the exact change (a diff or the new text), and one line on the friction it prevents. Order them by cost of the friction, worst first.

## 5. Land them

- **With a person**: show the list and ask which edits to apply. Apply those on one branch and open one PR with `publish-pr`. Edits to the vision, the glossary (`CONTEXT.md`) or a project-wide decisions doc are listed under **Reviewer notes** so the reviewer approves each one.
- **Routine run**: open no PR. Post the list as a comment on the epic or tracker issue in scope, or as the run's final message, with the unreachable sources named.

Done when every lesson is in the PR, declined by the user, or posted.

## Where a lesson goes

- **Skill**: a step was missing, wrong or ambiguous while running a skill. Fix the step or its completion criterion in place.
- **Navigation**: the agent took long to find something. Add a pointer to the docs catalog (`docs/agents/README.md`) or the doc the agent was already reading. `AGENTS.md` and `CLAUDE.md` load every turn, so they hold pointers only.
- **`docs/agents/`**: the agent misused the tracker, the docs store, labels or the domain docs. Fix the adapter file.
- **ADR** (`docs/adr/`): a technical choice was re-argued or made inconsistently. Propose a new ADR or an amendment.
- **Glossary, vision, decisions**: a term was used two ways, scope drifted, or a settled rule was reopened.
- **Automated check**: the mistake is mechanical (a banned API, an import shape, a file location). A lint rule, pre-commit hook or CI job beats any written rule. Read the repo's existing checks first; a check that exists but is unwired is the lesson.
- **Review rule**: a judgement call a reviewer should have caught. Add it to `REVIEW.md` at the repo root, creating it if missing; `review-pr` reads it.
- **Tooling and access**: an expensive or repeated tool call, or information the agent could not reach. Propose the script, connector or read access.
- **Pruning**: an instruction changed nothing or contradicted another. Propose deleting it.
