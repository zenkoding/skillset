---
name: codify
description: Turn a recurring mechanical mistake into a deterministic check (lint rule, type, test, CI step or codemod) with fixtures, wired into the repo's checks and opened as one PR. Run by retrospect for a repeated Codify candidate, or when the user asks to codify a rule.
compatibility: Requires Git, the repo's own check tooling, and authenticated gh to open the PR.
---

# Codify

Learn a mistake once. A check that runs on every change beats a review that might catch it next time, costs no tokens, and treats every author the same. The PR is the deliverable.

Arguments: the candidate as a `Codify: <rule>` line, a link to where it was found (review comment, debrief, investigation report), or the rule in plain words. With none, ask for it.

## 1. Confirm the candidate

Restate the rule in one line and collect every instance you were given or can find: the review comments, debriefs and reports that tagged it, and the code they point at. Check it against **codifiable** in [work sizing](references/work-sizing.md). When the rule needs judgment, stop: propose it as a `REVIEW.md` rule instead (`review` reads it) and say why.

Done when you hold the rule, at least one real instance, and a yes on codifiable.

## 2. Pick the cheapest home

Read the repo's checks first: linter configs (`.golangci.yml`, `eslint.config.*`, `ruff.toml`, `.semgrep*`), CI workflows, the `Makefile` or task runner, pre-commit hooks, and the check command `AGENTS.md` or `CLAUDE.md` names. Take the first option that can express the rule:

1. **Enable or configure** a rule the repo's linter already ships (e.g. `forbidigo`, `depguard`, `gocritic` in golangci-lint; `no-restricted-imports` in ESLint).
2. **Make it impossible**: a type, an unexported constructor or a narrower API, so the wrong code no longer compiles.
3. **A custom rule** in a tool the repo already runs: `ruleguard` or a `go/analysis` analyzer for Go, a semgrep rule, a custom ESLint rule.
4. **A test** that walks the code or config (a table-driven Go test over the tree is fine).
5. **A CI step or script**, only when nothing above fits.

Never add a new tool when an existing one can express the rule. A check that exists but isn't wired into the check command is the fix: wire it.

Done when you can name the home and why each cheaper option above it can't express the rule.

## 3. Write fixtures first

Before the rule, write a **bad** example the check must flag (the real instance, trimmed) and a **good** example it must pass, plus the nearest legitimate code that looks similar. Put them where the tool expects test data (`testdata/` for analyzers and ruleguard, the tool's own test layout otherwise).

Done when the fixtures exist and the empty or missing rule fails on the bad one.

## 4. Write the rule

Write the rule until it flags the bad fixture and passes the good and near-miss ones. Its failure message says what is wrong and what to write instead, in one line, and names the doc behind the rule when there is one. When the tool supports autofixes, add one for mechanical rewrites.

Done when the tool's own test run passes on the fixtures.

## 5. Run it on the whole repo

Run the check over the repo and read every hit:

- **False positive**: narrow the rule and add the case as a near-miss fixture. When it can't be narrowed, stop: the rule is not codifiable; fall back to a `REVIEW.md` rule.
- **Real violation**: fix it in this PR, with the autofix or a codemod when there are many. When the fixes would swamp the review (a large diff by [work sizing](references/work-sizing.md)), land the check with the tool's baseline or per-path excludes, list them under **Reviewer notes**, and propose the cleanup as a follow-up issue.

Then run the repo's full check command.

Done when the check runs from the repo's check command, exits 0, and every remaining exclude is listed.

## 6. Open the PR

Commit and open one PR with `publish`:

- **What**: the rule, the bad fixture and the message it prints.
- **Why**: the rule in one line and every instance that led here, linked (review comments, debriefs, investigation reports).
- **Reviewer notes**: the `Guard` line naming the new check, plus any baseline or excludes.

Done when the PR is open. When a retro called this skill (`retrospect` when installed), done when the change is on the retro's branch for its one PR.
