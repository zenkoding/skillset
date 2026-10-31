---
name: test-e2e
description: Run one PR's feature-plan "How to test" steps against a fresh build and post a verdict per step on the PR. Run as /test-e2e <PR number or URL> [mode].
disable-model-invocation: true
---

# Test E2E

Run a feature's **How to test** steps against a fresh build of the PR and post a **verdict** per step on the PR. You are the reviewer's eyes in the running app: a step passes only when you saw its expected result yourself. The PR comment is the deliverable.

Arguments: the PR, then the stack mode the repo's e2e doc defines (its default when omitted). The feature plan's format and the tracker operations are in `docs/agents/issue-tracker.md`.

## 1. Find the feature

```bash
gh pr view <pr> --json number,title,body,headRefName,headRefOid,isDraft,state,url
```

1. Stop without posting when the PR is a draft or closed.
2. Find the tracker identifier in the head branch name, then the PR title, then the body's **Why**. None found: post the **no plan** comment (step 4) and stop.
3. Read the feature from the tracker. Keep **How to test**, **Flags**, its status and its parent epic. With no plan in the tracker, fall back to the PR body's **How to test**.

Done when you hold the PR number, the head SHA and the steps. Empty or missing steps get the **no plan** comment, and you stop.

## 2. Bring up the stack

Check out the head SHA on a clean tree (`gh pr checkout <pr>`). Then follow `docs/agents/e2e.md` for the chosen mode ([format](references/e2e-doc.md)): build, start, wait for health, set flags.

Without `docs/agents/e2e.md`, infer the same steps from `README.md`, `CONTRIBUTING.md`, `AGENTS.md` and `CLAUDE.md`, and say in the report which commands you inferred and from where. When the docs don't settle a step, the run is **blocked** on it: never guess a command that writes outside the checkout.

Set every flag named in **Flags** to the stage the plan gives. A flag you can't set makes the steps that need it **blocked**, with the reason.

Done when the app answers its health check. A stack that doesn't come up is a **blocked** run: post it with the last 30 lines of the start output, stop whatever started, and end the run.

## 3. Run the steps

Run each step in order, exactly as written:

- **UI steps**: drive a browser at the app's URL with whichever browser tool the session has (Playwright, or the desktop app's browser). Read the page as text to check results; screenshot each failing step. With no browser tool, UI steps are **blocked** ("no browser in this environment").
- **Commands and API calls**: run them; keep the exit code and the relevant output lines.

Give each step one verdict:

| Verdict | When |
| --- | --- |
| pass | You saw the step's expected result |
| fail | You saw a different result. Record expected vs. actual in one line each |
| blocked | The step could not run: an earlier failure it depends on, an unset flag, missing data, no browser, an unclear instruction |

Run every step even after a fail. Treat an ambiguous step as **blocked** with the question it raises.

Done when every step has a verdict and every fail or blocked has its evidence.

## 4. Post the report

The report is one comment per PR. Find an earlier one: `gh pr view <pr> --json comments` and look for a body starting `## E2E:`. If there is one, edit it (`gh api repos/{owner}/{repo}/issues/comments/<id> --method PATCH -F body=@<file>`); else `gh pr comment <pr> --body-file <file>`. Write the body file in your own scratch or worktree directory.

```markdown
## E2E: <pass | fail | blocked> · <feature id> @ <short SHA>

Stack: <mode> · Flags: <as set, or none> · Setup: <docs/agents/e2e.md, or "inferred from <files>">

| # | Step | Verdict | Evidence |
| - | ---- | ------- | -------- |
| 1 | <step, shortened> | pass | <what you saw> |
| 2 | ... | fail | Expected ... / got ... |

<For each fail: steps to reproduce and the screenshot or output.>
```

The headline is **pass** only when every step passed, **fail** when any step failed, else **blocked**. The **no plan** comment is the same headline with `blocked` and one line: which plan section is missing, or that no identifier was found.

For each failing step, record the bug as the tracker adapter says: on a feature not yet Done, a checklist item under its **Open questions** (change no other section; skip an item already there for the same step); on a Done feature, a new bug issue in the same epic, unless one for this step is already open.

Done when the comment is on the PR. Stop the stack as the e2e doc says.
