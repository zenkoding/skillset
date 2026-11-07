# Triage rules

What `triage-issues` and `triage-prs` judge open work against, and when they may act without a person. Both skills read the same rules, so an issue and a PR about the same change get the same verdict.

## Labels

The roles come from the repo's tracker doc (`docs/contributing/issues.md`, "Labels and triage"); its label strings win over these defaults.

| Role | Default label | Means |
| --- | --- | --- |
| needs triage | `needs-triage` | Nobody has judged it yet |
| needs info | `needs-info` | Waiting on the reporter for a specific fact |
| ready for an agent | `ready-for-agent` | Fully specified and safe for an agent to pick up |
| ready for a human | `ready-for-human` | A person has to decide or do it |
| won't fix | `wontfix` | Will not be actioned |

When the tracker doc defines priority labels, set one from the severity below. Otherwise state the severity in the triage comment.

## Issue criteria

Judge each issue in this order. An earlier criterion outranks a later one: an issue that doesn't belong needs no reproduction. They are the [review criteria](review-criteria.md) applied to a request instead of a change.

1. **Product fit.** Does it belong in this project? Check it against the vision, its non-goals and `CONTEXT.md`.
2. **Rationale.** Is it a real problem or a useful feature, and not already handled?
   - **Duplicate**: an open or closed issue asks for the same outcome. Search titles and bodies for the key terms before deciding.
   - **Already fixed**: a merged PR or commit since the report changes the behaviour it describes (`git log -S '<term>'`, `gh pr list --state merged --search '<term>'`).
   - **Bugs**: reproducible from what the issue gives (steps, version, input)? Try it when that is cheap and safe; otherwise say what is missing.
3. **Severity.** How bad it is when left alone:
   - **Critical**: data loss, a security hole, or the core flow broken for everyone.
   - **High**: a core flow broken for some users, or a workaround that costs real time.
   - **Medium**: a secondary flow broken, or a clear improvement to a core one.
   - **Low**: cosmetic, or a nice-to-have.
4. **Sizing.** Predict the [work sizing](work-sizing.md) axes from the issue. Every axis low and the outcome stated clearly enough to test: an agent can take it. Any high axis, or an outcome only a person can pick: a human.

## Verdicts

| Verdict | When | Label |
| --- | --- | --- |
| **Ready for an agent** | Fits, real, every axis low, the outcome is testable | `ready-for-agent` |
| **Ready for a human** | Fits and real, but an axis is high or the outcome is a person's pick | `ready-for-human` |
| **Needs info** | A missing fact blocks the verdict, and only the reporter has it | `needs-info`, plus the exact question |
| **Close** | Evidence settles it (below) | close with `wontfix`, or as duplicate / completed |
| **Propose closing** | It looks like it should close, but that's a judgment call | `ready-for-human`, plus a comment proposing to close and why |

## Closing without a person

Close only when evidence settles it, never on judgment alone:

- a **duplicate**, linked to the issue it duplicates (close as duplicate);
- **already fixed**, linked to the PR or commit that fixed it (close as completed);
- a **stated non-goal**: the vision or `CONTEXT.md` names it as out of scope, quoted in the comment (close as not planned, `wontfix`);
- **needs info gone stale**: `needs-info` for 14 days with no answer from the reporter (close as not planned).

Anything else that looks closable, such as unclear fit, low value or a request nobody else wants, gets **Propose closing**. Every close and every label change gets a comment that says why in a sentence or two and names the evidence.
