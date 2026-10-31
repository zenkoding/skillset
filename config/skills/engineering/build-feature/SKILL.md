---
name: build-feature
description: "Build one tracker feature from its plan to an open PR a reviewer can check from the body alone. Run as /build-feature <feature id>."
disable-model-invocation: true
---

Build one **feature** from its **plan** to a PR into the default branch. The feature identifier is the argument; if none is given, ask for it.

`docs/agents/issue-tracker.md` says how to read and change the issue, where the plan lives, which plan sections you may edit, and the delivery rules. Use its operations for every tracker step below.

## Steps

1. **Claim.** Read the feature issue. State its title and parent epic. Go on when it is a feature or bug issue; for anything else, stop and say what it is. Go on when the parent epic is out of Backlog (its spec is approved); while it is in Backlog, stop and report that. Assign yourself. When the tracker moves the issue to In Progress once a branch carrying the identifier exists (step 2), check it moved and set it by hand only if it didn't. Done when the tracker shows In Progress.

2. **Branch.** Work on a branch whose name carries the identifier in lower case (the tracker's suggested branch name works). Keep the current branch if it already carries the identifier (a worktree session often starts on one); otherwise cut one from the current default branch. Done when `git branch --show-current` contains the identifier.

3. **Read the plan.** Read every section of the plan, the epic's description and the docs it links. Then open the repo docs that apply to this change from the docs catalog (`docs/agents/README.md`): usually the decisions doc, the architecture doc, ADRs in the area, and the glossary for names. Done when each **Done when** item maps to a check you know how to run.

4. **Resolve open questions.** Every item under the plan's **Open questions**, and every gap step 3 found, gets an answer before code:
   - Answer it from the code, the repo docs or the epic's spec when they settle it. Write the answer and its source under **Open questions**.
   - When only a person can answer it, or the answer would change **What to build** or **Done when**, add it under **Open questions**, comment the question on the issue, and stop. The issue stays In Progress and assigned to you, with the question as its last comment. Report the question.

   Done when every item under **Open questions** has an answer.

5. **Build test-first.** Read the plan's **Flags**; if it names a flag, put unfinished work behind it. Otherwise keep unfinished UI unlinked, so the default branch stays deployable after every PR. Work test-first (the `tdd` skill when installed). The **Done when** checks are the pre-agreed seams; a seam the plan does not name goes through step 4 first. Tick each **Done when** item on the issue once its check passes. Done when every item is ticked, or noted for the step 8 report with why it can't be checked from the branch.

6. **Check.** Run the check command that `AGENTS.md` or `CLAUDE.md` names (lint, tests, build) and fix what it reports until it exits 0. When it fails on code this feature didn't touch, stop and report the failure instead of fixing it. Done when the last run exits 0.

7. **Open the PR.** Commit, push, and open one PR with `publish-pr`, titled with the identifier first (`<ID>: <feature title>`) so the tracker links it. Its **How to test** starts from the plan's. Add the PR link to the issue. Done when the PR is open and its title carries the identifier.

   A feature whose diff would be large ([work sizing](references/work-sizing.md)) becomes a stack: each PR into its own base, each leaving the default branch deployable, each linked on the issue. The issue goes In Review when the first PR opens and Done when the last merges.

8. **Hand to review.** When the tracker moves the issue to In Review as the PR opens, check it moved and set it by hand only if it didn't. Mark the PR ready and add `ready-for-agent` when the repo uses `review-pr` routines. Report the PR URL, any **Done when** item still unticked and why, and anything not tested. Done when the tracker shows In Review.
