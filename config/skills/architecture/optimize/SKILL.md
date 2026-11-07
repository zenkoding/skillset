---
name: optimize
description: Take one performance issue (usually from find-hotspots) or a named code area and land a measured improvement, with before/after numbers for every claimed gain; trivial changes go to a PR, larger ones to brief. Use when asked to make something faster or cheaper and prove it; not for diagnosing a single slow-request bug.
argument-hint: "<issue | area> [--threshold N%]"
compatibility: Requires a writable checkout, the project's build and test tools, and a runnable benchmark or a way to write one. Parallel subagents are optional.
---

# Optimize

Make one hotspot measurably faster or cheaper without breaking behaviour. Every claimed gain comes with a before/after measurement; a change that does not beat noise is reverted.

Arguments: the issue, or a code area when there's no issue, and the smallest gain worth keeping (default 5%). Invoking `optimize` on an issue authorizes the branch and PR for a trivial change.

## 1. Agree the target

From the issue, take the metric, the baseline and the code it points at; re-check the baseline still holds at the current commit. Without an issue, agree the area and metric with the user. Read `docs/contributing/performance.md` for the telemetry store, budgets and benchmark commands.

Done when the target, the metric and the gain worth keeping are each one line.

## 2. Find or write the benchmark

Locate a benchmark that exercises the hotspot. If none does, write a focused one and commit it on its own first. Run it enough times to see the noise, and record the **baseline** with the exact command, machine and sample count.

Done when the baseline is reproducible within a known noise band.

## 3. Profile

Profile the benchmark with the language's standard profiler (CPU, then allocations or memory as the metric demands). Confirm the cost sits where the issue says; when it doesn't, comment on the issue with what you measured and stop.

## 4. Size the change

Form a hypothesis for why it is slow and the smallest change that tests it. Size that change with [work sizing](references/work-sizing.md).

- **Trivial** (every axis low, two-way door): go on.
- **Not trivial**: a new cache, a changed data layout, a concurrency model, a public API change. Hand the issue to `brief` with the baseline, the profile and the hypothesis, and stop.

## 5. Change, measure, keep or revert

1. Make the change; run the tests.
2. Re-run the benchmark with the same command and sample count as the baseline.
3. Keep it only if the tests pass and the gain clears both the noise band and the threshold. Commit it with the numbers in the message. Otherwise revert it and record why.

Give the hypothesis at most three attempts.

Done when the change is kept with numbers, or recorded as not worth it.

## 6. Publish

Open the PR with `publish`, linking the issue. How to test gives the benchmark command; the PR body carries baseline, after, and the gain with its noise band. When nothing was kept, comment on the issue with the attempts and their numbers instead.
