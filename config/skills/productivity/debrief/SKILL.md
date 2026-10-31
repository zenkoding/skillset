---
name: debrief
description: Debrief the user at the end of substantial work (code changed, PR opened, decision made) or when they ask for a debrief; records durable items in the project knowledge docs.
---

# Debrief

A debrief keeps the user's understanding in step with work an agent did, and
turns repeated friction into fixes. It is a short chat message, then filing.

Run it once per objective, at the end of work that changed code, opened a PR
or made a decision, and whenever the user asks. Skip it for
question-and-answer sessions and for follow-up fixes to work already
debriefed: CI fixes, review replies, scheduled runs.

## 1. Write the debrief

Post this in chat. Each section gets 1 to 3 bullets; drop empty sections.

```markdown
### Debrief: <objective>

**Decisions**
- Chose <X> over <Y> because <Z>. (<link to ADR or doc, written or proposed>)

**Learned**
- <fact about the code, tools or domain that was not obvious>

**Friction**
- <what slowed the work: missing doc, flaky tool, unclear instruction>

**Unusual choices**
- <what a reviewer would ask "why like this?" about>

**Worth understanding**
- `path:line`: <why it matters to the user's mental model> (1-2 items)

**Proposed fixes**
- <for friction seen before: the concrete change to a skill, instruction or doc>. Awaiting your OK.
```

- Decisions always read "chose X over Y because Z", self-contained, with no
  private note paths: `publish-pr`'s reviewer notes quote them verbatim.
- Every `path:line` in Worth understanding exists at the final commit; check
  it before posting.
- A friction item is **seen before** when an earlier note under
  `<store>/notes/` records the same friction: grep the notes for its key
  words. Only seen-before friction gets a Proposed fix.

## 2. File each item

Pick the destination per item. `<store>` is the repo's private project
knowledge store, by the absolute path the session-start hook prints.

```
session-only ─────────────> <store>/notes/<objective>/        write now
durable fact, private ────> <store>/*.md                       rewrite in place now
durable fact, shared ─────> docs/agents/*.md        ┐
ADR-worthy decision ──────> docs/adr/NNNN-*.md      ├ proposed in the debrief,
recurring friction ───────> skill / instruction fix ┘ written on the user's OK
```

- Append the debrief's Decisions and Friction bullets, dated, to the
  objective's note, so a later debrief can find seen-before friction.
- A decision is **ADR-worthy** when it is hard to reverse, surprising
  without context, and a real trade-off. Otherwise it stays in the note or a
  doc.
- Committed docs, ADRs, skills and instructions change only after the user
  approves, per the write policy in "Project knowledge".

Done when the debrief is posted, every session-only and private item is
written, and every shared item is either written with approval or listed as a
proposal.
