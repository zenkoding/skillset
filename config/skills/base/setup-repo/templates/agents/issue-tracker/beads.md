# Issue tracker: Beads

Read [Issues](../contributing/issues.md) first: it says where work lives, how branches and PRs point at tickets, the status flow and the labels. This doc adds what only an agent needs.

Use the `bd` CLI. Check `bd --help` for the flags this version supports.

## Operations

- **Create**: `bd create "<title>" -t <type> -p <priority> -d "<description>"`.
- **Read**: `bd show <id>`.
- **List**: `bd list --status open`; `bd ready` lists unblocked work.
- **Update status or assignee**: `bd update <id> --status in_progress --assignee <name>`.
- **Add an epic to a milestone**: `bd dep add <epic id> <milestone id> --type parent-child`.
- **Add a ticket to an epic**: `bd dep add <ticket id> <epic id> --type parent-child`.
- **Blocked by**: `bd dep add <id> <blocker id>`.
- **Close**: `bd close <id> --reason "..."`.

## Editing a plan

You may edit a plan's **Open questions** and tick its **Done when** boxes; other sections change only with a person's OK.
