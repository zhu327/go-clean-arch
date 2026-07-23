---
name: subagent-driven-development
description: Use when executing a written implementation plan in the current session. Dispatches implementation subagents by dependency-safe waves, verifies task compliance, and runs project-derived validation.
---

# Subagent-Driven Development

Execute an approved plan with fresh implementer subagents. The dependency graph controls scheduling: independent tasks may run in parallel only when their file lists do not conflict.

## Preflight

1. Read the plan and extract tasks, dependencies, files, acceptance criteria, tests, and validation commands.
2. Confirm every task has exact files, test cases, and acceptance criteria.
3. Confirm the coverage checklist is complete, dependencies have no cycles, and concurrent tasks do not create or modify the same files.
4. Identify the repository’s actual build, test, lint/typecheck, integration, and E2E commands. Do not assume any language-specific command.
5. If public behavior changes, confirm that the plan includes appropriate integration/E2E coverage when the repository supports it or requirements demand it.

If any preflight condition fails, stop and return the concrete gap to planning; do not guess.

## Execute by Waves

For each dependency wave:

1. Resolve HITL tasks through the user before dispatch.
2. Dispatch up to four non-conflicting AFK implementers in parallel. Give each the full task text, its allowed files, context, and actual contracts from completed dependencies.
3. Wait for all implementers in the wave.
4. Dispatch independent spec-compliance reviewers for every completed task.
5. If a reviewer finds an issue, have the implementer fix it and repeat review until it passes.
6. Run the validation commands identified in the plan and project. If integrated work fails, dispatch one integration fixer with the exact output and changed-file list, then revalidate. After three unsuccessful attempts, stop and escalate.
7. Mark the wave complete only after review and validation pass.

## Artifact Passing

For a task that depends on earlier work, read the real generated artifacts—not summaries—and provide the exact relevant contract, schema, configuration shape, or exported interface to the next implementer. If actual artifacts disagree with the approved plan, stop and report the mismatch.

## Acceptance Audit

After all waves, verify:

- all planned tasks are complete;
- all project-derived validations pass;
- required integration/E2E suites pass;
- the complete feature meets the approved goal;
- cross-task contracts are consistent;
- deviations are documented and approved.

Then dispatch the `code-reviewer` agent once for a global architecture, security, and maintainability review. The per-task reviewers assess specification compliance; the global reviewer assesses code quality.

## Safety Rules

Never:

- start implementation on the main branch without explicit user consent;
- run parallel implementers that modify the same file;
- ask an implementer to rediscover the plan instead of providing task text;
- skip a failed spec-review finding, validation failure, or required acceptance audit;
- replace independent review with implementer self-review;
- proceed to a dependent wave before its prerequisites pass.

Use `implementer-prompt.md` and `spec-reviewer-prompt.md` as templates.
