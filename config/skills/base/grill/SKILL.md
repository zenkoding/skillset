---
name: grill
description: Clarify scope before anything is built. Walk the design tree one decision at a time, recommend an answer for each, settle what the evidence can, and log every answer to the objective's notes. Called by brief, envision, scope-milestone and shape-epic; run alone as /grill <plan, idea or question>.
---

# Grill

Turn a plan with open questions into shared understanding. The goal is a short list of settled decisions the user agrees with, not an exhaustive interview: ask only what changes the work.

## 1. Seed the tree

The root is the goal in one sentence. Its branches come from:

- what the caller left open: brief's Unknown list, a template's sections, research gaps;
- each high axis of [work sizing](references/work-sizing.md):
  - **Complexity**: the design branches (components, contracts, data, failure modes);
  - **Taste**: the options a person picks from (UX, naming, output format, API shape);
  - **Risk**: rollback, blast radius, and every one-way door.

Grill a draft, not a blank page: when there is none, write a short one first (the problem, the proposal, the open decisions) so questions have something to push on. Run alone, read the plan or idea first and size it yourself.

## 2. Settle what the evidence can

Before asking, look for an answer in the code, the repo's docs, the tracker and the history (`git log -S '<term>'` on the doc or code it touches). A question the evidence settles is logged with its source, not asked. Ask only what the user, and nothing else, can decide.

## 3. Ask

- **One decision per question.** Lead with the recommended answer and a one-line why, then the alternatives.
- **Native question tool first.** When the host has one (Claude Code's `AskUserQuestion`, for example), put one question per card. A round can hold a few independent questions; a question that depends on another waits for the next round.
- **Fallback: a numbered round in chat.** Up to three independent questions, each as `1. <question> Recommend: <answer> (<why>). Else: <option>, <option>.` The user answers by number; "ok" accepts every recommendation in the round.
- **Taste is the user's pick.** Give the options and your recommendation, then wait for the pick. Never settle a Taste question on your own.
- **Risk names the door.** A one-way door needs the user's explicit yes, not an accepted default.
- **Check the premise.** Before asking about a mechanism, check it works the way the question assumes. "Should we refuse X or park it?" wastes a round when nothing can answer a parked X.
- **Recompute after each round.** An answer can reshape the tree and reopen something you treated as settled; rebuild the open list rather than carrying the old one.
- **Walk depth-first.** Each answer can open child questions. Take the branch that would change the plan most first, and stop descending when an answer is a routine, reversible detail you can decide yourself (log it as yours).

## 4. Log as you go

Append each answer as it lands to `<store>/notes/<objective>/grill.md`, by the absolute store path the session-start hook prints (see "Project knowledge" in your user instructions):

```markdown
- <YYYY-MM-DD> <question>: <answer>. Decided by <user | evidence: source | agent: routine | agent: recommended>.
```

Keep your own recommendations marked as yours (`agent: recommended`) until the user rules on them: later work quotes this log, and a proposal logged as a decision becomes one by accident. A resumed grill starts from this log and asks nothing already settled. When there is no store, keep the log in the chat and say so.

## 5. Close

Stop when the frontier is empty or the user calls it. Then:

1. Re-size the work: settled answers can raise or lower an axis, and a non-trivial call can turn trivial.
2. Post the result: the decisions (one line each), anything still open with who settles it, and the settled sizing with a why per high axis.
3. Ask the user to confirm this is the shared understanding. Nothing is built before they do.

Inside `brief`, the answers land in the plan as the [plan template](references/plan-template.md) says: decisions in Proposal and Risks, unsettled facts in Open questions, and the log in the companion.

Done when the user has confirmed the summary and every answer is in the log. Hand back to the skill that called the grill, or name the next step when run alone.
