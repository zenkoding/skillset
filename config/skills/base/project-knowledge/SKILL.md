---
name: project-knowledge
description: Set up or tidy a repo's private project knowledge store and its catalog. Use when the user asks to bootstrap the store, or to trim, split or merge its docs.
---

# Project knowledge

Partitions, the store path, precedence and write policy are in "Project knowledge" in your user instructions; this skill only changes the docs' shape. The private store is yours to edit, by the absolute path the session-start hook prints. Public docs are human-first (`CONTRIBUTING.md` and `docs/contributing/`); `/setup-repo` writes them when installed, and changes to them are proposals until the user says to write them.

## Catalog format

The catalog is `README.md` in the store: a title, one sentence on what the partition holds, then one line per doc:

```markdown
- [install.md](install.md): how the installer links, copies and merges. Open when touching the installer or settings.
```

Each line names what the doc covers and the trigger for opening it. Notes under `notes/` stay out of the catalog.

## Bootstrap

1. Find existing knowledge: `CONTRIBUTING.md`, `docs/`, `AGENTS.md`, `CLAUDE.md`, `.claude/docs/`, and the old private notes in `~/.local/share/context/notes/` that match this repo.
2. Write `<store>/README.md` cataloguing private docs; move matching old notes to `<store>/notes/<objective>/` when the user agrees. Offer a `<repo>/.docs` symlink to the store for humans (globally git-ignored); agents still write by the store's real path.
3. A fact a human contributor also needs belongs in the public docs, not the store: list it as a proposal for `CONTRIBUTING.md` or `docs/contributing/` (via `/setup-repo` when installed).
4. Done when every catalog path exists and every private doc is either listed or deliberately left out (say which and why).

## Trim

Go doc by doc through the store's catalog:

- Cut lines the code, config or `--help` already answer, stale facts, and session history. Keep intent, decisions, gotchas and conventions.
- Split a doc that covers more than one subject or no longer reads in one go; merge docs that cover the same subject.
- Rewrite each changed doc's catalog line in the same edit.

Done when every doc has exactly one catalog line, every catalog line points at an existing doc, and the catalog fits within 200 lines (the session-start hook cuts it there). Report what moved, merged or went, per doc.
