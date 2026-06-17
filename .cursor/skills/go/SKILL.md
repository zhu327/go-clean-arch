---
name: go
description: Orchestrates an end-to-end workflow from requirements to implementation, review, and simplification. Use when user says /go or asks for a full design→plan→execute→review→simplify flow.
disable-model-invocation: true
---

# go

Build features end-to-end with one automatic pipeline: design first, then implementation, then review, and cleanup.

## Flow

1. Read and follow `brainstorming` skill — only in this step you must confirm requirements with the user.
2. After the user approves the design, continue automatically with `writing-plans` skill.
3. After plan creation, continue automatically with `subagent-driven-development` skill (wave-parallel execution using the dependency graph from the plan).
4. After implementation, run E2E tests if the change touches API endpoints (see E2E Gate below).
5. After E2E, dispatch `code-reviewer` subagent and request complete fixes for identified issues.
6. After review/fixes, dispatch `code-simplifier` subagent once and report final result.

---

## Global Constraints

- Requirements are **locked after brainstorming approval**.
- Do not silently change scope, assumptions, or requirements in later stages.
- If deviation is necessary, **stop and report instead of proceeding**.

---

## Invocation behavior

- Do not ask for human confirmation between steps 2~5.
- If the user provides a task summary with `/go`, pass it into `brainstorming` as the starting context.
- If `brainstorming` identifies missing context, ask follow-up questions and wait for clarification before moving to implementation planning.
- If there are implementation blockers during execution, stop and ask for resolution instead of guessing.

---

## Validation Hooks

- After `subagent-driven-development` completes each wave, run available validation if applicable:
  - Compilation / build checks (e.g. `go build ./...`)
  - Unit tests (e.g. `go test ./...`)
  - Static analysis / linting (e.g. `go vet ./...`)
- After `code-reviewer` fixes, re-run validation to ensure correctness.
- Only proceed to next stage if validation passes or no validation is available.

## E2E Gate

If the change adds or modifies API endpoints (handlers, routes, DTOs):

1. Read and follow `e2e-testing` skill to add E2E tests for new/changed endpoints.
2. Run `make e2e` — all tests must pass (existing + new).
3. If E2E fails, fix before proceeding to code review.

Skip this gate if the change is purely internal (domain logic, repository, no API surface change).

---

## Step 1: Brainstorming

- Read and follow `.cursor/skills/brainstorming/SKILL.md`.
- Explore project context and clarify intent as specified by brainstorming rules.
- Define clear requirements, constraints, and expected outputs.
- Output spec and wait for user approval before moving forward.

---

## Step 2: Planning

- Read and follow `.cursor/skills/writing-plans/SKILL.md`.
- Produce a concrete implementation plan based strictly on approved requirements.
- Do not introduce new requirements or alter scope.
- Save the resulting plan and proceed immediately to execution once the plan is available.

---

## Step 3: Execute

- Read and follow `.cursor/skills/subagent-driven-development/SKILL.md`.
- Parse the plan's dependency graph and compute execution waves.
- Dispatch independent tasks in parallel within each wave.
- Run two-stage review (spec compliance + code quality) after each wave.
- Do not ask for design or process approval again during execution.
- Respect locked requirements from brainstorming.

---

## Step 3.5: E2E Gate (conditional)

- If the implementation touches API endpoints, read and follow `.cursor/skills/e2e-testing/SKILL.md`.
- Run `make e2e` and confirm all tests pass.
- If no API surface was changed, skip this step.

---

## Step 4: Review

- Dispatch the `code-reviewer` subagent for the whole change set (including E2E tests).
- Use review scope: baseline commit and current commit.
- If any fixable Critical/Important issues are found:
  - Apply fixes directly
  - Re-run validation hooks (including `make e2e` if API tests were added)
- Repeat until:
  - No fixable Critical/Important issues remain, OR
  - No further safe fixes can be applied

---

## Step 5: Simplify

- Preconditions:
  - No outstanding Critical/Important issues
  - Validation hooks pass

- Dispatch the `code-simplifier` subagent against the final touched implementation files.
- Preserve existing behavior and external interfaces.
- Prefer minimal, readable refactors only.
- Do not introduce new abstractions unless clearly beneficial.

---

## Notes

- `brainstorming` is the only step allowed to ask clarifying questions.
- Steps 2-5 run as an autonomous chain — do not ask for confirmation between them.
- `subagent-driven-development` uses the dependency graph from `writing-plans` to maximize parallel execution of independent tasks.
- E2E tests are **mandatory** for new API endpoints and **recommended** for modified endpoints.
