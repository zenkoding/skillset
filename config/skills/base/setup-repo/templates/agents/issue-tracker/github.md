# Issue tracker: GitHub

Read [Issues](../contributing/issues.md) first: it says where work lives, how branches and PRs point at tickets, the status flow and the labels. This doc adds what only an agent needs.

The repo is `<owner>/<repo>`. Use the `gh` CLI; it infers the repo inside a clone.

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

## Resolving a bare number

GitHub shares one number space across issues and PRs: resolve a bare `#<n>` with `gh issue view`, then `gh pr view`.

## Editing a plan

You may edit a plan's **Open questions** and tick its **Done when** boxes; other sections change only with a person's OK.
