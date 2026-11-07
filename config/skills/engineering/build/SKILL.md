---
name: build
description: "Build one tracker ticket from its plan to an open PR a reviewer can check from the body alone. Run as /build <ticket ID, URL or description>."
disable-model-invocation: true
---

Build one **ticket** (whatever the tracker calls a unit of work: a feature, task, story or bug) from its **plan** to a PR into the default branch.

`docs/agents/issue-tracker.md` says how to read and change the ticket, where the plan lives, which plan sections you may edit, and the delivery rules. Use its operations for every tracker step below.

## Steps

1. **Find the ticket.** The argument can be:
   - **an identifier or a tracker URL**: read that ticket;
   - **free text**: search the tracker for matching tickets and show the best few. Ask the user to confirm which one they meant. When nothing matches, ask what they meant and whether it is new work; for new work, create the ticket first, then go on with it.
   - **nothing**: ask for one of the above.

   Done when exactly one ticket is confirmed.

2. **Claim.** Read the ticket. State its title and its parent epic when it has one. Go on when it is buildable work (a feature, task, story or bug); for anything else, such as an epic, stop and say what it is. Go on when the parent epic is out of Backlog (its spec is approved); while it is in Backlog, stop and report that. Assign yourself. When the tracker moves the ticket to In Progress once a branch carrying the identifier exists (step 3), check it moved and set it by hand only if it didn't. Done when the tracker shows In Progress.

3. **Branch.** Work on a branch whose name carries the identifier in lower case (the tracker's suggested branch name works). Keep the current branch if it already carries the identifier (a worktree session often starts on one); otherwise cut one from the current default branch. Done when `git branch --show-current` contains the identifier.

4. **Read the plan.** Read every section of the plan, the parent epic's description and the docs it links. Then open the repo docs that apply to this change from the docs catalog (`docs/agents/README.md`): usually the decisions doc, the architecture doc, ADRs in the area, and the glossary for names. Done when each **Done when** item maps to a check you know how to run.

5. **Resolve open questions.** Every item under the plan's **Open questions**, and every gap step 4 found, gets an answer before code:
   - Answer it from the code, the repo docs or the epic's spec when they settle it. Write the answer and its source under **Open questions**.
   - When only a person can answer it, or the answer would change **What to build** or **Done when**, add it under **Open questions**, comment the question on the ticket, and stop. The ticket stays In Progress and assigned to you, with the question as its last comment. Report the question.

   Done when every item under **Open questions** has an answer.

6. **Build test-first.** Read the plan's **Flags**; if it names a flag, put unfinished work behind it. Otherwise keep unfinished UI unlinked, so the default branch stays deployable after every PR. Work test-first with `tdd`. The **Done when** checks are the pre-agreed seams; a seam the plan does not name goes through step 5 first. Tick each **Done when** item on the ticket once its check passes. Done when every item is ticked, or noted for the step 9 report with why it can't be checked from the branch.

7. **Check.** Run the check command that `AGENTS.md` names (lint, tests, build) and fix what it reports until it exits 0. When it fails on code this ticket didn't touch, stop and report the failure instead of fixing it. Done when the last run exits 0.

8. **Open the PR.** Commit, push, and open one PR with `publish`. `publish` writes the title; its **Why** links the ticket, and its **How to test** starts from the plan's. The branch name carries the identifier, so trackers that link by branch still link the PR. Add the PR link to the ticket. Done when the PR is open and its Why links the ticket.

   A ticket whose diff would be large ([work sizing](references/work-sizing.md)) becomes a stack: each PR into its own base, each leaving the default branch deployable, each linked on the ticket. The ticket goes In Review when the first PR opens and Done when the last merges.

9. **Hand to review.** `publish` opens the PR as a draft and flips it to ready, with its review label, once it checks out and CI is green. When the tracker moves the ticket to In Review as the PR opens, check it moved and set it by hand only if it didn't. Report the PR URL, any **Done when** item still unticked and why, and anything not tested. Done when the tracker shows In Review.
