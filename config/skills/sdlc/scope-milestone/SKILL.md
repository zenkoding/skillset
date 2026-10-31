---
name: scope-milestone
description: Scope the next milestone. Compare where the project is with the vision, grill the goal, write the milestone brief, create the tracker milestone and seed its epics in Backlog. Run as /scope-milestone.
disable-model-invocation: true
---

A product-lane skill. It turns "what should we do next" into one milestone: a brief in the docs store, a milestone in the tracker with its "done when", and epics waiting in Backlog for `shape-epic`. Use the glossary's terms (milestone, epic, feature). Every tracker and docs-store operation below runs as the adapters say: `docs/agents/issue-tracker.md` and `docs/agents/product-docs.md`.

The person decides the goal, the non-goals, the "done when" and which epics exist. You find the facts and draft.

## 1. Gather where we are

When the previous milestone is closed and no `retrospect` ran on it, offer one first (when installed): its product lessons are input for this milestone.

Read, without asking the person for any of it:

- The vision doc, the project-wide decisions doc when the repo has one, and the glossary.
- The project's milestones and epics in the tracker: what is Done, In Progress and in Backlog.
- The docs store: briefs of earlier milestones and research tagged for the milestone being scoped.
- What runs on the default branch today: recent merged PRs and the project home, enough to say what the product does.

Done when you can write, with a link or PR for each claim, what the product does today and which parts of the vision it does not yet reach. That is the **gap**.

## 2. Grill the goal

Run a grilling skill when installed (otherwise interview the person one question at a time) on the gap. The design tree starts at the goal and branches into every section of [`milestone-brief.md`](milestone-brief.md): owner and target date, goal, rough shape, non-goals, done when with appetite, and epics. The template holds the shape rules for each; hold every answer to them.

Two rules the template cannot carry:

- **Non-goals** inherit every item in the vision's "Not in v1" and add the milestone's own. If the person wants to drop a vision non-goal, that is a vision change: it goes through `envision`, not this brief.
- **Epics** cover the "done when": every check is covered by at least one epic, and an epic that covers no check is cut or its check is missing.

Done when the frontier is empty and the person confirms the shared understanding.

## 3. Write the brief

List the docs store's overview docs for this milestone first; if a brief exists, update it instead of creating another. Otherwise fill the template and create the brief as `docs/agents/product-docs.md` says: type Overview, the new milestone, today's date, tracker link empty until step 4.

Show the person the brief and wait for approval. Edits go into the doc before anything is created in the tracker.

## 4. Create the milestone and seed the epics

1. List the milestones, then the milestone's epics, in the tracker. Create the milestone only if it is missing, and only the epics from the brief's table that are missing. A re-run then completes a half-finished run instead of duplicating it.
2. Create the milestone: named by its stage (`PoC`, `MVP`), description the "done when" sentence and the brief link.
3. Create each missing epic in Backlog in that milestone: description the outcome in one line and the brief link. Add blocks between epics where the table gives an order.
4. Do whatever else `docs/agents/product-docs.md` asks for a new milestone or epic (an index page, a tag option), so `shape-epic` can file its docs under it.
5. Fill the epic links in the brief's table, and set the brief's tracker link to the milestone.

If a call fails, stop and report what exists and what does not, by name and URL.

Done when each of these links resolves:

- the brief's epics table → each epic
- each epic's description → the brief
- the milestone's description → the brief
- the brief's tracker link → the milestone

Leave the epics in Backlog: a person moves an epic out once its spec is approved.

## 5. Hand over

Report the brief URL, the milestone URL and the epic identifiers. Then tell the person:

- the vision's goals should link this milestone; that edit goes through `envision`.
- the next command is `/shape-epic <epic id>`, starting with the first epic in the table that nothing blocks.
