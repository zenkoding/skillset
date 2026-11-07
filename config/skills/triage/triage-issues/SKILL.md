---
name: triage-issues
description: Validate open issues against product fit, rationale, severity and work sizing, then label them, ask for missing facts, and close the ones evidence settles. Use when asked to triage issues or clean up the backlog, or from an issue triage routine.
argument-hint: "[owner/repo] [#n ...] [--dry-run]"
compatibility: Requires authenticated gh (or the tracker CLI `docs/contributing/issues.md` names) and a readable checkout of the repo.
---

# Triage issues

Turn the open backlog into work someone can pick up, a question for the reporter, or a closed issue with its reason on it. Act on what evidence settles; hand every judgment call to a person with a recommendation.

Arguments: the repo (default: the current clone), issue numbers to limit the run, and `--dry-run` to print the verdicts without changing anything. Invoking this skill authorizes the labels, comments and closes the [triage rules](references/triage-rules.md) allow, for this run.

## 1. Load the rules and the backlog

Read `docs/contributing/issues.md` for the tracker's operations and label strings, and the context the criteria need: the vision (`VISION.md` or what `CONTRIBUTING.md` points to), its non-goals, and `CONTEXT.md`.

List open issues. Triage those labelled `needs-triage`, those with no triage label at all, and those in `needs-info` (to check for an answer or for staleness). Skip issues already labelled `ready-for-agent`, `ready-for-human` or `wontfix` unless named in the arguments. Work oldest first.

Done when you hold the rules, the context and the list to triage.

## 2. Judge each issue

Read the issue with its comments. Apply the [issue criteria](references/triage-rules.md#issue-criteria) in order, and stop at the first one that settles the verdict:

1. **Product fit** against the vision, non-goals and `CONTEXT.md`.
2. **Rationale**: search for duplicates, check whether a merged change already fixed it, and for a bug check it reproduces from what the issue gives. Reproduce only when it is cheap and safe: a test, a CLI call, a read-only request.
3. **Severity**: critical, high, medium or low.
4. **Sizing**: predict the [work sizing](references/work-sizing.md) axes.

For a `needs-info` issue, first check whether the reporter answered. If they did, triage it afresh. If not, check whether it has gone stale.

Record a verdict for each issue, with the evidence for it: links to the duplicate or the fixing PR, the quoted non-goal, the reproduction you ran and its result.

Done when every issue in the list has a verdict and its evidence.

## 3. Act

With `--dry-run`, skip to step 4. Otherwise, for each issue, through the tracker doc's operations:

- **Ready for an agent / for a human**: swap `needs-triage` for the verdict's label, set the priority label when the tracker defines one, and comment with the severity, the sizing in one line, and what "done" looks like when the issue doesn't say.
- **Needs info**: label `needs-info` and comment with the one question that blocks the verdict, addressed to the reporter.
- **Close**: only for the cases under [closing without a person](references/triage-rules.md#closing-without-a-person). Comment with the reason and the evidence, then close with the matching reason (duplicate, completed, not planned).
- **Propose closing**: label `ready-for-human` and comment with why it looks closable and what would keep it open.

Never edit an issue's body, and never close one for a reason not listed in the rules. When an action fails, stop acting on that issue and report it.

Done when every verdict has its label and comment, and every close is done.

## 4. Report

Print one line per issue, grouped by verdict: number, title, verdict, severity and the reason in a few words. Ready-for-agent issues come first, critical ones on top. End with the count per verdict and the issues that need a person, so they can be read in one pass.
