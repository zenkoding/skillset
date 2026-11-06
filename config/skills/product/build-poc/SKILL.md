---
name: build-poc
description: Build a throwaway proof of concept that answers open questions about external behaviour or feasibility before a design is committed. Use when the answer can only come from running code; skip when docs or existing code already settle it.
compatibility: Requires a writable workspace and the tools to build and run the code under test. Parallel subagents are optional.
---

# Build PoC

Turn open questions into answers backed by code that ran. The PoC is evidence, not a product: optimise for learning speed, and expect to delete the code.

## 1. Frame the questions

Collect the questions or hypotheses from the user, an issue, or a research report. Rewrite each as a **falsifiable** question with the observation that would answer it ("Does the API return results in insertion order under concurrent writes?"). Rank them by how much a wrong guess would cost the design, and drop any that docs or existing code already answer.

Done when each remaining question names the observation that settles it.

## 2. Plan the smallest experiments

For each question, pick the smallest runnable experiment: a script, a test against the real dependency, a minimal service. Prefer real systems over mocks, since mocks encode the assumptions under test. Group questions that share a setup.

## 3. Run

Build each experiment in `poc/<topic>/` or a scratch directory the user names, kept out of production code. Run it and capture the raw output. When an observation contradicts the hypothesis, run one variation to rule out a setup mistake before accepting it.

Independent experiments can run in parallel subagents when the host supports them, each given its question, its experiment plan, and the output location. Otherwise run them one at a time.

Done when every question is **answered** (with captured output), **refuted**, or **blocked** (with what is missing).

## 4. Report

Save `docs/spikes/<topic>/REPORT.md` unless the user names another place. For each question: the answer, the evidence (commands and output), and the confidence. Close with the constraints this places on the design and any questions still open.

Return the answers, the design constraints, and the report path. Turning the PoC into production code, or filing GitHub issues for open questions, is a separate request.
