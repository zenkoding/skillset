# Optional generic worker briefs

Use only when useful workers are actually available; otherwise execute locally. Each brief must include PR identity, exact checkout/revision, user authorization, applicable project rules, concrete input evidence, allowed write scope, output artifact or structured response and stopping condition. Include document contents when workers cannot access the filesystem. Do not assume inherited conversation or request a registered role.

| Task | Inputs and allowed work | Required result |
|---|---|---|
| Scout | PR metadata, current-head checks/logs, unresolved review threads, body and merge state; read-only | CI_FAILURES with relevant logs/SHA; PENDING_COMMENTS with source/path/line and resolution evidence; PR_BODY; MERGE_STATE; inaccessible evidence explicitly unknown |
| Project checks | Applicable instructions and build configuration; read-only | Exact PRE_COMMIT_COMMANDS, PRE_FLIGHT_COMMANDS, COMPLIANCE_DOCS; missing commands distinguished from suggested existing targets; human judgment steps listed separately |
| CI fix | Actual failure output, relevant source/check command and compliance docs; only assigned files | Minimal patch, actual check result, unresolved cause and changed paths |
| Review fix | Clear/context-supported feedback, source context, PR purpose and compliance rules; only assigned files | Each addressed request mapped to its change/evidence; flagged requests remain unresolved |
| Verify | Final revision, actual project commands and required manual gates; only expected test-generated files | Per command PASS/FAIL/SKIP/TIMEOUT plus output; MANUAL gates; revision checked; no success if required work is missing |

Run overlapping writers sequentially (CI fix before review fix when both touch the same files). Inspect returned patches and evidence. A failed worker or absent expected result cannot satisfy a gate; retry only when useful new evidence supports it, or report the blocker. Same-agent verification is self-review. If independent verification is explicitly required, it remains unmet without a separate reviewer.
