# Scope and source guidance

For broad questions, group the requested dimensions into useful questions, not predetermined workers. Choose sources by what can resolve each question: specifications for guarantees; code and tests for observed implementation; project decisions for intended local use; studies for measured effects; current documentation for supported versions. Historical evidence may explain a decision but does not establish current behavior.

For code research, cite concrete files/symbols/revisions and observed execution paths; read-only analysis does not authorize source edits. For project context, inspect relevant design docs, ADRs, runbooks and strategic documents. Distinguish current approved decisions from draft proposals. Note missing referenced documents rather than inventing their contents.

For external research, inspect original content rather than relying on snippets. Record URL/title, author or organization when known, publication/access dates and version. Prefer direct measurements and primary specifications to marketing claims. Match benchmark workload, dataset, hardware, concurrency, percentile and time window before comparing numbers. Verify current facts with current sources; do not use fixed historical cutoffs.

When sources disagree, test whether they differ in version, region, configuration, workload or definition. State any remaining conflict beside the affected conclusion. A vendor claim and a measured incident may both describe different configurations. Do not average incompatible results. If a key source is unavailable, narrow the claim or report the missing verification.

A useful research plan names questions, evidence needed, likely sources and the completion condition. Keep it proportional to the scope. Do not require web, code and project searches when only one applies, or expand the topic merely to populate a template.

Carry only dimensions the source specifies into a result: read versus write, percentile, unit, workload, geography and operating mode. Leave omitted dimensions unknown; do not silently specialize a benchmark or guarantee to a plausible operation. A reference beside a statement does not validate added qualifiers.
