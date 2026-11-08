# Issues

Work is tracked in [Beads](https://github.com/gastownhall/beads), in the repo's `.beads/` directory.

## Where work lives

- **Milestones**: issues of type `epic` with the label `milestone`. A milestone's epics are its children. The milestone brief lives in the docs store (see [docs](docs.md)) and is linked from the milestone's description.
- **Epics**: issues of type `epic`; their tickets are children.
- **Tickets**: issues of type `task`, `feature` or `bug`. A ticket's plan is its description, or a plan doc linked from it (see [docs](docs.md)).

## Filing a ticket

File and update tickets with the `bd` CLI; `bd --help` lists its commands.

## From a branch or PR to its ticket

Branches and PR bodies carry the bead id (`<prefix>-<n>`).

## Status

Status flow: open → in_progress → closed.

## Labels and triage

| Role | Label | Meaning |
| --- | --- | --- |
| needs triage | `needs-triage` | A maintainer needs to evaluate it |
| needs info | `needs-info` | Waiting on the reporter |
| ready for an agent | `ready-for-agent` | Fully specified; an agent can pick it up |
| ready for a human | `ready-for-human` | Needs a person to implement or decide |
| won't fix | `wontfix` | Will not be actioned |
