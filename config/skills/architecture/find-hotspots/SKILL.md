---
name: find-hotspots
description: Find where a project spends the most time, memory or money, from its telemetry store (or benchmarks and profiles when it has none), and file each measured hotspot as an issue with baseline numbers. Use when asked what to optimize, to find performance problems, or from a performance routine.
argument-hint: "[area or metric] [--budget N]"
compatibility: Requires a readable checkout; access to the telemetry store docs/contributing/performance.md names, or the project's benchmark and profiling tools; authenticated gh or the tracker CLI for filing.
---

# Find hotspots

Find the few places where a fix would move a metric people care about, measure each, and file it so `optimize` can take it. A hotspot is a measured share of cost, never a guess from reading code.

Arguments: an area or metric to focus on (default: every metric in the performance doc) and a budget (default 5, per the [find rules](references/find-rules.md#budget)).

## 1. Load where performance is measured

Read `docs/contributing/performance.md`: the telemetry store, the metrics with their queries and budgets, profiles and benchmarks. When it is missing, say so, suggest `/setup-repo performance`, and go on with the benchmarks and profilers the build files and CI already use.

Done when you know, per metric, where to read it and what its budget is.

## 2. Measure

- **From telemetry**: run each metric's query over a window long enough to cover normal load (default 7 days). Break it down by the dimension that points at code: route, operation, query, job. Note what is over budget, what got worse (compare with the window before), and what dominates total cost.
- **From profiles**: for the top entries, pull CPU and allocation profiles from the same window and find the functions that carry the cost.
- **Without telemetry**: run the benchmarks enough times to see noise, then profile the slowest.

Record every number with its source, window or command, so the issue can be checked later.

Done when you have a ranked list of candidates, each tied to code and carrying numbers.

## 3. Weigh each hotspot

Keep a candidate when it is over budget, regressed, or carries a large share of total cost, and when the cost sits in this repo's code (not a dependency you can't change, or the network). Rate severity: **high** when it breaks a budget users feel, **medium** when it regressed or costs real money, **low** for headroom.

Done when every kept hotspot has a location, numbers and a severity.

## 4. File and report

File each hotspot per the [find rules](references/find-rules.md#filing), with the lane label `performance`. Evidence holds the **baseline**: the metric, its value, the query or benchmark command, the window or machine, and the profile excerpt that points at the code. Then report per the find rules.
