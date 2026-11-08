# Issue tracker operations

How agents work with GitHub Issues on `zenkoding/skillset`. Use the `gh` CLI; it infers the repo inside a clone. What the labels and work items mean is in `docs/contributing/issues.md`.

- **Create**: `gh issue create --title "..." --body-file <file> [--label ...] [--milestone ...]`. Shape the body after the matching form in `.github/ISSUE_TEMPLATE/`, with the same headings and the same labels.
- **Read**: `gh issue view <n> --comments`.
- **List**: `gh issue list --state open --json number,title,labels,assignees,milestone`.
- **Comment**: `gh issue comment <n> --body-file <file>`.
- **Label**: `gh issue edit <n> --add-label "..."` / `--remove-label "..."`.
- **Assign**: `gh issue edit <n> --add-assignee @me`.
- **Close**: `gh issue close <n> --comment "..."`.
- **Sub-issue**: `gh api --method POST repos/zenkoding/skillset/issues/<parent>/sub_issues -F sub_issue_id=<child database id>` (the id from `gh api repos/zenkoding/skillset/issues/<n> --jq .id`).
- **Blocked by**: `gh api --method POST repos/zenkoding/skillset/issues/<n>/dependencies/blocked_by -F issue_id=<blocker database id>`.

GitHub shares one number space across issues and PRs: resolve a bare `#<n>` with `gh issue view`, then `gh pr view`.

You may edit a plan's **Open questions** and tick its **Done when** boxes. Other sections change only with a person's OK.
