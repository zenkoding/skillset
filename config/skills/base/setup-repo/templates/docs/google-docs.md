# Docs

Product docs (plans, specs, milestone briefs, research) live in Google Drive, in the [<folder>](<url>) folder. Engineering docs stay in the repo: `VISION.md`, `ARCHITECTURE.md`, `CONTEXT.md`, `docs/rfcs/`.

## Where each doc goes

| Doc | Folder | Name |
| --- | --- | --- |
| Plans and specs | `<folder>/Plans` | `Plan: <title>` |
| Milestone briefs | `<folder>/Milestones` | `Milestone: <name>` |
| Research | `<folder>/Research` | `Research: <title> (<YYYY-MM-DD>)` |

The first lines of each doc carry its milestone, epic and ticket link.

## Creating a doc

Agents use the Google Drive connector: search the folder first and update an existing doc rather than create a duplicate; create the doc in its folder; link it from the ticket. People use the app.

## Decisions

A decision edits the doc it changes. For repo docs, the PR body records the why and the rejected alternative, and the squash merge keeps it in git (`git log -S '<term>'`). For Drive docs, add a dated Changelog line.
