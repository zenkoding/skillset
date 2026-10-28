---
name: optimize-performance
description: Find and land measured performance improvements in a code area using the project's benchmarks and profiles. Use when asked to make something faster or cheaper and prove it; not for diagnosing a single slow-request bug.
compatibility: Requires a writable checkout, the project's build and test tools, and a runnable benchmark or a way to write one. Parallel subagents are optional.
---

# Optimize performance

Land changes that make the target measurably faster without breaking behaviour. Every claimed gain comes with a before/after measurement; a change that does not beat noise is reverted.

## 1. Agree the target

Establish the code area, the metric (latency, throughput, allocations, memory), and the budget: how many hotspots to attempt and the smallest gain worth keeping (default: up to 5 hotspots, at least 5% improvement). Ask the user only when the area or metric is genuinely unclear.

## 2. Find the benchmark

Locate the project's existing benchmarks and how they run (build files, test conventions, CI config). If none covers the target, write a focused one and commit it separately first. Run it enough times to see the noise; record the **baseline** with the exact command, machine, and sample count.

Done when the baseline is reproducible within a known noise band.

## 3. Profile

Profile the benchmark with the language's standard profiler (CPU, then allocations or memory as the metric demands). Rank **hotspots** by share of cost. Base the list on profile data, not on reading code.

## 4. Optimize one hotspot at a time

For each hotspot in rank order:

1. Form a hypothesis for why it is slow and the smallest change that tests it.
2. Make the change; run the tests.
3. Re-run the benchmark with the same command and sample count as the baseline.
4. Keep the change only if tests pass and the gain clears both the noise band and the agreed threshold. Commit it on its own with the measured numbers in the message. Otherwise revert it and record why.

Give each hotspot at most three attempts. When the host supports subagents, independent hotspots in separate files may be investigated in parallel, each in its own worktree; merge and re-measure them one at a time. Otherwise work through them in order.

Done when every hotspot in the budget is either landed with numbers or recorded as not worth it.

## 5. Report

Save a report to `docs/performance/<date>-<area>.md` unless the user names another place: baseline, each attempt with before/after numbers and verdict, the final cumulative gain, and remaining hotspots worth a follow-up. Open a GitHub issue with `gh issue create` for follow-ups only when the user asks.

Return the cumulative gain, the commits, and the report path.
