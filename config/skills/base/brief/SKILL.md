---
name: brief
description: Turn a high-level problem, idea or ticket into a plan the user has signed off on, so build has something concrete to work against. Researches, drafts, grills, writes the plan and its companion, then asks for sign-off. Run as /brief <idea | ticket | problem>.
---

# Brief

The start-of-work counterpart of `debrief`, and the planning orchestrator. A brief takes a high-level problem and ends with a **plan** the user has signed off on, plus a **companion** for the agent that builds it. The ask can be code work or anything else (a doc, a decision, an investigation); the code-only steps say so.

It runs the [planning loop](references/planning-loop.md) at ticket level; the output is a plan in the [plan template](references/plan-template.md).

```
/brief <problem>
  1 read the ask and its context
  2 research the big unknowns       predict the sizing → route
  3 draft the plan in chat          short, to argue with
  4 grill the open decisions        settle the sizing
  5 write the pair                  plan      → docs store, linked on the issue
                                    companion → <store>/notes/<objective>/
  6 sign-off                        the user approves the plan → build
```

## 1. Read the ask

- **Ticket** (an identifier or URL): read the issue, its comments, its parent, and every doc it links. `docs/agents/issue-tracker.md` says how to reach the tracker, and `docs/contributing/docs.md` where docs live.
- **Idea or problem** (free text): take it as given. Search the tracker for an existing issue on it and read that too.

Name the **objective** in a few words; it names the notes directory. Then read the context the loop's step 1 lists, plus the private store's catalog, and, for code work, the code the ask touches: entry points, callers, tests, recent history and recent PRs on the area.

Done when you can say what exists today and where the ask lands in it, with a `path:line` or link per claim.

## 2. Research and route

Research as the loop's step 2 says: answer small unknowns inline, and send an unknown that could change the plan and is too big to answer inline to `research`. Then predict the axes of [work sizing](references/work-sizing.md) and pick the route:

- **Defect**: behaviour differs from what the code or docs promise. Next is `fix-bug` when installed, to establish the cause before any fix; without it, find the cause before changing anything. The fix then comes back here, sized like any other work.
- **Trivial**: bounded scope, every axis low. Skip the grill. Write the **short plan** (Context, Proposal, Approach) on the issue, or in chat when there is none, and go to step 6.
- **Non-trivial**: any high axis. Run steps 3 to 6.

Tell the user the route and why in one line; a different call from the user wins.

## 3. Draft

Post the draft plan in chat, short enough to argue with: Context, Proposal and Approach in a few lines each, the predicted sizing, and the open decisions.

## 4. Grill

Run `grill` on the draft, seeded with the open decisions and each high axis. Then re-size: the settled axes and the why for each go in the plan's Risks section. When the grill turned the work trivial, say so and write the short plan instead.

## 5. Write the pair

- **The plan**, in the plan template, in the docs store where `docs/contributing/docs.md` says. Link it from the tracker issue; create the issue first when there is none and the user agrees. Without a docs store, the plan goes in the issue body.
- **The companion** at `<store>/notes/<objective>/companion.md`, by the absolute store path the session-start hook prints: the commit, code pointers, per-step instructions and tests, the sizing call, and the grill's decision log. On a resumed objective, update it in place.

The plan never links the companion. When they disagree, the plan wins; edit both in the same pass.

## 6. Sign-off

Show the user the plan and ask for sign-off. Apply wording rounds to both documents. On sign-off, the next step is `build` when installed (with the ticket), otherwise test-first (the `tdd` skill when installed); for work that isn't code, do it directly.

Done when the user has signed off the plan, the plan is linked on the issue, and the companion is written.
