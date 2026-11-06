# Planning loop

brief, envision, scope-milestone and shape-epic run the same loop at different levels: a ticket's plan, the vision, a milestone brief, an epic spec. Each skill names its own output and adds its own step; this file holds what they share.

```
1 read the ask and its context
2 research the unknowns          predict the sizing
3 draft, short                   to argue with
4 grill the open decisions       settle the sizing
5 write the doc and its companion
6 sign-off                       a person approves the doc
```

## 1. Read

Read the ask and everything it points at: the issue and its comments, linked docs, earlier notes for the same objective under `<store>/notes/<objective>/`, `CONTRIBUTING.md` and the docs it points to for the area, and the decisions already made there (`git log -S '<term>'`). Separate what is already decided from what is still open: re-deciding a settled question wastes the user's time.

Name the problem in one sentence before researching. If you can't, say so and ask; everything after inherits that sentence.

## 2. Research

List what you'd need to know to write the doc, and answer each:

- **Facts from the source, not from notes about it.** Verify each load-bearing claim against the code at the current commit and record that commit in the companion; line numbers mean nothing without it.
- **Facts apart from inference.** "I read this" and "this follows from that" fail differently; keep them marked apart in your notes and in both documents.
- **The blast radius, not just the mechanism.** Who else reaches this code path, which of them are already broken, and which would break only after the change.
- **Delegate breadth, keep depth.** A sweep across many files whose answer is a conclusion goes to `research`, in a background agent when the host has one. A lookup in code you already hold doesn't. Never delegate the one fact the design turns on without checking it yourself.

Then **predict** the axes of [work sizing](work-sizing.md). The prediction picks the route the skill gives (brief: trivial, defect or non-trivial).

Done when every unknown is answered with a source or carried into the grill as a question.

## 3. Draft

Post a draft in chat before writing anything final, short enough to argue with: its job is to surface the decisions, not to look finished.

## 4. Grill

Run `grill` on the draft, seeded with what research left open and each high axis. The user's answers can reshape the draft: recompute what is still open after each round rather than carrying the old list.

Then **settle** the sizing: answers can raise or lower an axis, and a non-trivial call can turn trivial. The settled axes, with a why for each high one, go in the doc's risks; the companion keeps the call.

## 5. Write the doc and its companion

Write the skill's doc where `docs/contributing/docs.md` says docs live, and link it from the tracker issue. Write the **companion** to `<store>/notes/<objective>/companion.md`: code pointers at the recorded commit, per-step instructions, tests, the sizing call and the decision log (the plan template's companion section lists them). The doc never links the companion.

Read the doc's steps back as a sequence and check each leaves the system working, including against changes later steps make.

Style, for every doc: plain declaratives for a competent reader new to the area; roles, not names; no file paths; unverified claims marked where they are made; a diagram only when the structure is the point.

## 6. Sign-off

Show the user the doc and ask for sign-off. Expect wording rounds: treat a correction as a pattern and apply it everywhere it fits, in both documents, before coming back. When the doc and companion disagree, the doc wins; edit both in the same pass.

Done when the user has approved the doc. Nothing is built before that.
