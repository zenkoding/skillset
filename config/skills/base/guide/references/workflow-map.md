<!-- every skill named here: when installed -->
# Workflow map

Every skill, by loop. Name a skill to the user only when it is installed; a loop with nothing installed is left out.

## Product loop (product set)

Where the project is going, a milestone at a time.

```
envision ──> scope-milestone ──> shape-epic ──> (build loop, per feature) ──> retrospect
VISION.md    milestone brief     epic spec +                                   lessons → skills,
                                 feature tickets                               docs, next milestone
```

- `envision`: turn a rough idea into a `VISION.md` change, with research and a mockup or PoC (`build-poc`).
- `scope-milestone`: pick the next milestone from the gap between today and the vision.
- `shape-epic`: turn one epic into a spec and vertical-slice feature tickets.
- `retrospect`: after an epic or milestone, turn friction into skill, doc and check fixes; runs `codify` for repeated mistakes.
- `create-rfc`, `update-rfc`: write or revise an RFC for a contract other code depends on.

## Build loop (base and engineering sets)

One ticket, from problem to merged PR.

```
brief ──> grill ──> build ──> publish ──> review ──> address ──> (merge) ──> debrief
plan +    settled   tdd,      PR with     labelled   fixes,                  decisions,
companion decisions checks    reviewer    findings   re-review               friction
                              notes
```

- `brief`: research, draft, grill and sign off a plan; trivial work gets a short plan on the issue.
- `grill`: settle open decisions one question at a time, with a recommendation each.
- `build`: build one ticket test-first (`tdd`) to an open PR.
- `publish`: open the PR with What, Why, How to test and Reviewer notes; mark it ready once CI is green.
- `review`: review a PR and post labelled findings (Crucial, Important, Suggestion, Question, Hint, Remark).
- `address`: fix what review found and what CI reports, refresh the PR body, hand back for re-review. Never merges.
- `test`: run the plan's How to test steps against a fresh build and post a verdict per step.
- `debrief`: at the end, report decisions, lessons and friction, and file them.

## Engineering (engineering set)

- `fix-bug`: find the cause of a defect before anyone fixes it.
- `optimize-performance`: measured performance work.
- `triage-prs`: which open PRs need attention, and the command for each.
- `codify`: turn a repeated mechanical mistake into a check.

## Setup and knowledge (base set)

- `setup-repo`: write or migrate `CONTRIBUTING.md` and `docs/contributing/`; `e2e` and `review` parts on demand.
- `research`: a source-backed report, saved privately.
- `project-knowledge`: set up or tidy the private store.
- `guide`: this map, and what to run next.
