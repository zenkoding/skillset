# Review criteria

What a change is judged against, in order of importance. `review` judges a pull request against them. `address` judges each review finding against the same criteria, so a finding is acted on only when it holds up against the same bar the review used. An earlier criterion outranks a later one: a change that doesn't belong needs no bug hunt.

The sources: the vision (`VISION.md`), `CONTEXT.md`, `ARCHITECTURE.md`, the PR's intent (its linked issue or plan, and the PR body), the repo's rules (`REVIEW.md`, `AGENTS.md`, `CONTRIBUTING.md` and the docs it points to) and the decision history of the area (`git log -S '<term>'`).

1. **Product fit.** Does the change belong in this project? Check it against the vision, non-goals and `CONTEXT.md`.
2. **Rationale.** Does it add a useful feature, make a necessary enhancement, or fix a real issue? Give each acceptance criterion and each Reviewer notes Must check item a verdict: met, partial, missing, or not checkable from the diff. Note scope nobody asked for.
3. **Approach.** Does it reach its goal in a reasonable way? Weigh complexity (less is better), consistency with the codebase and its architecture, and verbosity (less code is better while it stays easy to read). Check the repo's written rules here. When Complexity is high, check that each new abstraction or cross-component change is earned. When Taste is high and nothing (issue, PR body, git history) says who picked the option or why, ask. When the change alters a decision, check its **Why** says Chose X over Y because Z and that it edits the living doc the decision changes.
4. **Security and bugs.** Trace each changed path for a concrete input that breaks it: error paths, boundaries, leaks on failure, concurrency, callers of changed behaviour. Run a separate security pass when Risk is high. Note missing or weak tests for new behaviour.
