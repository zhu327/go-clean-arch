---
name: writing-plans
description: Write implementation plans with vertical slices, dependency graphs, interface contracts, and test-case descriptions. Use when you have a spec or requirements for a multi-step task, before touching code. The plan defines the skeleton (contracts + verification) and leaves implementation to the TDD implementer. Produces plans executable by subagent-driven-development.
---

# Writing Plans

## Overview

Write implementation plans that define the **skeleton** of the work: which files to touch for each task, the interface contracts to create, the test cases that must pass, and how to verify the result. Give the implementer the whole plan as bite-sized vertical-slice tasks. DRY. YAGNI. TDD.

The plan describes **what** to build (contracts, file paths, dependencies) and **how it is verified** (acceptance criteria, test cases). It does NOT write the implementation code — that is the implementer's job, done test-first. Writing code in the plan wastes context and forces the implementer to regenerate (and often silently diverge from) the same code.

Assume the implementer is a skilled developer who knows almost nothing about our toolset or problem domain, and needs explicit contracts and test scenarios to work autonomously.

**Announce at start:** "I'm using the writing-plans skill to create the implementation plan."

**Save plans to:** `docs/plans/YYYY-MM-DD-<feature-name>.md`

## Process

### 1. Gather Context

Before writing the plan, explore the relevant code directories and existing implementation patterns:

- Identify the domain modules involved and read their current structure
- Note naming conventions, existing interfaces, and patterns in use
- Check for ADRs or design docs that constrain the approach
- Use the project's domain glossary vocabulary in task titles and descriptions

### 2. Draft Vertical Slices

Break the plan into **vertical slices (tracer bullets)**. Each task is a thin but COMPLETE path through all integration layers end-to-end, NOT a horizontal slice of one layer.

<vertical-slice-rules>
- Each slice delivers a narrow but COMPLETE path through every layer it touches (Domain → UseCase → Repository/Gateway → Handler → Route → Tests)
- A completed slice is demoable or verifiable on its own
- Prefer many thin slices over few thick ones
- A slice like "implement all domain entities" is WRONG — instead "implement Create Certificate end-to-end" is correct
</vertical-slice-rules>

### 3. Quiz the User

Present the proposed breakdown as a numbered list before writing the full plan. For each slice, show:

- **Title**: short descriptive name
- **Type**: HITL (needs human decision) / AFK (agent can execute autonomously)
- **Blocked by**: which other slices (if any) must complete first
- **Layers touched**: Domain / UseCase / Adapter / etc.

Ask the user:

- Does the granularity feel right? (too coarse / too fine)
- Are the dependency relationships correct?
- Should any slices be merged or split further?
- Are the correct slices marked as HITL and AFK?

**STOP HERE.** Do not write the full plan or proceed to Step 4 until the user explicitly approves the breakdown and answers the quiz. 在这里停止。在用户明确批准拆分并回答问题之前，不要编写完整的计划或进入第 4 步。

### 4. Write the Full Plan

Expand each approved slice into the full Task Structure (see below), add the Plan Coverage Checklist, and save the plan document.

### 5. Plan Coverage Check

Before handing off to execution, explicitly verify the plan covers the approved requirements and has the guardrails needed for autonomous execution. If any checklist item fails, revise the plan before proceeding — do not hand off a known-incomplete plan.

## Task Granularity

- **A Task** represents one Vertical Slice — a complete end-to-end feature path (typically 30-60 minutes).

Because a Vertical Slice touches multiple layers, a single Task MUST list the interface contracts and test cases for each layer it touches (e.g., Repository, UseCase, Handler). The plan defines WHAT to build (contracts) and HOW it is verified (acceptance criteria + test cases); it does NOT write the implementation code. The implementer applies TDD per test case (red-green-refactor) — that discipline lives in the `implementer-prompt`, so do not repeat it as numbered steps in every task.

## Plan Document Header

**Every plan MUST start with this header:**

```markdown
# [Feature Name] Implementation Plan

> **For Cursor:** Execute this plan using `subagent-driven-development` (wave-parallel subagents in the current session).

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

## Task Dependency Graph

Tasks marked ✅ AFK can be executed by agents autonomously.
Tasks marked 🙋 HITL require human decision before proceeding.

\```
Task 1 (AFK) ──┐
                ├── Task 4 (AFK)
Task 2 (AFK) ──┘
Task 3 (HITL) ───── Task 5 (AFK)
\```

| Task | Type | Blocked by | Parallelizable with |
|------|------|------------|---------------------|
| 1    | AFK  | None       | 2, 3                |
| 2    | AFK  | None       | 1, 3                |
| ...  | ...  | ...        | ...                 |

---
```

## Task Structure

```markdown
### Task N: [Vertical Slice Name] (e.g., "Create Certificate End-to-End")

**Type:** AFK / HITL
**Blocked by:** Task X, Task Y / None - can start immediately
**Layers touched:** Domain, UseCase, Repository, Handler

**Goal:** [One or two sentences: what this slice delivers end-to-end and why]

**Acceptance Criteria:**
- [ ] Criterion 1 (functional behavior from user perspective)
- [ ] Criterion 2 (edge case or boundary condition)
- [ ] Criterion 3 (non-functional requirement if any)

**Files:**
- Create: `internal/{domain}/domain/entity.go`
- Create: `internal/{domain}/adapter/repository/repo.go`
- Create: `internal/{domain}/adapter/repository/repo_test.go`
- Create: `internal/{domain}/usecase/manager.go`
- Create: `internal/{domain}/usecase/manager_test.go`
- Create: `internal/{domain}/adapter/delivery/http/handler/handler.go`
- Modify: `internal/{domain}/adapter/delivery/http/router/router.go`
- Modify: `internal/di/wire.go`

---

#### Interface Contracts

Define the contracts the implementer must create — interfaces, struct shapes, and function signatures. These are the skeleton; the implementer writes the bodies via TDD. List only what this task introduces or changes. If a contract already exists in another module, reference its file path instead of redefining it.

\```go
// domain/entity.go
type Entity struct {
    ID   string
    Name string
    // ...fields with intent, not boilerplate
}

// usecase/interfaces.go
type Repository interface {
    Create(ctx context.Context, entity *domain.Entity) error
}

type Manager interface {
    Create(ctx context.Context, entity *domain.Entity) error
}

// adapter/repository/repo.go
func NewRepository(db *gorm.DB) *Repository

// adapter/delivery/http/handler/handler.go
func NewHandler(manager Manager) *Handler
func (h *Handler) Create(c *gin.Context)
\```

> Signatures only. No method bodies, no test code. These contracts are what downstream tasks (and the spec reviewer) will check against.

#### Test Cases to Cover

Describe the scenarios the tests must verify — behavior and expected outcomes, not full test code. The implementer writes each test BEFORE its implementation (TDD). Organize by layer.

**Repository layer:**
- Create succeeds and persists the entity (verify via re-query or a test DB/container)
- Create returns the DB error transparently on failure

**UseCase layer:**
- Create delegates to the Repository and returns its error
- Create validates input (e.g., rejects empty Name) before calling the Repository

**HTTP Handler layer:**
- POST returns 200 with the created resource on valid input
- POST returns 400 on missing/invalid required fields
- Handler delegates to Manager and maps errors to the right response codes

#### Layer Guidance

Brief notes on each layer's responsibility in this slice, so the implementer understands intent without being handed the code:

- **Domain:** define the `Entity` and any invariants
- **Repository:** persist via GORM; follow existing repo patterns in this module
- **UseCase:** orchestrate, validate input, own business rules; depend on the `Repository` interface (mocked in tests)
- **Handler:** bind request, call Manager, format response; register the route in the router; wire into DI

---

#### Validation

Run after the slice is implemented:

\```bash
go build ./...
go test ./internal/{domain}/...
go vet ./{domain}/...
\```

If this slice adds/modifies API endpoints, an E2E test task (see below) must also pass: `make e2e`.
```

### E2E Test Tasks

If the feature adds or modifies API endpoints, the plan MUST include one or more E2E test tasks (do NOT treat E2E as a separate gate that runs after the whole feature — it is a task in the plan, executed by the same subagent-driven-development flow). Each E2E task:

- Is typed AFK and blocked by the endpoint task(s) it covers
- Lists the endpoints + scenarios to cover (use the `e2e-testing` skill for the methodology)
- Has acceptance criteria like: `make e2e` passes (existing + new tests)

## Plan Coverage Checklist

Every plan MUST include this checklist before the Execution Handoff. Fill it in based on the approved brainstorming requirements/design and the final task list:

```markdown
## Plan Coverage Checklist

- [ ] Every approved requirement maps to at least one task
- [ ] Every task has clear acceptance criteria
- [ ] Every task lists behavior-focused test cases
- [ ] Every task lists exact Create/Modify file paths
- [ ] New or modified API endpoints have E2E test task(s)
- [ ] The dependency graph has no cycles
- [ ] Parallelizable tasks do not modify the same files
- [ ] No task is purely horizontal unless it is unavoidable infrastructure
- [ ] Known assumptions or deviations from the approved design are documented
```

If a checklist item does not apply, mark it `N/A` with a short reason instead of silently omitting it.

## Remember
- Exact file paths always
- Interface contracts + test-case descriptions in the plan, NOT implementation code — define the skeleton, not the body; leave the code to the TDD implementer
- A one-line minimal code example is allowed only when a contract cannot be made clear in prose — never full method bodies
- Exact validation commands where relevant
- Reference relevant skills with @ syntax
- DRY, YAGNI, TDD (enforced by the implementer, not re-written as numbered steps per task)
- If the plan adds/modifies API endpoints, it MUST include E2E test tasks (use `e2e-testing` skill) — E2E is a planned task, not a separate post-implementation gate
- Include and complete the Plan Coverage Checklist before handoff
- Vertical slices, not horizontal layers — each task must be independently verifiable
- Declare dependencies explicitly — enable parallel execution where possible

## Execution Handoff

After saving the plan, hand off to `subagent-driven-development` for execution in the current session:

- Read and follow `.cursor/skills/subagent-driven-development/SKILL.md`
- Fresh subagent per task, wave-parallel execution using the dependency graph
- Spec-compliance review gate per task; global architecture/quality review at the end

Within the `/go` pipeline this handoff is automatic (no confirmation). When running `writing-plans` standalone, confirm with the user before starting execution:

```
AskQuestion({
  title: "开始执行",
  questions: [{
    id: "start_execution",
    prompt: "Plan saved to docs/plans/<filename>.md. Start subagent-driven execution now?",
    options: [
      { id: "yes", label: "Yes - execute now with subagent-driven-development" },
      { id: "no", label: "Not yet - I'll review the plan first" }
    ]
  }]
})
```
