# Issues

Work is tracked in [GitHub Issues](https://github.com/zenkoding/skillset/issues) on `zenkoding/skillset`.

## Filing an issue

Open a [new issue](https://github.com/zenkoding/skillset/issues/new/choose) and pick the bug report or feature proposal form. Blank issues are off.

Propose new work through a form rather than opening a pull request first. Triage covers issues only, and external pull requests are not treated as feature requests.

## How work is organized

- **Milestones** are GitHub milestones. The milestone brief lives in the repo (see [docs](docs.md)) and is linked from the milestone's description.
- **Epics** are issues labelled `epic`, with their tickets as sub-issues.
- **Tickets** are issues. A ticket's plan is in its issue body, or in a plan doc linked from it.

## Picking up work

Assign the issue to yourself before you start. Include the issue number in the branch name when there is one (`fix/123-short-name`). Put `Closes #<number>` in the PR body, so that merging the PR closes the issue.

## Labels and triage

| Label | Meaning |
| --- | --- |
| `needs-triage` | A maintainer needs to evaluate it |
| `needs-info` | Waiting on the reporter |
| `ready-for-agent` | Fully specified; an agent can pick it up |
| `ready-for-human` | Needs a person to implement or decide |
| `wontfix` | Will not be actioned |

New issues start as `needs-triage`, and the forms also set `bug` or `enhancement`. The maintainer triages issues as they arrive.
