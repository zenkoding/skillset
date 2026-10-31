---
name: shape-epic
description: Shape one Backlog epic into an approval-ready spec and a set of vertical-slice features in the tracker. Run as /shape-epic <epic id>.
disable-model-invocation: true
---

Take one epic from a one-line outcome to a spec a person can approve and a set of features an agent can build. Read `docs/agents/issue-tracker.md` and `docs/agents/product-docs.md` first: they say where epics, features and docs live and how to write to them. Use glossary terms throughout.

The argument is an epic identifier. If none is given, ask for one.

## 1. Load the epic

Read the epic issue, its comments, its milestone, and every doc it links. Read the milestone brief and note the **Covers check** row or rows that name this epic: the spec must cover them. Read the vision doc, the project-wide decisions doc when there is one, and the ADRs in the area the epic touches. List the epic's sub-issues and its docs in the docs store.

Stop and ask the person when:

- the epic is not in Backlog: its spec is approved, and changing it is the person's call;
- a spec doc or features already exist: ask whether to update them in place. Update, never create a second spec or a duplicate feature.

Done when you can state the epic's outcome in one line and name every source you read.

## 2. Research the unknowns

List what you would need to know to write **Done when** and cut features: how the code works today, what an external API allows, what prior art exists. Answer each from the code and primary sources. A question too big to answer inline goes to `research-scope`, in a background agent when the host has them; its findings go to the docs store as a Research doc for this epic, not the repo.

Done when every unknown is answered with a source or carried into the grill as a question.

## 3. Grill the person

Run a grilling skill when installed (otherwise interview the person one question at a time), seeded with what step 2 left open. Cover, at least: the problem story, what is out of scope, how each outcome is checked, appetite, flags and rollout, and anything that would force a feature across epics.

Done when every branch of the design tree has an answer or is recorded as an open question in the spec.

## 4. Record decisions

Each settled call that later work must respect is a **decision**. Give each one home:

- How the system is built: an ADR in `docs/adr/`.
- A change to a contract other code or people depend on (an API, a schema, an artifact format): an RFC, written with `create-rfc`. Otherwise an ADR is enough.
- A project-wide rule: an entry in the decisions doc. A person must approve that edit in review.
- A product call only this epic needs: the spec's **Rabbit holes**.

ADRs, RFCs and decisions-doc edits land through a branch and a PR opened with `publish-pr`. The spec links each one.

Done when every decision from the grill has one home and each repo change is in an open PR.

## 5. Write the spec

Fill [`epic-spec.md`](epic-spec.md). Its **Done when** covers every milestone **Covers check** row for this epic. Create it in the docs store as a Spec doc with the epic's milestone, epic and tracker link, as `docs/agents/product-docs.md` says. Link the spec from the epic issue.

Done when the spec fits on one page (hard cap two), every **Done when** row names its check, and **Non-goals** gives a reason for each item.

## 6. Cut features

Split the epic into **vertical slices**: each feature takes a narrow path through every layer it needs (schema, API, UI, tests) and is demoable or checkable on its own. Size each for one agent session and one PR or PR stack.

- Each feature leaves the default branch deployable once merged. Unfinished UI stays unlinked or behind a flag; new API fields land before anything reads them.
- Prefactoring comes first, as its own feature.
- A wide mechanical change (a rename across the codebase) uses expand–contract: add the new form, migrate callers in batches, remove the old form, each its own feature.
- Blocks stay inside the epic. Needing a block on another epic's feature means the epics are cut by layer: tell the person and regroup rather than add it.
- An open question that needs research or a person's call becomes a decision issue (`Decide …` / `Check …`) under the epic, blocking the features that need its answer.

Present the breakdown as a numbered list: title, blocked by, flag, and the behaviour it delivers. Ask whether the granularity and blocks are right. Iterate until the person approves.

Done when the person approves the list and every **Done when** row in the spec is covered by at least one feature.

## 7. Publish

Create the issues blockers first, so each block names a real identifier:

1. Decision issues: sub-issues of the epic, titled `Decide …` or `Check …`, labelled as the tracker adapter says.
2. Features: sub-issues of the epic, named and labelled as the tracker adapter says.

Put every issue in the epic's milestone, in Backlog: the features wait with the epic until a person approves the spec. Set blocks with the tracker's native relation.

A feature's description is its plan, with every section of the tracker adapter's plan template filled. Typically: **Why** ties it to the epic; **What to build** has no file paths or code; every **Done when** item names its check; **How to test** gives steps on a fresh checkout; **Flags** comes from the spec's Features table; **Open questions** is empty unless a decision issue blocks it, and then names that issue.

Fill the spec's **Features** table with the new identifiers. Leave the epic in Backlog: a person approves the spec by moving it.

Done when every issue exists with its blocks, every plan section is filled, and the spec's Features table links each feature. Report the spec URL, the issue list, the open decision PRs and anything left open.
