# Plan template

The plan is the document a person reviews and signs off before anything is built. It lives where `docs/contributing/docs.md` says docs live (repo, Notion, Google Docs…) and is linked from the tracker issue. It is authoritative: when it and the companion disagree, the plan wins, and both change in the same pass.

## Full plan

Seven sections, in this order. Drop a section only when it genuinely has nothing in it, and say so rather than padding it.

```markdown
# Plan: <what the change is>

## 1. Context

## 2. Motivation

## 3. Proposal

## 4. Approach

## 5. Risks

## 6. Open questions

## 7. References
```

The title carries the document's type (`Plan:`) so a reader knows what they opened. No status, audience or metadata block under it: those age badly, and the first thing a reader wants is section 1.

- **1. Context** is the before state. Describe the system as it is, then what goes wrong in it. The reader may not know this area: explain the pieces before using them. Finish with how the problem shows up for whoever hits it: the symptom, not the mechanism.
- **2. Motivation** is the why. Rank the reasons by what it costs to leave the problem alone. Say who is affected and how badly. End with one line on whether this is a fix or a new direction, and what has already been agreed.
- **3. Proposal** is the after state. Lead with the change in one sentence, then the components that make it up. Say explicitly what does *not* change: reviewers assume more is moving than is.
- **4. Approach** is the way from before to after, as steps that map to PRs: one step, one PR. Each step leaves the system working on its own, and says in a line why it does. A step that can't stand alone either folds into its neighbour or gets an intermediate behaviour that can (a feature refused is often safe where a feature half-wired is not). Give each step its rough size.
- **5. Risks** covers two horizons. First, what is degraded *between* steps, as a table (`between step n and m | what is degraded | verdict`), and whether that means merging the steps as one. Second, what is risky *once the work is done*: new blast radius, behaviour that changed rather than moved, authorization surfaces, anything left unrepaired on purpose. The settled [work sizing](work-sizing.md) axes go here, with a why for each high one. Decisions the user made are written as decisions ("we accept this"), not as open items.
- **6. Open questions** holds only what is genuinely unsettled. When the design is agreed and only facts remain unverified, say so and list the facts, each tied to the step it affects. Never leave a settled decision here, and never leave an unverified fact asserted in an earlier section.
- **7. References** holds links: issues, PRs, research, related plans, each with a line on why it matters.

## Short plan

For trivial work: three sections, written on the tracker issue itself.

```markdown
### Context
<where we are and what's wrong, two or three sentences>

### Proposal
<the after state, and what does not change>

### Approach
1. <step, one per PR>
```

## Style

Write the plan in the style for docs people read from the user instructions, when they have one. Without it, at least: plain declaratives for a competent engineer new to the area, roles rather than names, no file paths or line numbers (those go in the companion), and unverified claims marked where they are made.

## The companion

The agent's half, in `<store>/notes/<objective>/companion.md`, never committed and never linked from the plan:

- the commit every code pointer was checked at;
- code pointers: paths, symbols, line numbers;
- per-step instructions and the tests each step adds;
- the sizing call (predicted, then settled) and why;
- the decision log: each answer from the grill, marked as the user's decision or as the agent's recommendation until the user rules on it.
