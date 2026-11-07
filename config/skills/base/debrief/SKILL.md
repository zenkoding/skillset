---
name: debrief
description: Debrief the user at the end of substantial work (code changed, PR opened, decision made) or when they ask for a debrief; records durable items in the project knowledge docs.
---

# Debrief

A debrief keeps the user's understanding in step with work an agent did, and turns repeated friction into fixes. It is a short chat message, then filing.

Run it once per objective, at the end of substantial work ([work sizing](references/work-sizing.md)), and whenever the user asks. Skip it for question-and-answer sessions and for follow-up fixes to work already debriefed: CI fixes, review replies, scheduled runs.

## 1. Write the debrief

When the objective has a `companion.md` from `brief` in its notes, read it and the plan it belongs to first: Decisions and Friction name where the work left the plan.

Post this in chat. Each section gets 1 to 3 bullets; drop empty sections.

```markdown
### Debrief: <objective>

**Decisions**
- Chose <X> over <Y> because <Z>. (<the doc it changed, or the PR whose body records it>)

**Learned**
- <fact about the code, tools or domain that was not obvious>

**Friction**
- <what slowed the work: missing doc, flaky tool, unclear instruction>

**Unusual choices**
- <what a reviewer would ask "why like this?" about>

**Worth understanding**
- `path:line`: <why it matters to the user's mental model> (1-2 items)

**Codify candidates**
- Codify: <the rule, one line>. Seen at <where the mistake happened>.

**Proposed fixes**
- <for friction seen before: the concrete change to a skill, instruction or doc>. Awaiting your OK.
```

- Decisions always read "chose X over Y because Z", self-contained, with no private note paths: the PR's reviewer notes (`publish` when installed) quote them verbatim.
- Every `path:line` in Worth understanding exists at the final commit; check it before posting.
- Friction that spans several objectives is a retro's job (`retrospect` when installed): say so in the debrief instead of proposing a one-off fix.
- When the objective finished the last open feature of an epic or milestone, offer a retro for it (`retrospect` when installed).
- A Codify candidate is a mechanical mistake from this objective that is **codifiable** ([work sizing](references/work-sizing.md)): one a reviewer, a test or the user caught, or one the agent made and fixed. List it; building the check is a retro's call (`retrospect` when installed).
- A friction item is **seen before** when an earlier note under `<store>/notes/` records the same friction: grep the notes for its key words. Only seen-before friction gets a Proposed fix.

## 2. File each item

Pick the destination per item. `<store>` is the repo's private project knowledge store, by the absolute path the session-start hook prints.

```
session-only ─────────────> <store>/notes/<objective>/        write now
durable fact, private ────> <store>/*.md                       rewrite in place now
durable fact, shared ─────> CONTRIBUTING.md,          ┐
                            docs/contributing/*.md    │
decision ─────────────────> the living doc it changes ├ proposed in the debrief,
                            + the why in the PR body  │ written on the user's OK
recurring friction ───────> skill / instruction fix   ┘
```

- Append the debrief's Decisions, Friction and Codify candidates bullets, dated, to the objective's note, so a later debrief can find seen-before friction and a retro (`retrospect` when installed) can count repeated candidates.
- A decision that changes how the project works edits the doc it affects; its why and the rejected alternative go in the PR body, which the squash merge keeps in git. A decision that changes nothing durable stays in the note.
- A shared fact a human contributor needs goes in `CONTRIBUTING.md` or `docs/contributing/`; only agent-only facts go in `AGENTS.md` or `docs/agents/`.
- Committed docs, skills and instructions change only after the user approves, per the write policy in "Project knowledge".

Done when the debrief is posted, every session-only and private item is written, and every shared item is either written with approval or listed as a proposal.
