---
name: find-vulns
description: Hunt for security vulnerabilities in a code area (injection, broken auth, secrets, unsafe deserialization, vulnerable dependencies) and report each evidenced one privately on public repos. Use when asked for a security audit or to find vulnerabilities, or from a security routine.
argument-hint: "[path or area] [--budget N]"
compatibility: Requires a readable checkout; authenticated gh for advisories and issues. Dependency scanners are used when the project has them.
---

# Find vulnerabilities

Find security problems before someone else does, prove each one, and report it where only maintainers can see it. A vulnerability filed in public is a disclosure.

Arguments: the area to audit (default: the whole repo, trust boundaries first) and a budget (default 5, per the [find rules](references/find-rules.md#budget)).

## 1. Map the attack surface

Find where untrusted input enters (HTTP handlers, CLI args, files, environment, messages, third-party APIs) and where it reaches something that matters (queries, shell, file paths, templates, deserializers, auth decisions, secrets). Read `SECURITY.md` when there is one for scope and the reporting channel.

Check the repo's visibility: `gh repo view --json visibility`. It decides where findings go in step 4.

Done when you have the list of entry points and the sinks each reaches.

## 2. Hunt

Trace each entry point to its sinks, looking for:

- **Injection**: SQL, shell, path traversal, template, header, log injection.
- **Auth**: missing or bypassable checks, privilege escalation, insecure direct object references, session and token handling.
- **Secrets and crypto**: secrets in code, logs or errors; weak or home-made crypto; missing verification of signatures or TLS.
- **Unsafe handling**: deserialization of untrusted data, unbounded reads or allocations, SSRF, unsafe redirects.
- **Dependencies**: run the project's scanner (`govulncheck`, `npm audit`, `pip-audit` or what CI uses) and keep only what the code actually reaches.

Done when every entry point has been traced to its sinks.

## 3. Prove

For each candidate, show it is exploitable with the smallest proof that stays local: a unit test, a request against a local build, a reachable call path from the scanner. Never probe a deployed system. A candidate you can't prove stays "suspected" in the report.

Done when every finding to report has a local proof, or a stated reason it can't have one.

## 4. Report privately

Follow the [find rules](references/find-rules.md) for evidence, budget and dedupe, with one change: the destination.

- **Public repo**: open a draft repository security advisory, not an issue: `gh api -X POST repos/<owner>/<repo>/security-advisories` with the summary, the find rules' body as the description, the severity and the affected package when there is one. Its fix goes through the advisory's temporary private fork, never a public PR.
- **Private repo**: file an issue per the find rules with the lane label `security`.
- **No permission for advisories**: stop and give the user the findings in chat for them to report; never fall back to a public issue.

Then report per the find rules, naming each advisory or issue.
