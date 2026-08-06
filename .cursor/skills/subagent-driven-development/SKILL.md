---
name: subagent-driven-development
description: Execute an approved go-clean-arch plan with parallel implementation subagents when independent slices have stable contracts and non-overlapping write sets. Do not use when coordination cost exceeds plausible parallel benefit.
---

# Parallel Plan Execution

Use implementer subagents only for genuinely independent planned work. The coordinator may implement directly when tasks are coupled, small, or depend on evolving context.

## Entry criteria

Confirm that:

- every slice has a goal, acceptance criteria, allowed files, and focused validation;
- dependencies are acyclic and prerequisite contracts are stable;
- concurrent slices do not modify the same files, shared interfaces, Wire sets, generated mocks, Swagger outputs, or other shared artifacts;
- relevant `.cursor/rules/` and repository commands have been identified;
- parallel execution likely saves more time than context transfer and integration.

Otherwise execute sequentially or return the concrete planning gap. Do not invent contracts to enable parallelism.

## Execute by waves

For each dependency wave:

1. Resolve blocking user decisions.
2. Dispatch the smallest useful number of non-conflicting implementers.
3. Provide the complete task, allowed files, relevant layer rules, and exact contracts from prerequisites.
4. Inspect actual changes, not only reports.
5. Run focused tests for each slice and an appropriate wave integration check.
6. Use the narrowest owner for integration fixes; stop after repeated failures rather than looping indefinitely.

Use `implementer-prompt.md` as an adaptable template. `examples.md` illustrates dependency waves and artifact passing, not mandatory ceremony.

## Independent compliance review

Use a separate spec reviewer for high-risk slices, trust/ownership boundaries, material public-contract changes, or behavior that tests cannot verify convincingly. Otherwise implementer self-check, coordinator inspection, and passing tests are sufficient.

## Artifact passing

Read generated artifacts and provide exact relevant declarations, schemas, configuration, or contracts to downstream work. If actual output differs from the approved plan, resolve it before dispatching dependent tasks.

## Final integration

After all waves, confirm acceptance criteria and cross-slice contracts; run applicable integrated validation once; inspect the combined diff for scope, generated artifacts, and overlap; report deviations and residual risk.

Changeset-level code review belongs to the invoking workflow; do not dispatch a duplicate global reviewer here.

## Safety rules

Never run parallel writers on overlapping files, allow silent write-set expansion, or claim validation not run. Follow repository and user authorization for destructive operations and branch policy.
