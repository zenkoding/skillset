# docs/agents/e2e.md format

The adapter `test-e2e` reads to run a repo's app. One section per stack mode; the first is the default. Every command runs from the repo root on a fresh checkout.

```markdown
# E2E stack

## <mode, e.g. host>

When to use: <e.g. cloud routines; no containers>

**Build**
<commands>

**Start**
<command; how to run it in the background and where its output goes>

**Health**
<command that exits 0 once the app is ready, with a timeout>

**URL**
<base URL for UI steps and API calls>

**Flags**
<how to set a feature flag to a stage, or "none">

**Stop**
<command>
```

Keep secrets out: name the environment variables a mode needs and where they come from, never their values.
