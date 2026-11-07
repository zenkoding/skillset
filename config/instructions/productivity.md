## Session start

Before the first substantive response, invoke the `i-have-adhd` skill. If this host has no skill invocation tool, read `~/.agents/skills/i-have-adhd/SKILL.md` and apply it. Keep it active for the whole session, including after compaction, until I say "stop adhd mode" or "normal mode".

## Writing docs for people

Applies to every document a person reads: plans, specs, briefs, RFCs, READMEs and contributing docs. It adds to the writing style above.

- Write for a competent engineer who doesn't know this area. Plain declaratives, "we" for what the team does, full sentences with a subject and a verb. No epigrammatic fragments, no rhetorical set-ups, no idioms where a literal description fits.
- Name roles, not people ("the maintainers of the installer"). Drop dates, except one that identifies a state, such as the commit a fact was checked at.
- Keep implementation detail out unless the doc is about it: file paths, symbols and line numbers belong in agent notes or the code. A crucial algorithm earns a short code block; boilerplate does not.
- Mark an unverified claim where it is made, not only in a list of open questions.
- Draw a diagram only when the structure is the point: a component diagram for architecture (grouped by how a reader should think, at one level), a sequence diagram for a flow across participants, especially where the bug is the ordering. A list of steps is a sentence, not a diagram. Check every diagram renders.
