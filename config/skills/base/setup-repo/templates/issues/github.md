# Issues

Work is tracked in GitHub Issues on `<owner>/<repo>`. Use the `gh` CLI; it infers the repo inside a clone.

## Where work lives

- **Milestones**: GitHub milestones.
- **Epics**: issues labelled `epic`; their tickets are sub-issues.
- **Tickets**: issues. A ticket's plan is its issue body, or a plan doc linked from it (see [docs](docs.md)).

## Operations

- **Create**: `gh issue create --title "..." --body-file <file> [--label ...] [--milestone ...]`. Shape the body after the matching form in `.github/ISSUE_TEMPLATE/` when there is one.
- **Read**: `gh issue view <n> --comments`.
- **List**: `gh issue list --state open --json number,title,labels,assignees,milestone`.
- **Comment**: `gh issue comment <n> --body-file <file>`.
- **Label**: `gh issue edit <n> --add-label "..."` / `--remove-label "..."`.
- **Assign**: `gh issue edit <n> --add-assignee @me`.
- **Close**: `gh issue close <n> --comment "..."`.
- **Sub-issue**: `gh api --method POST repos/<owner>/<repo>/issues/<parent>/sub_issues -F sub_issue_id=<child database id>` (the id from `gh api repos/<owner>/<repo>/issues/<n> --jq .id`).
- **Blocked by**: `gh api --method POST repos/<owner>/<repo>/issues/<n>/dependencies/blocked_by -F issue_id=<blocker database id>`.

## From a branch or PR to its ticket

A PR links its ticket with `Closes #<n>` in the body. A branch carries the number when there is one (`fix/123-...`). GitHub shares one number space across issues and PRs: resolve a bare `#<n>` with `gh issue view`, then `gh pr view`.

## Plans and status

- An agent may edit a plan's **Open questions** and tick its **Done when** boxes; other sections change only with a person's OK.
- Status is open or closed, plus the labels below. A merged PR with `Closes #<n>` closes the ticket.

## Labels and triage

| Role | Label | Meaning |
| --- | --- | --- |
| needs triage | `needs-triage` | A maintainer needs to evaluate it |
| needs info | `needs-info` | Waiting on the reporter |
| ready for an agent | `ready-for-agent` | Fully specified; an agent can pick it up |
| ready for a human | `ready-for-human` | Needs a person to implement or decide |
| won't fix | `wontfix` | Will not be actioned |

New issues start as `needs-triage`. <Who triages, and how often.>
