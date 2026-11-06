# Epic spec template

The body of an epic's spec doc in the docs store. One page, hard cap two. Fill every heading; write `none` rather than dropping one. Keep every heading in this order so people and agents find the same thing in the same place. Epic, milestone and tracker link are page properties (or front matter, per `docs/agents/product-docs.md`), not body lines.

```markdown
Date: <YYYY-MM-DD>. Purpose: <one line on what this spec settles>.
Sources: <milestone brief, research docs, ADRs this spec rests on>

# <Title: what someone can do once this ships>

**Next action:** <who does what next, one line>
**Appetite:** <days or weeks; when it runs out we stop and re-decide>

## Problem
One specific story of who hits the problem today and what it costs them. A story, not a category.

## Done when
| Criterion | Verified by |
| --- | --- |
| <observable criterion> | <test path, CI check, URL + status, metric threshold, or "person judges"> |

## Scope
- <what this epic delivers>

## Non-goals
- <what it does not deliver>: <why>

## Rabbit holes
- <known trap, or a decision already made so nobody reopens it>

## Features
| Feature | Blocked by | Flag and stage |
| --- | --- | --- |
| [<id>](<url>) <imperative title> | <id, or none> | <flag, stage; or none> |

## Decisions
- <ADR / RFC / decisions entry>: [<title>](<repo link>), or none

## Changelog
- <YYYY-MM-DD> <one line>
```

## Rules for the Done when table

- Every row names its check: a test path, a CI check name, a URL with the expected status, or a metric threshold. A row only a person can judge says `person judges`.
- A row that cannot be tested before ship gets a date or a fail condition.
- By convention, word rows checked after deploy as an `http`, `metric` or `check` kind, after ADR 0012. That ADR is proposed and may change; follow it once it settles.
- A person changes Criterion. Nothing tracks whether a row is met yet; the feature plans' checkboxes carry that until a skill does.

## Flag stages

When a feature hides behind a flag: stage 1 team only, stage 2 on by default, stage 3 flag removed.
