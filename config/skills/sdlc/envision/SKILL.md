---
name: envision
description: Take a rough product idea through research, a mockup or PoC, and a grill, to an updated vision in the vision brief format, proposed as a PR. Run as /envision <idea>.
disable-model-invocation: true
---

Product-lane skill: turn a rough idea into a change to the project's vision. The output is an edited vision doc (`docs/vision.md` unless the docs catalog names another) in the **vision brief** format ([`vision-brief.md`](vision-brief.md)). Vision changes wait for a person, so this skill opens a PR and stops for approval.

Read the vision doc, the project-wide decisions doc when the repo has one, the glossary (`CONTEXT.md`, per `docs/agents/domain.md`) and `docs/agents/product-docs.md` before step 1. Name things with glossary terms.

## Steps

1. **Frame the idea.** Restate it in two or three sentences: who it is for, what hurts, what changes for them. Name where it touches the current vision: which sections it extends, contradicts or leaves alone. Done when the user confirms the framing.

2. **Research.** List the open questions the vision change rests on (market, users, competitors, prior art, technical feasibility). Investigate each against primary sources with the `research-scope` skill, in background agents when the host has them, and write the findings to one Research doc in the docs store as `docs/agents/product-docs.md` says, with every property it asks for (type, milestone, epic, tracker link, date). With no docs store, save the report where `research-scope` defaults. Each claim carries a confidence mark; the doc ends with what could not be verified. Done when every listed question is answered or marked unknown.

3. **Make it concrete.** Build one thing a person can look at: a UI mockup or a logic PoC (`build-poc`, or a prototype skill when installed). It is throwaway and lives on a throwaway branch; it answers "does this feel right?", not "does it work in production". Link it from the Research doc and, later, the PR. Skip only when the user says the change is wording-only (a principle or non-goal reworded). Done when the user has seen it and said what they think.

4. **Grill.** Run a grilling skill when installed (otherwise interview the user one question at a time) on the drafted change, seeded with the research and the mockup. Push hardest on: who the first audience is, the diagnosis and guiding policy, what we won't do, which principles break ties, and the top risks. Record every settled answer. Done when the design tree has no open frontier, or the user calls it.

5. **Write the brief.** Edit the vision doc to follow the template, within its page budget. Keep the existing wording wherever the grill did not change it. On a first run, map the doc's current headings onto the template's sections and fill the new ones from the grill. Add a dated Changelog line naming the change and linking the research. Done when every `<...>` placeholder is replaced or its section deliberately left empty, and every claim traces to the grill or the research.

6. **Propose.** Create a branch from the default branch, named with the vision epic's tracker identifier when there is one, commit the vision change, and push. If the change adds or renames a glossary term, commit that edit on the same branch: both need the same person's approval, so they go in one PR. Open the PR with `publish-pr`; its **Why** links the vision epic, or says why there is none. Link the research and the mockup in the body. Done when the PR is open and its URL is reported; stop there for a person's review.

Milestones and epics are the next skills' job (`scope-milestone`, `shape-epic`). The brief's Goals section links milestones that already exist and names new ones only as one-line candidates.
