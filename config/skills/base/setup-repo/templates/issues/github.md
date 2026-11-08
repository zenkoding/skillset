# Issues

Work is tracked in GitHub Issues on `<owner>/<repo>`.

## Where work lives

- **Milestones**: GitHub milestones. The milestone brief lives in the docs store (see [docs](docs.md)) and is linked from the milestone's description.
- **Epics**: issues labelled `epic`; their tickets are sub-issues.
- **Tickets**: issues. A ticket's plan is its issue body, or a plan doc linked from it (see [docs](docs.md)).

## Filing a ticket

Open a new issue on GitHub and pick the matching form when the repo has one.

## From a branch or PR to its ticket

A PR links its ticket with `Closes #<n>` in the body. A branch carries the number when there is one (`fix/123-...`).

## Status

Status is open or closed, plus the labels below. A merged PR with `Closes #<n>` closes the ticket.

## Labels and triage

| Role | Label | Meaning |
| --- | --- | --- |
| needs triage | `needs-triage` | A maintainer needs to evaluate it |
| needs info | `needs-info` | Waiting on the reporter |
| ready for an agent | `ready-for-agent` | Fully specified; an agent can pick it up |
| ready for a human | `ready-for-human` | Needs a person to implement or decide |
| won't fix | `wontfix` | Will not be actioned |

New issues start as `needs-triage`. <Who triages, and how often.>
