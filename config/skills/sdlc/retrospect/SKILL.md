---
name: retrospect
description: Look back on a closed epic or milestone (or a session or recent PRs) and turn its lessons into fixes to the process (skills, docs, checks) and to the product (vision, briefs, backlog). Offered when an epic or milestone closes; also run by hand or from a routine.
---

# Retrospect

A product-lane ceremony: the retro at the end of an epic or milestone. It finds two kinds of **lessons**:

- **Process**: places where the agent's environment (skills, instructions, docs, checks, tooling) made the work slower or worse than it had to be.
- **Product**: what the shipped work taught about the product: a goal missed or met differently, scope that drifted, a non-goal that turned out wrong, an idea for what comes next.

A lesson ends as an edit or a tracker item, never as a note. `debrief` records friction per objective; `retrospect` reads across the objectives of an epic or milestone and fixes what keeps recurring.

## 1. Pick the scope

The usual scope is an epic or milestone that just closed: `debrief` offers a retro when its objective finished an epic's last feature, and `scope-milestone` offers one for the previous milestone before scoping the next. A user may also name a session, an epic, a milestone or a time range; with none, use the current session. A routine run with no scope takes the PRs merged since the last `retrospect` PR.

Done when you can name the work in scope, where it was tracked, and the goal it set out to meet (the epic spec's or milestone brief's "done when").

## 2. Read the sources

Read every source in scope that you can reach, not a sample. List the ones you could not reach; they go in the report.

1. **Tracker.** The epic or issues and their comments, as `docs/agents/issue-tracker.md` says. Open every doc they link through the matching adapter (`docs/agents/product-docs.md` for a docs store): specs, briefs, research, handovers. Read handovers and results first.
2. **Project knowledge.** The public catalog (`docs/agents/README.md`, `CONTEXT.md`, `docs/adr/`) and what changed in it during the work in scope, the vision doc and the decisions doc when the repo has them, plus any debrief notes a user asked to commit (search the repo for `### Debrief:` headings).
3. **PRs.** Each PR in scope with its review comments and reviews (`gh pr view <n> --comments`, `gh api repos/{owner}/{repo}/pulls/<n>/reviews`), its CI failures, and reverts or fixup commits in git history.
4. **Local only, when present.** The private store's catalog and its briefs and debriefs under `<store>/notes/` (the store path is in "Project knowledge" in your user instructions); session transcripts: Claude Code under `~/.claude/projects/<project>/*.jsonl` (`<project>` is the repo path with `/` replaced by `-`), Codex under `~/.codex/sessions/`.

Done when every reachable source is read and the unreachable ones are listed.

## 3. Collect friction

Walk the sources in order. For the product, compare what shipped with the goal from step 1: each "done when" check met, partly met or missed, scope added or cut on the way, and what users, reviewers or tests showed that the plan didn't expect. For the process, note each moment the work went sideways: a long search for a file, a mistake a person or reviewer caught, a step a skill got wrong or left out, a rule someone restated by hand, a repeated review finding, an expensive tool call, information the agent could not reach. Quote the moment and say where it happened. Friction seen in two or more objectives ranks first.

## 4. Turn friction into lessons

Pick the cheapest target that stops each friction next time, or carries each product finding forward ([where a lesson goes](#where-a-lesson-goes)). Read the target file before writing the edit, so the edit fits what is there and does not repeat it. Drop friction no edit would have prevented. Write each edit in the `writing-for-agents` style when that skill is installed.

For each lesson give: the target path, the exact change (a diff or the new text), and one line on the friction it prevents. Order them by cost of the friction, worst first.

## 5. Land them

- **With a person**: show the list, process and product apart, and ask which to apply. Apply the process edits on one branch and open one PR with `publish-pr`. Product lessons go where the product lane keeps them: a vision change through `envision`, input for the next milestone into its brief through `scope-milestone`, new work as Backlog issues in the tracker. Edits to the vision, the glossary (`CONTEXT.md`) or a project-wide decisions doc are listed under **Reviewer notes** so the reviewer approves each one.
- **Routine run**: open no PR. Post the list as a comment on the epic or tracker issue in scope, or as the run's final message, with the unreachable sources named.

Done when every lesson is in the PR, handed to the product lane, declined by the user, or posted.

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

Product lessons:

- **Vision**: a goal, audience or non-goal the work proved wrong. Hand it to `envision`.
- **Next milestone**: a "done when" missed or cut, or a gap the work exposed. It becomes input for the next milestone brief via `scope-milestone`.
- **Backlog**: an idea or follow-up the work turned up. File it as a Backlog issue as `docs/agents/issue-tracker.md` says, after the user agrees.
