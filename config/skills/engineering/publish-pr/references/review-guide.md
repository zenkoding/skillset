# Review guide

An optional `### Review guide` section at the bottom of a PR body, after How
it was tested. It tells the reviewer where the risk is and what to confirm,
so attention goes to the few files and decisions that matter.

## When to offer

Offer only when creating the PR, and add the guide only on a yes. Updates refresh an existing guide but never offer one. An explicit request ("publish with a review guide") adds it at any size, on create or update. Offer when either holds:

- **500+ counted changed lines.** Count additions plus deletions, excluding
  lockfiles, generated or vendored files, pure renames and test fixtures.
- **A risky area at any size**: public contracts or APIs, migrations or
  schema, auth, secrets or permissions, dependency manifests, or code many
  callers share.

The offer is one line: "This PR is large/risky (<reason>). Add a review guide?"

## Template

Fill every section that has something real to say; drop the rest.

```markdown
### Review guide

**Verdict**: ⚠ careful review: <one line on where the risk is>

**Main decisions**
- <chose X over Y because Z> ([ADR-NNNN](https://github.com/<owner>/<repo>/blob/<head-branch>/docs/adr/NNNN-....md))

**How it works**
<components touched and how they connect; a small diagram when it helps>

**Contentious points**
- <unusual choice a reviewer will ask "why like this?" about>

**Must check**
- [ ] <breaking contract / migration / security-sensitive change / dep bump>: confirm <what>

**Where to focus**
- High: `path` <why>
- Medium: `path` <why>
- Skim: <the rest, summarised>

**Merge risk and rollback**
<what could break, blast radius, how to roll back>

**Test by hand**
- [ ] <flow or command to try>
```

## Verdicts

- **✓ fast track**: mechanical or isolated; reading the diff once is enough.
- **◐ quick review**: real logic, contained blast radius, covered by tests.
- **⚠ careful review**: touches a risky area, changes behaviour callers rely
  on, or tests leave a gap.

## Rules

- **Main decisions** quotes the session's `debrief` Decisions bullets when
  there are any, and links `docs/adr/` entries the branch adds. Never link
  `.docs/` paths or the private store (`~/.local/share/agents/`): they are
  private and absent from the repo.
- Every path in Where to focus is in the PR diff.
- Links in the guide are absolute `https://github.com/<owner>/<repo>/blob/<head-branch>/` URLs; relative links in a PR body resolve under `/pull/` and 404.
- Every Must check item names what to confirm.
- When `publish-pr` updates the PR, rewrite the guide against the final diff;
  the rest of the body follows the usual preservation rules.
