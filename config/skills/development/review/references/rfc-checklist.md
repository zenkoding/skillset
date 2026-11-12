# RFC checklist

For a PR that adds or changes an RFC (under `docs/rfcs/`, or wherever `CONTRIBUTING.md` says RFCs live). The RFC is the change, so the four areas of the review apply to the design it proposes, not to code. Use these checks in place of the code-focused ones; label findings as usual.

1. **Product fit.** The problem is real and belongs in this project: check it against `VISION.md` and the non-goals. The RFC says who is affected and why it matters now.
2. **Rationale.** Scope and exclusions are explicit. The motivation is a concrete problem with scenarios, not a solution in search of one.
3. **Approach.**
   - The proposed UX (or "How it works") teaches the design through a concrete example before implementation detail.
   - The architecture names components and ownership boundaries, then contracts, data, permissions, failure behaviour and concurrency where they apply.
   - Drawbacks are honest, and the alternatives are real choices with the reason each lost.
   - It is consistent with `ARCHITECTURE.md`, `CONTEXT.md` terms and the decisions already in git history (`git log -S '<term>'`). A contradiction is either argued for or fixed.
4. **Risks and gaps.**
   - Validation lists acceptance criteria and separates observed evidence from checks still required.
   - Unresolved questions are listed with what would settle each; no assumption is presented as a decision.
   - Rollout or migration is covered when the change touches stored data, a public contract or many callers.

Also check that it reads as one self-contained design: no transcript of how it was drafted, terminology near the end, and references that support the claims they sit next to.
