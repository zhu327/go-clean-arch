---
name: go
description: Orchestrates an end-to-end workflow from requirements to implementation, review, and simplification. Use when user says /go or asks for a full design→plan→execute→review→simplify flow.
disable-model-invocation: true
---

# go

Build features end-to-end with one automatic pipeline: design first, then implementation, then review, and cleanup.

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
2. Implement directly using TDD (read and follow `test-driven-development` skill)
3. Run validation: `go build ./...`, `go test ./...`, `go vet ./...`
4. Dispatch `code-reviewer` subagent for the change

**Full Pipeline** (multi-task features, new modules, cross-cutting changes):
Follow the full 5-step flow below.

---

## Flow (Full Pipeline)

1. Read and follow `brainstorming` skill — only in this step you must confirm requirements with the user.
2. After the user approves the design, continue automatically with `writing-plans` skill. If the change adds/modifies API endpoints, the plan MUST include E2E test tasks (E2E is a planned task, not a separate gate) and a completed Plan Coverage Checklist.
3. After plan creation, continue automatically with `subagent-driven-development` skill (plan preflight, wave-parallel execution using the dependency graph, per-task spec review, wave validation, final feature acceptance audit).
4. After implementation passes the acceptance audit, dispatch `code-reviewer` subagent and request complete fixes for identified issues.
5. After review/fixes, dispatch `code-simplifier` subagent once and report final result.

---

## Global Constraints

- Requirements are **locked after brainstorming approval**.
- Do not silently change scope, assumptions, or requirements in later stages.
- If deviation is necessary, **stop and report instead of proceeding**.

---

## Invocation behavior

- Do not ask for human confirmation between steps 2~5 (Full Pipeline) or during the Fast Track.
- If the user provides a task summary with `/go`, pass it into `brainstorming` as the starting context.
- If `brainstorming` identifies missing context, ask follow-up questions and wait for clarification before moving forward.
- The complexity gate is assessed by you after brainstorming — do not ask the user which track to use unless truly ambiguous.
- If there are implementation blockers during execution, stop and ask for resolution instead of guessing.

---

## Validation Hooks

- **Fast Track:** After implementation, run `go build ./...`, `go test ./...`, `go vet ./...` before dispatching code-reviewer.
- **Full Pipeline:** After `subagent-driven-development` completes each wave, run available validation if applicable:
  - Compilation / build checks (e.g. `go build ./...`)
  - Unit tests (e.g. `go test ./...`)
  - Static analysis / linting (e.g. `go vet ./...`)
- Before global code review, complete the SDD final feature acceptance audit.
- After `code-reviewer` fixes, re-run validation to ensure correctness.
- Only proceed to next stage if validation and the acceptance audit pass.

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
- Include and complete the Plan Coverage Checklist from `writing-plans`.
- Save the resulting plan and proceed immediately to execution once the plan is available.

---

## Step 3: Execute

- Read and follow `.cursor/skills/subagent-driven-development/SKILL.md`.
- Run SDD plan preflight before dispatching implementers.
- Parse the plan's dependency graph and compute execution waves.
- Dispatch independent tasks in parallel within each wave.
- Run spec-compliance review after each wave (architecture/code-quality is reviewed once, globally, in Step 4).
- Complete the final feature acceptance audit before Step 4.
- Do not ask for design or process approval again during execution.
- Respect locked requirements from brainstorming.

---

## Step 4: Review

- Dispatch the `code-reviewer` subagent for the whole change set (including E2E tests) only after the SDD acceptance audit passes.
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
- E2E tests for new API endpoints are planned as tasks in Step 2 and executed in Step 3 — **mandatory** for new endpoints, **recommended** for modified.
- The acceptance audit is a lightweight completeness check, not a second code-quality review.
