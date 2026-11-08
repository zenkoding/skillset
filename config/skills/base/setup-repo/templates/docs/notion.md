# Docs

Product docs (plans, specs, milestone briefs, research) live in Notion under [<workspace page>](<url>). Engineering docs stay in the repo: `VISION.md`, `ARCHITECTURE.md`, `CONTEXT.md`, `docs/rfcs/`.

## Where each doc goes

All product docs go in the `<database>` database, one page per doc, with these properties:

| Property | Values |
| --- | --- |
| Type | Plan, Spec, Overview, Research, Notes |
| Milestone | <milestone> |
| Epic | <epic> |
| Tracker link | the ticket URL |
| Date | the day it was created |

<An index page or tag to add for a new milestone or epic, if any.>

## Creating a doc

Search the database first and update an existing page rather than create a duplicate. Create a page with every property above and link it from the ticket.

## Decisions

A decision edits the doc it changes. For repo docs, the PR body records the why and the rejected alternative, and the squash merge keeps it in git (`git log -S '<term>'`). For Notion docs, add a dated Changelog line.
