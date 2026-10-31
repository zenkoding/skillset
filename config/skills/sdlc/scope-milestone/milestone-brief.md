# Milestone brief template

Fixed headings; one fact per bullet; epic links only in the epics table. Replace every `<...>`. The person approves Goal, Where we are, Non-goals and Done when; the epics table and changelog change freely after that.

Shape rules:

- **Goal**: one sentence in problem terms, not a list of features.
- **Where we are**: each claim about today carries a PR, issue or URL.
- **Rough shape**: 5 to 10 bullets or one diagram; design axioms, if any, as the last bullets.
- **Non-goals**: each with a reason.
- **Done when**: 3 to 5 checks a person can run, each naming how it is checked. The first row is the one-sentence "done when" the tracker milestone carries.
- **Appetite**: a window in weeks, up to a quarter, and what happens when it runs out.
- **Epics**: noun phrases, sentence case; each an end-to-end change someone can see.
- **Length**: one to two pages. Link to the vision and earlier briefs instead of restating them.

```markdown
Date: <YYYY-MM-DD>. Owner: <person>. Target date: <YYYY-MM-DD, or "none">.
Vision: <link to docs/vision.md>. Previous milestone: <brief link, or "none">.

**Next action:** <who does what, one line>

## Goal

<One sentence in problem terms.>

## Where we are vs vision

- **Today:** <what the product does now, with links>
- **Vision:** <the parts of the vision this milestone moves toward, linked>
- **Gap:** <what this milestone closes>

## Rough shape

- <bullet>

## Non-goals

- Everything in the vision's "Not in v1".
- <non-goal>: <reason>

## Done when

| Check | How it is checked |
| --- | --- |
| <observable check> | <URL, command, demo or production observation> |

**Appetite:** <window in weeks>. When it runs out: <cut epics, or extend by a written decision>.

## Epics

| Epic | Outcome | Covers check | Blocked by |
| --- | --- | --- | --- |
| [<noun phrase>](<tracker link>) | <one line> | <row from Done when> | <epic, or "none"> |

## Risks and open questions

- <short>

## Changelog

- <YYYY-MM-DD>: <one line>
```
