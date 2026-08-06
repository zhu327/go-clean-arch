---
name: subagent-driven-development
description: Execute an approved plan with parallel implementation subagents when independent slices have stable contracts and non-overlapping write sets. Do not use when coordination cost exceeds plausible parallel benefit.
---

# Parallel Plan Execution

Use fresh implementer subagents only when a written plan contains genuinely independent work. The coordinator may implement directly when tasks are coupled, small, or dependent on evolving context.

## Entry criteria

Before dispatching, confirm:

- each slice has a clear goal, acceptance criteria, allowed files, and focused validation;
- dependencies are acyclic and prerequisite contracts are stable;
- concurrent slices do not create or modify the same files;
- the repository's actual commands and conventions have been identified;
- parallel execution is likely to save more time than context transfer and integration cost.

If these conditions do not hold, execute sequentially or return the concrete planning gap. Do not invent contracts to enable parallelism.

## Execute by dependency waves

For each wave:

1. Resolve any blocking user decision.
2. Dispatch only non-conflicting implementers. Use the smallest useful number rather than a fixed quota.
3. Give each agent the complete task, allowed files, relevant project conventions, and exact contracts from completed prerequisites.
4. Wait for the wave and inspect actual changes—not only agent reports.
5. Run focused validation for each slice and an integration check appropriate to the wave.
6. Fix integration failures with the narrowest owner. Stop and escalate after repeated failures rather than looping indefinitely.

Use `implementer-prompt.md` as a starting template, adapting it to the task.

## Independent compliance review

Require a separate spec-compliance reviewer when a slice is high-risk, crosses a trust/ownership boundary, changes a public contract materially, or cannot be validated convincingly by tests. Otherwise implementer self-check, coordinator inspection, and passing tests are sufficient.

Use `spec-reviewer-prompt.md` for those reviews. Review only the affected slice and repeat only when findings require a fix.

## Artifact passing

For dependent work, read generated artifacts and provide the exact relevant declaration, schema, configuration, or contract. If it differs from the approved plan, stop and resolve the mismatch before dispatching downstream work.

## Final integration

After all waves:

- confirm planned acceptance criteria and cross-slice contracts;
- run the repository's applicable integrated validation once;
- inspect the combined diff for scope and accidental overlap;
- report deviations, unavailable checks, and residual risk.

Changeset-level code review is owned by the invoking workflow; do not dispatch a duplicate global reviewer here.

## Safety rules

Never run parallel writers on overlapping files, let an agent expand its allowed write set silently, or claim validation that was not run. Destructive operations and main-branch policies follow the repository's instructions and explicit user authorization.
