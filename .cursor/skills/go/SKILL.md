---
name: go
description: Orchestrates an end-to-end workflow from requirements to implementation, review, and simplification. Use when user says /go or asks for a full design→plan→execute→review→simplify flow.
---

# go

Build features end-to-end with one automatic pipeline: design first, then implementation, then review, and cleanup.

## Flow

1. Run `brainstorming` first, and only in this step you must confirm requirements with the user.
2. After the user approves the design, continue automatically with `writing-plans`.
3. After plan creation, continue automatically with `executing-plans`.
4. After implementation, run `code-reviewer` (subagent) and request complete fixes for identified issues.
5. After review/fixes, run `code-simplifier` (subagent) once and report final result.

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

- After `executing-plans`, run available validation if applicable:
  - Compilation / build checks (e.g. `go build ./...`)
  - Unit tests (e.g. `go test ./...`)
  - Static analysis / linting (e.g. `go vet ./...`)
- After `code-reviewer` fixes, re-run validation to ensure correctness.
- Only proceed to next stage if validation passes or no validation is available.

---

## Step 1: Brainstorming

- Announce: "I'm using the brainstorming skill to finalize requirements."
- Explore project context and clarify intent as specified by brainstorming rules.
- Define clear requirements, constraints, and expected outputs.
- Output spec and wait for user approval before moving forward.

---

## Step 2: Planning

- Announce: "I'm using the writing-plans skill to create the implementation plan."
- Produce a concrete implementation plan based strictly on approved requirements.
- Do not introduce new requirements or alter scope.
- Save the resulting plan and proceed immediately to execution once the plan is available.

---

## Step 3: Execute

- Announce: "I'm using the executing-plans skill to implement this plan."
- Execute plan tasks in the current session.
- Do not ask for design or process approval again during execution.
- Respect locked requirements from brainstorming.

---

## Step 4: Review

- Dispatch the `code-reviewer` subagent for the whole change set.
- Use review scope: baseline commit and current commit.
- If any fixable Critical/Important issues are found:
  - Apply fixes directly
  - Re-run validation hooks
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
- `writing-plans`, `executing-plans`, `code-reviewer` (subagent), and `code-simplifier` (subagent) are to be run as an autonomous chain in this order.
- If multiple implementation tracks conflict, prioritize the planning/execution chain (`executing-plans`) as the source of truth.
