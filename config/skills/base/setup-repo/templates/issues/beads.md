# Issues

Work is tracked in [Beads](https://github.com/gastownhall/beads), in the repo's `.beads/` directory, with the `bd` CLI. Check `bd --help` for the flags this version supports.

## Where work lives

- **Milestones**: issues of type `epic` with the label `milestone`. A milestone's epics are its children (`bd dep add <epic id> <milestone id> --type parent-child`). The milestone brief lives in the docs store (see [docs](docs.md)) and is linked from the milestone's description.
- **Epics**: issues of type `epic`; their tickets are children (`bd dep add <ticket id> <epic id> --type parent-child`).
- **Tickets**: issues of type `task`, `feature` or `bug`. A ticket's plan is its description, or a plan doc linked from it (see [docs](docs.md)).

## Operations

- **Create**: `bd create "<title>" -t <type> -p <priority> -d "<description>"`.
- **Read**: `bd show <id>`.
- **List**: `bd list --status open`; `bd ready` lists unblocked work.
- **Update status or assignee**: `bd update <id> --status in_progress --assignee <name>`.
- **Blocked by**: `bd dep add <id> <blocker id>`.
- **Close**: `bd close <id> --reason "..."`.

## From a branch or PR to its ticket

Branches and PR bodies carry the bead id (`<prefix>-<n>`).

## Plans and status

- An agent may edit a plan's **Open questions** and tick its **Done when** boxes; other sections change only with a person's OK.
- Status flow: open → in_progress → closed.

## Labels and triage

| Role | Label | Meaning |
| --- | --- | --- |
| needs triage | `needs-triage` | A maintainer needs to evaluate it |
| needs info | `needs-info` | Waiting on the reporter |
| ready for an agent | `ready-for-agent` | Fully specified; an agent can pick it up |
| ready for a human | `ready-for-human` | Needs a person to implement or decide |
| won't fix | `wontfix` | Will not be actioned |
