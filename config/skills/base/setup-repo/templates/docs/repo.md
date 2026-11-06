# Docs

Docs live in this repo as Markdown.

## Where each doc goes

| Doc | Path |
| --- | --- |
| Vision | `VISION.md` |
| Architecture | `ARCHITECTURE.md`, growing into `docs/architecture/<component>.md` |
| Glossary | `CONTEXT.md` |
| RFCs | `docs/rfcs/<NNNN>-<slug>.md` |
| Plans, specs, milestone briefs | `docs/plans/<YYYY-MM-DD>-<slug>.md` |
| Research | `docs/research/<YYYY-MM-DD>-<slug>.md`, only when a person asks for it to be shared; it defaults to private notes |

## Creating a doc

Start each doc with a title line that carries its type (`# Plan: ...`, `# Spec: ...`). Put its milestone, epic and ticket link on one line under the title. Add it through a PR; link it from the ticket.

## Decisions

A decision edits the doc it changes. The PR body records the why and the rejected alternative, and the squash merge keeps it in git: search with `git log -S '<term>'`. <An existing `docs/adr/` is frozen history: read it, don't add to it.>
