---
name: plan-execute
description: Plan and execute multi-step changes with dependent steps, cross-layer contracts, parallel slices, or high risk. Lightweight in-session plans and direct implementation do not need this skill.
---

# Plan & Execute

Plan only when dependencies, contracts, risk, or handoffs justify a durable artifact; otherwise use a lightweight in-session plan:

```markdown
## Goal
## Assumptions
## Steps
1. [change + focused verification]
## Final validation
```

Save durable plans to `docs/plans/YYYY-MM-DD-<feature>.md`.

## Full plan (parallel / high-risk / cross-session)

First inspect the actual code, tests, Makefile targets, and conventions; record unknowns instead of inventing architecture. Break work into vertical slices — one complete, independently verifiable behavior, never a horizontal "write all handlers" phase.

Plan header: **Goal** (one sentence), **Architecture** (fit with existing layers), **Validation** (Makefile commands), plus a dependency table (Task | Type AFK/HITL | Blocked by | Parallelizable with).

Per task, include only what applies:

```markdown
### Task N: [slice name]
Type / Blocked by / Areas (layer terms)
Goal: one or two sentences
Acceptance criteria: observable behavior + key error/edge cases
Files: exact Create/Modify paths
Contracts: signatures or schemas downstream tasks need — no implementations
Tests: primary + boundary behavior; E2E only if the harness covers it
Validation: focused test command + make target
Risk controls (high-risk only): failure modes, rollout/migration order, rollback trigger and steps, data compatibility
```

Coverage check before handoff: every requirement mapped, no dependency cycles, parallel tasks share no files, commands come from the Makefile.

## Wave execution (parallel slices)

Dispatch subagents only when all hold: per-slice goal/criteria/allowed-files/validation; acyclic dependencies with stable prerequisite contracts; non-overlapping write sets; parallelism saves more than context-transfer costs. Otherwise execute sequentially.

Per wave:

1. Resolve blocking user decisions before dispatch.
2. Dispatch the fewest useful non-conflicting implementers. Give each the full task text, allowed files, project conventions, and exact contracts from completed prerequisites.
3. Inspect actual diffs, not agent reports; run each slice's focused validation plus one integration check per wave.
4. Fix integration failures with the narrowest owner; escalate after repeated failures rather than looping.

The coordinator verifies each slice's acceptance criteria against the actual diff — there is no separate slice-level spec reviewer. High-risk or public-contract changesets get exactly one independent `code-reviewer` pass at the end (per AGENTS.md), never per slice. Never: parallel writers on the same file, silent write-set expansion, or claimed validation that did not run. For dependent work, pass the generated artifact's exact declarations — if they differ from the plan, stop and resolve before dispatching downstream.

## Handoff

Confirm acceptance criteria, run integrated validation once (`make all` or the risk-appropriate subset), inspect the combined diff, report deviations and residual risk. Changeset review belongs to the invoking workflow; do not duplicate it.
