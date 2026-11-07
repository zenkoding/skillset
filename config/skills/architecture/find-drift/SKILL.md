---
name: find-drift
description: Compare the code against the architecture the repo documents (ARCHITECTURE.md, docs/architecture/, CONTEXT.md, recorded decisions) and file each evidenced deviation as an issue. Use when asked to check architecture, find drift or audit module boundaries, or from an architecture routine.
argument-hint: "[path or area] [--budget N]"
compatibility: Requires a readable checkout; authenticated gh or the tracker CLI for filing. Parallel subagents are optional.
---

# Find drift

Find where the code has moved away from the architecture the repo says it has, and file each deviation so `refactor` can bring it back. The documented architecture is the standard; your own taste is not.

Arguments: the area to check (default: the whole repo) and a budget (default 5, per the [find rules](references/find-rules.md#budget)).

## 1. Load the agreed architecture

Read `ARCHITECTURE.md` (or `docs/architecture/`), `CONTEXT.md` for names, `CONTRIBUTING.md` and the docs it points to, and the decisions recorded in merged PR bodies for the area (`git log -S '<term>'`). Turn them into a list of checkable rules, each with its source: layers and which way dependencies point, which component owns which data, where an entry point lives, the names the domain uses, the patterns the repo chose over alternatives.

When the repo documents no architecture, stop: there is nothing to drift from. Suggest writing `ARCHITECTURE.md` (`setup-repo` covers it), and report the de facto structure you saw in a few lines.

Done when every rule you'll check has a quoted source.

## 2. Check the code against each rule

- **Dependencies**: imports that point the wrong way between layers or reach into another component's internals.
- **Ownership**: data written or business rules decided outside the component that owns them.
- **Placement**: code in the wrong layer, entry points outside where the docs say they live.
- **Names**: types, packages or functions that use a different word than `CONTEXT.md` for the same thing.
- **Patterns**: a second way of doing something the repo already decided how to do.

Use the language's tools where they help (import graphs, `go list -deps`, linters). For a large repo, when the host supports subagents, give each one a slice of the rules.

Done when every rule has been checked across the area.

## 3. Weigh each deviation

Keep a deviation only when it is located and its rule is quoted. Drop it when a later decision changed the rule; say in the report that the doc is stale instead. Rate severity by what the drift costs: **high** when it couples components the architecture keeps apart or puts data ownership in doubt, **medium** when it spreads a second pattern, **low** for naming and placement.

Done when every kept deviation has a location, its rule and a severity.

## 4. File and report

File each deviation per the [find rules](references/find-rules.md#filing), with the lane label `architecture`, the broken rule quoted in Evidence, and the smallest move that restores it in Suggested fix. Then report per the find rules, plus the stale docs you found.
