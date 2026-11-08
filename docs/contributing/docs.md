# Docs

Docs live in this repo as Markdown. The repo is public, so nothing internal or private goes in them.

## Where each doc goes

Create each doc when there is something to put in it.

| Doc | Path |
| --- | --- |
| Vision | `VISION.md` |
| Architecture | `ARCHITECTURE.md`, growing into `docs/architecture/<component>.md` |
| Glossary | `CONTEXT.md` |
| RFCs | `docs/rfcs/<YYYY-MM-DD>-<slug>.md` |
| Plans | The ticket's issue body. Use `docs/plans/<YYYY-MM-DD>-<slug>.md`, linked from the ticket, when a plan is too long for an issue. |
| Milestone briefs | `docs/milestones/<YYYY-MM-DD>-<slug>.md` |
| Research | `docs/research/<YYYY-MM-DD>-<slug>.md`, only when a person asks for it to be shared; it defaults to private notes |

The date is the day the doc was created, so each directory sorts by age.

## Creating a doc

Start each doc with a title line that carries its type (`# Plan: ...`, `# RFC: ...`). Put its milestone, epic and ticket link on one line under the title. Add it through a PR; link it from the ticket.

## Decisions

A decision edits the doc it changes. The PR body records the why and the rejected alternative, and the squash merge keeps it in git: search with `git log -S '<term>'`.
