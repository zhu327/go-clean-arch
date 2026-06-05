---
description: End-to-end Go workflow — brainstorm, plan, execute, review, simplify
argument-hint: "<task description>"
---

# go

Build features end-to-end with one automatic pipeline: design first, then implementation, then review, and cleanup.

## Go Mode Coordinator Contract

While `/go` is active, the main agent is the workflow coordinator, not the implementation worker.

- Coordinate the pipeline, maintain gate status, dispatch child skills/agents, validate results, and report blockers.
- Do not directly implement production changes during the execution phase.
- Implementation must happen through `subagent` calls via `subagent-driven-development`, even when the approved plan has only one task.
- Maintain a visible Go Gate Checklist (todo tool or equivalent visible checklist) from the start of go mode through the final response.
- Final completion may only be claimed after all required gates are `completed`, `blocked: <reason>`, or `skipped: <valid reason>`.
- GO-4, GO-6, and GO-7 may not be skipped; unavailable agents/tools are blockers, not skipped gates.
- GO-5 may not be skipped when validation commands are available; unavailable validation tooling is a blocker unless no applicable validation exists.
- Skipping review or simplification because the task is small is not a valid reason.

## Required Go Gates

Maintain these gates visibly and include their final statuses in the final report:

- `GO-1: Brainstorming design approved` — status: `completed`, `blocked: <reason>`, or `skipped: <valid reason>`
- `GO-2: Design document saved` — status: `completed`, `blocked: <reason>`, or `skipped: <valid reason>`
- `GO-3: Implementation plan saved` — status: `completed`, `blocked: <reason>`, or `skipped: <valid reason>`
- `GO-4: Subagent execution completed` — status: `completed` or `blocked: <reason>`; may not be skipped
- `GO-5: Validation completed` — status: `completed`, `blocked: <reason>`, or `skipped: <valid reason>` only when no applicable validation commands are available
- `GO-6: Whole-change review completed and required fixes handled` — status: `completed` or `blocked: <reason>`; may not be skipped
- `GO-7: Simplifier completed` — status: `completed` or `blocked: <reason>`; may not be skipped
- `GO-8: Final report delivered` — status: `completed`, `blocked: <reason>`, or `skipped: <valid reason>`

## Forbidden Shortcuts

In go mode, never:

- Edit implementation code before the planning gate completes.
- Execute planned implementation directly in the coordinator/main agent.
- Replace `subagent-driven-development` with manual sequential execution.
- Skip subagent execution, even for small or single-task changes.
- Skip validation when build, test, lint, or E2E commands are available.
- Skip whole-change review.
- Skip simplification.
- Treat whole-change review or simplification as optional polish.
- Silently change approved scope, assumptions, requirements, or external behavior after brainstorming approval.
- Continue past a blocker by guessing.

## Child-Skill Overrides in Go Mode

When child-skill instructions conflict with this coordinator contract, these go-mode overrides take precedence:

### Brainstorming override

- Use `question` for requirement clarification and incremental design validation.
- Save the approved design document.
- After the user approves the design, do not ask whether they are ready to set up for implementation.
- Return control to the go coordinator so it can automatically invoke `writing-plans`.

### Writing-plans override

- Use the approved brainstorming design as locked input.
- Do not ask the plan-breakdown quiz.
- Do not ask the execution-choice question.
- Write the full implementation plan directly.
- Include a dependency graph suitable for `subagent-driven-development`.
- Mark tasks as AFK/HITL.
- If a HITL task needs a human decision not covered by the approved design, stop and report the blocker to the `/go` coordinator.
- Include `**GO_EXECUTION_READY:** true` in the plan header only when the plan can be executed autonomously.

### Subagent-driven-development override

- At least one `subagent` call is mandatory.
- If the plan has one task, dispatch one implementation subagent rather than implementing in the coordinator.
- If independent tasks exist, use dependency-aware parallel waves.
- Do not ask new HITL questions during execution; unresolved human decisions block GO-4.
- Return a structured execution report containing tasks run, agents used, files changed, validation commands and results, blockers, and deviations.

### Review override

- `code-reviewer` produces findings only; it must not ask the user how to proceed while `/go` is active.
- The `/go` coordinator owns the fix loop: automatically dispatch implementation/fix subagents for safe P0/P1 fixes, and block on ambiguous or unsafe fixes.

## Flow

1. Run `brainstorming` first, and only in this step you must confirm requirements with the user.
2. After the user approves the design, continue automatically with `writing-plans`.
3. After plan creation, continue automatically with `subagent-driven-development` (wave-parallel execution using the dependency graph from the plan).
4. After implementation, run E2E tests if the change touches API endpoints (see E2E Gate below).
5. After E2E, run `code-reviewer` agent (via `subagent` tool) and request complete fixes for identified issues.
6. After review/fixes, run `code-simplifier` agent (via `subagent` tool) once and report final result.

---

## Global Constraints

- Requirements are **locked after brainstorming approval**.
- Do not silently change scope, assumptions, or requirements in later stages.
- If deviation is necessary, **stop and report instead of proceeding**.

---

## Invocation behavior

- Do not ask for human confirmation between steps 2~5.
- If the user provides a task summary with `/go`, pass it into `brainstorming` as the starting context.
- If `brainstorming` identifies missing context, ask follow-up questions (via `question`) and wait for clarification before moving to implementation planning.
- If there are implementation blockers during execution, stop and ask for resolution instead of guessing.

---

## Validation Hooks

- After `subagent-driven-development` completes each wave, run available validation if applicable:
  - Compilation / build checks: `bash("go build ./...")`
  - Unit tests: `bash("go test ./...")`
  - Static analysis / linting: `bash("go vet ./...")`
- After `code-reviewer` fixes, re-run validation to ensure correctness.
- After `code-simplifier` changes code, re-run validation to ensure behavior is preserved.
- Only proceed to next stage if validation passes or no validation is available.

## E2E Gate

If the change adds or modifies API endpoints (handlers, routes, DTOs):

1. Confirm the implementation plan included E2E coverage; if it did not, use `/skill:e2e-testing` to add the missing tests before validation.
2. Run `bash("make e2e")` — all tests must pass (existing + new).
3. If E2E fails, fix before proceeding to code review.

Ownership: `writing-plans` identifies E2E needs, `subagent-driven-development` creates or updates tests during implementation, and this gate validates the result. Skip this gate only if the change is purely internal (domain logic, repository, no API surface change).

---

## Step 1: Brainstorming

- Announce: "I'm using the brainstorming skill to finalize requirements."
- Load and follow `/skill:brainstorming`
- Explore project context (use `read`, `bash` for git/files) and clarify intent.
- Define clear requirements, constraints, and expected outputs.
- Output spec and wait for user approval before moving forward.

---

## Step 2: Planning

- Announce: "I'm using the writing-plans skill to create the implementation plan."
- Load and follow `/skill:writing-plans`
- Produce a concrete implementation plan based strictly on approved requirements.
- Do not introduce new requirements or alter scope.
- Save the resulting plan (use `write` tool) and proceed immediately to execution.

---

## Step 3: Execute

- Announce: "I'm using the subagent-driven-development skill to implement this plan."
- Load and follow `/skill:subagent-driven-development`
- Parse the plan's dependency graph and compute execution waves.
- Dispatch independent tasks in parallel within each wave via the `subagent` tool's `tasks` payload.
- Collect each wave's execution results, deviations, and validation output for the later whole-change review gate; do not introduce a separate final review owner during execution.
- Do not ask for design or process approval again during execution.
- Respect locked requirements from brainstorming.

---

## Step 3.5: E2E Gate (conditional)

- If the implementation touches API endpoints, announce: "Running E2E gate — validating endpoint coverage."
- If E2E tests are missing or incomplete, load and follow `/skill:e2e-testing` to add them.
- Run `bash("make e2e")` and confirm all tests pass.
- If no API surface was changed, skip this step.

---

## Step 4: Review

- Dispatch the `code-reviewer` agent for the whole change set:
  ```json
  {
    "agent": "code-reviewer",
    "task": "Review all changes in this implementation..."
  }
  ```
- Use review scope: baseline commit and current commit.
- If the `code-reviewer` agent is unavailable, stop and report the blocker instead of skipping review.
- If any safe fixable P0/P1 (P0 Critical / P1 High) issues are found:
  - Dispatch a suitable implementation/fix subagent to apply the fixes.
  - If no suitable fix subagent is available, stop and report the blocker instead of editing directly.
  - Re-run validation hooks (including `bash("make e2e")` if API tests were added)
- Repeat until:
  - No safe fixable P0/P1 (P0 Critical / P1 High) issues remain, OR
  - No further safe fixes can be handled through a suitable fix subagent
- Unresolved real P0/P1 (P0 Critical / P1 High) issues block simplification and final completion until reported and resolved.
- This is a completion gate, not optional polish.

---

## Step 5: Simplify

- Preconditions:
  - No outstanding P0 Critical / P1 High issues
  - Validation hooks pass

- Dispatch the `code-simplifier` agent:
  ```json
  {
    "agent": "code-simplifier",
    "task": "Simplify the implementation files touched by this change..."
  }
  ```
- If the `code-simplifier` agent is unavailable, stop and report the blocker instead of skipping simplification.
- Preserve existing behavior and external interfaces.
- Prefer deletion and simplification over new abstractions.
- Prefer minimal, readable refactors only.
- Do not introduce new abstractions unless clearly beneficial.
- If the simplifier changes code, re-run validation hooks before final reporting.
- This is a completion gate, not optional polish.

---

## Notes

- `brainstorming` is the only step allowed to ask requirement/design clarifying questions via `question`; unresolved later human decisions become blockers, not new scope discussions.
- `writing-plans`, `subagent-driven-development`, `e2e-testing` (conditional), `code-reviewer` (agent), and `code-simplifier` (agent) are to be run as an autonomous chain in this order.
- `subagent-driven-development` uses the dependency graph from `writing-plans` to maximize parallel execution of independent tasks.
- E2E tests are **mandatory** for new API endpoints and **recommended** for modified endpoints.

---

Task from user:

$@
