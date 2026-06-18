---
description: End-to-end Go workflow — brainstorm, plan, execute, review, simplify
argument-hint: "<task description>"
---

# go

Build features end-to-end with one automatic pipeline: design first, then implementation, then review, and cleanup.

## Go Gate Checklist

Maintain this checklist visibly (via `todo` tool) throughout the session. Mark each gate as you complete it. Include final statuses in the final report.

### Full Pipeline Gates

- `GO-1: Brainstorming approved` — user confirmed requirements
- `GO-2: Design document saved` — written to `docs/plans/`
- `GO-3: Implementation plan saved` — written to `docs/plans/`
- `GO-4: Execution completed` — all tasks done, acceptance audit passed
- `GO-5: Validation passed` — `go build`, `go test`, `go vet` all green
- `GO-6: Review completed` — code-reviewer dispatched, fixes applied
- `GO-7: Simplification completed` — code-simplifier dispatched
- `GO-8: Final report delivered`

### Fast Track Gates (smaller subset)

- `GO-1: Brainstorming approved`
- `GO-2: Design document saved`
- `GO-4: TDD implementation completed`
- `GO-5: Validation passed`
- `GO-6: Review completed`
- `GO-7: Simplification completed`
- `GO-8: Final report delivered`

**Rules:**
- GO-5, GO-6, GO-7 may NOT be skipped — unavailable tools are blockers, not skip reasons
- "Task is small" is NOT a valid reason to skip review or simplification
- Final completion may only be claimed after all gates are `completed` or `blocked: <reason>`

## Forbidden Shortcuts

In go mode, never:

- Edit implementation code directly in the coordinator — always dispatch via `subagent`
- Skip validation when build/test/lint commands are available
- Skip review or simplification for any reason
- Silently change approved scope after brainstorming approval
- Continue past a blocker by guessing
- Claim completion before all required gates are done

## Complexity Gate

After brainstorming approval, assess the task scope before choosing a path:

| Criteria | Fast Track | Full Pipeline |
|----------|-----------|---------------|
| Files touched | 1–3 | 4+ |
| Layers affected | 1–2 | 3+ (cross-layer) |
| New domain/module | No | Yes |
| Architecture decision | No | Yes |
| Estimated tasks | 1–2 | 3+ |

**Fast Track** (small changes, bug fixes, single-slice features):
1. Skip `writing-plans` and `subagent-driven-development`
2. Implement directly using TDD (load and follow `/skill:test-driven-development`)
3. Run validation: `bash("go build ./... && go test ./... && go vet ./...")`
4. Dispatch `code-reviewer` agent via `subagent`
5. Dispatch `code-simplifier` agent via `subagent`

**Full Pipeline** (multi-task features, new modules, cross-cutting changes):
Follow the full 5-step flow below.

---

## Flow (Full Pipeline)

1. Load and follow `/skill:brainstorming` — only in this step you must confirm requirements with the user.
2. After the user approves the design, continue automatically with `/skill:writing-plans`. If the change adds/modifies API endpoints, the plan MUST include E2E test tasks (E2E is a planned task, not a separate gate) and a completed Plan Coverage Checklist.
3. After plan creation, continue automatically with `/skill:subagent-driven-development` (plan preflight, wave-parallel execution using the dependency graph, per-task spec review, wave validation, final feature acceptance audit).
4. After implementation passes the acceptance audit, dispatch `code-reviewer` agent via `subagent` and request complete fixes for identified issues.
5. After review/fixes, dispatch `code-simplifier` agent via `subagent` once and report final result.

---

## Global Constraints

- Requirements are **locked after brainstorming approval**.
- Do not silently change scope, assumptions, or requirements in later stages.
- If deviation is necessary, **stop and report instead of proceeding**.

---

## Invocation behavior

- Do not ask for human confirmation between steps 2~5 (Full Pipeline) or during the Fast Track.
- If the user provides a task summary with `/go`, pass it into `brainstorming` as the starting context.
- If `brainstorming` identifies missing context, ask follow-up questions (via `question`) and wait for clarification before moving forward.
- The complexity gate is assessed by you after brainstorming — do not ask the user which track to use unless truly ambiguous.
- If there are implementation blockers during execution, stop and ask for resolution instead of guessing.

---

## Validation Hooks

- **Fast Track:** After implementation, run `bash("go build ./... && go test ./... && go vet ./...")` before dispatching code-reviewer.
- **Full Pipeline:** After `subagent-driven-development` completes each wave, run available validation if applicable:
  - Compilation / build checks: `bash("go build ./...")`
  - Unit tests: `bash("go test ./...")`
  - Static analysis / linting: `bash("go vet ./...")`
- Before global code review, complete the SDD final feature acceptance audit.
- After `code-reviewer` fixes, re-run validation to ensure correctness.
- Only proceed to next stage if validation and the acceptance audit pass.

---

## Step 1: Brainstorming

> **Gate: GO-1 (Brainstorming approved) + GO-2 (Design saved)**

- Load and follow `/skill:brainstorming`.
- Explore project context (use `bash` for git/files, `read` for code) and clarify intent.
- Define clear requirements, constraints, and expected outputs.
- Output spec and wait for user approval before moving forward.
- Save design document → mark GO-2 complete.
- After user approves → mark GO-1 complete.
- **Assess Complexity Gate** → choose Fast Track or Full Pipeline.

---

## Step 2: Planning

> **Gate: GO-3 (Plan saved)**

- Load and follow `/skill:writing-plans`.
- Produce a concrete implementation plan based strictly on approved requirements.
- Do not introduce new requirements or alter scope.
- Include and complete the Plan Coverage Checklist from `writing-plans`.
- Save the resulting plan → mark GO-3 complete.
- Proceed immediately to execution.

---

## Step 3: Execute

> **Gate: GO-4 (Execution completed) + GO-5 (Validation passed)**

- Load and follow `/skill:subagent-driven-development`.
- Run SDD plan preflight before dispatching implementers.
- Parse the plan's dependency graph and compute execution waves.
- Dispatch independent tasks in parallel within each wave via the `subagent` tool's `tasks` payload.
- Run spec-compliance review after each wave (architecture/code-quality is reviewed once, globally, in Step 4).
- Complete the final feature acceptance audit.
- Run final validation: `bash("go build ./... && go test ./... && go vet ./...")` → mark GO-5 complete.
- Mark GO-4 complete after acceptance audit passes.
- Do not ask for design or process approval again during execution.
- Respect locked requirements from brainstorming.

---

## Step 4: Review

> **Gate: GO-6 (Review completed)**

- Dispatch the `code-reviewer` agent for the whole change set (including E2E tests) only after GO-4 and GO-5 pass.
- Use review scope: baseline commit and current commit.
- If any fixable Critical/Important issues are found:
  - Apply fixes directly via `subagent`
  - Re-run validation hooks (including `bash("make e2e")` if API tests were added)
- Repeat until:
  - No fixable Critical/Important issues remain, OR
  - No further safe fixes can be applied
- Mark GO-6 complete.

---

## Step 5: Simplify

> **Gate: GO-7 (Simplification completed)**

- Preconditions:
  - GO-6 complete (no outstanding Critical/Important issues)
  - Validation hooks pass

- Dispatch the `code-simplifier` agent against the final touched implementation files.
- Preserve existing behavior and external interfaces.
- Prefer minimal, readable refactors only.
- Do not introduce new abstractions unless clearly beneficial.
- If simplifier changes code, re-run validation to confirm behavior preserved.
- Mark GO-7 complete.
- Deliver final report → mark GO-8 complete.

---

## Notes

- `brainstorming` is the only step allowed to ask clarifying questions.
- Steps 2-5 run as an autonomous chain — do not ask for confirmation between them.
- `subagent-driven-development` uses the dependency graph from `writing-plans` to maximize parallel execution of independent tasks.
- E2E tests for new API endpoints are planned as tasks in Step 2 and executed in Step 3 — **mandatory** for new endpoints, **recommended** for modified.
- The acceptance audit is a lightweight completeness check, not a second code-quality review.

---

Task from user:

$@
