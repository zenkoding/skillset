# Docs

Product docs (plans, specs, milestone briefs, research) live in Google Drive, in the [<folder>](<url>) folder. Engineering docs stay in the repo: `VISION.md`, `ARCHITECTURE.md`, `CONTEXT.md`, `docs/rfcs/`.

## Where each doc goes

| Doc | Folder | Name |
| --- | --- | --- |
| Plans | `<folder>/Plans` | `<YYYY-MM-DD>-<title>` |
| Specs | `<folder>/Specs` | `<YYYY-MM-DD>-<title>` |
| Milestone briefs | `<folder>/Milestones` | `<YYYY-MM-DD>-<milestone>` |
| Research | `<folder>/Research` | `<YYYY-MM-DD>-<title>` |

The date is the day the doc was created, so each folder sorts by age. The doc's own title line carries its type (`Plan: <title>`).

The first lines of each doc carry its milestone, epic and ticket link.

## Creating a doc

Agents use the Google Drive connector: search the folder first and update an existing doc rather than create a duplicate; create the doc in its folder; link it from the ticket. People use the app.

## Decisions

A decision edits the doc it changes. For repo docs, the PR body records the why and the rejected alternative, and the squash merge keeps it in git (`git log -S '<term>'`). For Drive docs, add a dated Changelog line.
