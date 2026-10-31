---
name: brief
description: Brief the user at the start of work on an idea, ticket or problem: collect context, research the unknowns, present understanding and a plan, then hand off to grilling, building or bug investigation. Run as /brief Start working on <idea | ticket | problem>.
---

# Brief

The start-of-session counterpart of `debrief`. A brief turns a one-line ask into shared understanding and a plan before any code changes, and picks the next phase. It is a short chat message plus a note `debrief` reads at the end.

## 1. Read the ask

- **Ticket** (an identifier or URL): read the issue, its comments, its parent, and every doc it links, through `docs/agents/issue-tracker.md` and `docs/agents/product-docs.md`.
- **Idea or problem** (free text): take it as given. Search the tracker for an existing issue on it and read that too.

Name the **objective** in a few words; it names the notes directory.

## 2. Collect context

Read, without asking the user for any of it:

- the project knowledge catalogs (`docs/agents/README.md` and the private store's), `CONTEXT.md`, and the ADRs in the area;
- earlier notes for the same objective under `<store>/notes/<objective>/` (a resumed objective starts from its last debrief);
- the code the ask touches: entry points, callers, tests, recent history;
- recent PRs and open issues on the same area.

Done when you can say what exists today and where the ask lands in it, with a `path:line` or link per claim.

## 3. Research the unknowns

List what you'd need to know to plan the work. Answer each from the code and primary sources inline. An unknown that could change the plan and is too big to answer inline goes to `research-scope`, in a background agent when the host has one. Stop when the remaining unknowns are ones only the user, or a grill, can settle.

## 4. Present the brief

```markdown
### Brief: <objective>

**Understanding**: <the problem, who it's for, the goal: 2-3 sentences>

**Known**
- <fact> (<path:line, doc or issue link>)

**Unknown**
- <open question>; <who or what settles it>

**Plan**
1. <step> (~<time>)

**Call**: <trivial | non-trivial | defect>: <why>

**Next**: <the next phase>. Go?
```

Make the call:

- **Defect**: the ask reports behaviour that differs from what the code or docs promise. Next is `fix-bug`, to establish the cause before any fix; the fix is then trivial or non-trivial like any other work.
- **Trivial**: bounded scope, one obvious approach, no architecture, UX/UI or taste decision. Next is building: `build-feature` for a tracker feature, otherwise test-first (the `tdd` skill when installed).
- **Non-trivial**: open-ended, needs taste decisions, or has options to choose from. Next is `grilling` (or one question at a time when it isn't installed), seeded with the Unknown list.

When in doubt between trivial and non-trivial, call it non-trivial: a short grill costs less than rework. Keep each list to about five items; the plan names times.

Wait for the user's go. A go starts the named phase; a different call from the user wins.

## 5. File the brief

Write the brief, as posted, to `<store>/notes/<objective>/brief.md` by the absolute store path the session-start hook prints (see "Project knowledge" in your user instructions). On a resumed objective, append a dated section instead of overwriting. Nothing is committed.

Done when the brief is posted, filed, and the user has said go or redirected it.
