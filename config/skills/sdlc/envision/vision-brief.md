# Vision brief template

The format for `docs/vision.md`. About 1 to 1.5 pages in total. Living: re-read at each milestone close; every change gets a changelog line and a person's approval. No front matter. Keep the headings as written so people and agents find the same section in the same place, and replace every `<...>`.

```markdown
# <Product> vision

<One sentence: the world once this exists.>

**Status:** <what is open right now, and the next decision and who makes it>

## Pitch

<3 to 5 sentences a customer would understand.>

## Problem

<Who hurts and how. 2 to 4 bullets, each a concrete pain, not a category.>

## Who it's for first

<The first audience, specific enough to say no to others. Roles inside it, one line each.>

## Diagnosis and guiding policy

- **Diagnosis:** <why the problem exists; the one thing that, if true, explains it>
- **Guiding policy:** <our approach to it, in one or two sentences>
- **What we won't do:** <the tempting approaches this policy rules out>

## What we're building

<3 to 5 numbered principles. Each is a tie-breaker an agent can apply to a real choice, in bold, then one or two sentences.>

## Not in v1

<Explicit non-goals, one line each, with a reason where it is not obvious.>

## Goals

<What success looks like, as observable outcomes. Each links the milestone that delivers it, or is marked "candidate" when no milestone exists yet. Include the evidence bar to clear before building, if any.>

## Top risks

<Three at most: what could make this wrong, and how we'd notice.>

## Changelog

- <YYYY-MM-DD>: <one line: what changed and why; link the research doc>
```

## Section rules

- **Principles** are tie-breakers, not values. Test: can an agent use it to pick between two real options? If not, rewrite or cut it.
- **Non-goals** may grow freely; removing one is a vision change that needs its own changelog line.
- **Goals** link down to milestones; a milestone never restates the vision, it links up.
- Vision and milestones change at different speeds. Plans, dates and epic lists belong in milestone docs, not here.
