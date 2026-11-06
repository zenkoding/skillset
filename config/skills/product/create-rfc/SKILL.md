---
name: create-rfc
description: Research, interview, and develop an engineering RFC from an idea, notes, or an existing draft. Use when asked to create or substantially develop an RFC; ordinary questions and small document edits do not need this workflow.
---

# Create RFC

Produce one self-contained RFC that engineers outside the conversation can review. Research establishes facts; the interview resolves consequential decisions. Maintain the same document throughout, using the [Rust-based RFC template](references/rfc-template.md). This skill ends with RFC review, not implementation or publication.

## Establish the starting point

Read the supplied material, relevant repository guidance, and existing RFC conventions. Identify the problem, audience, constraints, settled decisions, and consequential gaps. An existing draft is the starting artifact: preserve its accepted decisions and focus on contradictions or missing reasoning.

Distinguish the broader feature opportunity from this RFC's design scope. An umbrella issue may describe use cases that require several RFCs; identify which this RFC covers and preserve deferred findings in existing private notes or authorized tracker records without designing sibling proposals. Follow an issue → acceptance → RFC sequence where the repository requires it; otherwise continue the authorized RFC work without introducing that gate.

Separate design scope from delivery milestones. Establish from the request whether the user wants a full design delivered incrementally or a deliberately limited proposal; ask only when unclear. A smaller first release need not narrow the RFC.

Follow the repository's RFC location and numbering convention. Otherwise use `rfc/NNN-title.md`, selecting the next available number starting at `001`. Outside a repository, use a suitable user-visible working directory. Keep subsequent revisions in the same file. If the destination is shared, prepare a private copy until publication is authorized.

Create an initial RFC with established facts and clearly marked open questions. Use it as the single design document rather than creating separate glossaries, context documents, ADRs or research reports. Place terminology near the end, followed by a references section when needed. Required private execution notes may hold resumption state, but must not become an alternative design source. Keep incidental experiment artifacts out of the RFC's narrative.

## Research and interview in rounds

Map the design as decisions and their dependencies. In each round, ask the independent consequential questions whose prerequisites are settled; defer questions that depend on unanswered ones. Use manageable batches, stable question numbers, concrete scenarios, and a recommended answer with its main tradeoff. Use the host's question tool when suitable, otherwise chat. Allow the user to propose a different answer.

Investigate factual questions yourself. Read relevant code, primary documentation, and available connected sources before asking the user to resolve uncertainty. Use bounded independent research agents when available and useful; otherwise perform the research directly. Keep synthesis and decisions with the lead. Do not delegate a product decision to a research agent.

Small, reversible local experiments may resolve a design question. State what the experiment tests, preserve existing work and data, and stop once the evidence answers it. Ask before substantial environment setup or implementation; keep external mutations within explicit authorization. An RFC request alone does not authorize deploying a prototype, changing shared systems, or messaging others.

Distinguish source-supported capability, observed behavior, proposed design, and unverified integration. Evidence from one build or fixture establishes only what was exercised. Keep research provenance in working context or required private notes. In the RFC, cite only sources that help a reviewer assess a material claim, following the reference policy below. Disclose a relevant coverage gap rather than converting it into a universal limitation.

After each answer, update the RFC immediately and recompute the open decisions. After scope changes or accepted exceptions, reconcile prose, command examples, schemas, diagrams and open questions together. Remove superseded behavior and check that each exception has the same boundary everywhere. Preserve the current design rather than appending an interview transcript. Do not reopen settled choices without new evidence or a contradiction; explain that evidence when reconsideration is needed.

Grill behavior, scope, interfaces and meaningful tradeoffs. Choose routine, reversible implementation details with engineering judgment. Leave technical checks that require implementation in unresolved questions, with the concrete evidence needed to resolve them. Exhaustive questioning about incidental details is not the objective.

## Write for fellow engineers

Read the template when drafting or restructuring. Lead with the proposal and its effect, then follow the template's progression through motivation, scope, proposed UX and proposed architecture. Begin architecture with components and ownership boundaries before individual flows. Adapt headings to the proposal and repository conventions. Use direct, familiar language and connected paragraphs. Include technical detail where it helps assess a contract, tradeoff or failure mode; define domain-specific terms in the terminology section. Label illustrative schema fields or routes as provisional until selected.

For repository RFCs, keep reader-facing references within that repository unless the user requests external references. External documents, meetings and issue trackers may inform research; incorporate their relevant decisions and reasoning so readers can understand the proposal without opening them. Use descriptive links to repository files or pinned source revisions near material claims. Keep raw hashes, local checkout paths and source inventories in working notes rather than the prose. Apply the user's reference policy when working outside a repository.

Use Mermaid when a diagram explains a relationship more clearly than prose:

- Component views for architecture, with subgraphs for meaningful ownership or deployment boundaries. Keep external systems opaque unless their internals are verified, relevant to the decision and appropriate for the audience.
- Sequence diagrams for request, callback and background flows; show ownership and observable interactions. Include handler or persistence details only when they explain a relevant contract.
- State diagrams for lifecycle transitions and terminal states.

Include concrete examples in the first substantial draft. For CLI proposals, demonstrate the main journeys with commands and expected behavior. For resource proposals, prefer one annotated example covering the proposed fields over overlapping snippets and field inventories; distinguish writable configuration, server observations and concurrency metadata. Mark unknown mappings explicitly rather than implying schema completeness. Keep examples consistent with permissions, identity, lifecycle and error semantics. Do not invent an existing framework or component merely to label a diagram. Use tables for genuinely comparable concepts, not as a substitute for explanatory prose.

Write a snapshot of the proposed design. Give each rule or rationale one clear home; elsewhere use only the short restatement needed to understand a flow. Retain alternatives only when they explain a real engineering tradeoff. Remove abandoned implementation commentary, conversation history, repeated defensive contrasts, and task-management restrictions from the reader-facing RFC. Localize material limitations and distinguish accepted behavior from implementation questions and future work. Combine or omit template sections when that improves the reading order; scale length to the decisions reviewers must assess.

## Complete and review

The interview is complete when consequential decisions are resolved or explicitly deferred with the user's agreement, and remaining technical questions have concrete validation criteria. Do not treat silence as agreement.

Perform a final editorial pass from the perspective of an engineer who has not seen the conversation. Check that the proposal and rationale stand on their own, references follow the chosen policy, and repeated evidence or caveats have been consolidated. Repeat the consistency pass against accepted decisions: prose, commands, schemas, diagrams and open questions must agree on scope, ownership, lifecycle rules and exceptions. Check that shortening preserves requirements and that settled decisions are no longer presented as open questions. Verify local links and example syntax where practical; validate or render diagrams when tooling is available. State any material validation limitation without implying a check ran. Use an independent bounded review when it adds meaningful confidence.

Present the complete RFC for the user's review, naming any remaining decisions or technical gaps. Ask whether it captures the shared understanding; revise the same artifact until it does. Do not require another approval for each routine draft edit. Implementation, commit/push, PR creation and shared publication are separate actions requiring their own authorization; honor authorization already given.
