# Issues

Work is tracked in Linear: team `<team>`, project `<project>`. Agents use the Linear connector's tools; people use the app.

## Where work lives

- **Milestones**: project milestones. The milestone brief lives in the docs store (see [docs](docs.md)) and is linked from the milestone's description.
- **Epics**: parent issues in the project; their tickets are sub-issues.
- **Tickets**: issues. A ticket's plan is its description, or a plan doc linked from it (see [docs](docs.md)).

## Operations

- **Create**: create an issue in team `<team>`, project `<project>`, with its milestone and parent. Check for a duplicate first by listing issues with a keyword query.
- **Read**: get the issue, then list its comments (the description doesn't include them).
- **List**: list issues filtered by project, status, label or assignee; follow the page cursor.
- **Comment**, **label** (add or remove one label without replacing the set), **assign**, **set status**, **set blocked by**: the matching issue update.

## From a branch or PR to its ticket

Branches carry the identifier in lower case (`<team key>-123-...`), which Linear's GitHub integration uses to link the PR and move the issue. A PR body names the identifier in its **Why**.

## Plans and status

- An agent may edit a plan's **Open questions** and tick its **Done when** boxes; other sections change only with a person's OK.
- Status flow: <Backlog → Todo → In Progress → In Review → Done>. The integration moves a ticket to In Progress when a branch with its identifier exists and to In Review when the PR opens; check it moved and set it by hand only if it didn't.

## Labels and triage

| Role | Label | Meaning |
| --- | --- | --- |
| needs triage | `needs-triage` | A maintainer needs to evaluate it |
| needs info | `needs-info` | Waiting on the reporter |
| ready for an agent | `ready-for-agent` | Fully specified; an agent can pick it up |
| ready for a human | `ready-for-human` | Needs a person to implement or decide |
| won't fix | `wontfix` | Will not be actioned |
