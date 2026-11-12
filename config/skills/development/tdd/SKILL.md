---
name: tdd
description: Work test-first in red, green, refactor slices, with the plan's Done when checks as the seams and tests at public boundaries. Use when building a feature or fixing a bug test-first, or when build reaches its test-first step.
---

# Test-driven development

Grow the change one behaviour at a time: a failing test, the least code that passes it, then a cleanup with the suite green. The tests describe what the code does through its public interface, so they survive refactors.

## 1. Find the seams

- **From the plan.** The plan's **Done when** checks are the agreed seams: each becomes one or more tests at the boundary it names. A seam the plan doesn't name (a new public function, endpoint or flag) goes back to the plan's Open questions before you test against it.
- **Without a plan** (a quick fix or the user's direct ask), name the seams in one line each and confirm them with the user before writing tests.
- **Commands from the repo.** Take the test and check commands from the repo's docs (`AGENTS.md`, `CONTRIBUTING.md`, the testing docs or the `Makefile`), never a guess. When none is written down, ask once.

Done when every seam has a named boundary and you know the command that runs its tests.

## 2. Red

Write one failing test for the next behaviour on the list. Start with the simplest case that forces real code, not the edge cases.

- Test behaviour through the public interface: inputs, outputs, observable effects. Never assert on private state or call order.
- Use real collaborators where they're cheap. Fake only the edges: network, clock, filesystem, randomness, third-party services.
- Use table-driven tests where the language makes them idiomatic (Go, for one), one row per behaviour.
- Name the test after the behaviour, in the reader's words.

Run it and watch it fail for the reason you expect. A test that passes at once, or fails for another reason (a typo, a missing import), proves nothing: fix the test until it fails on the missing behaviour.

## 3. Green

Write the least code that makes the test pass. Hard-coding is fine when the next test will force the general case. Run the whole suite, not just the new test. Done when everything passes.

## 4. Refactor

With the suite green, clean up: remove duplication, rename for the reader, extract what the tests now show is one idea. Change tests and code in separate steps, and run the suite after each. Never add behaviour here: new behaviour starts at red.

Repeat red, green, refactor until every seam's behaviour is covered.

## Bug fixes

A bug fix starts with a test that reproduces the bug through the public interface. Run it on the base branch and confirm it fails there for the reported reason; that failure is the proof the test covers the bug. Then fix, and keep the test as the regression guard.

When the bug can't be reproduced in a test, say so, name the seam that would have to exist to test it, and stop before fixing blind.

## Done

Done when every **Done when** check (or confirmed seam) has a passing test, the repo's check command exits 0, and no test asserts on internals.
