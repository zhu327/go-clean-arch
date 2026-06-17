---
name: writing-plans
description: Use after requirements or design are understood and before coding a multi-step change. Creates concrete implementation plans with vertical slices, dependency graphs, file lists, TDD steps, and validation commands.
---

# Writing Plans

## Overview

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Document everything they need to know: which files to touch for each task, code, testing, docs they might need to check, how to test it. Give them the whole plan as bite-sized tasks. DRY. YAGNI. TDD.

Assume they are a skilled developer, but know almost nothing about our toolset or problem domain. Assume they don't know good test design very well.

**Announce at start:** "I'm using the writing-plans skill to create the implementation plan."

**Save plans to:** `docs/plans/YYYY-MM-DD-<feature-name>.md` (use the `write` tool)

## When invoked by /go

When this skill is invoked as part of `/go`, the approved brainstorming design is locked input. Do not silently change approved requirements or reopen design choices.

Go mode overrides the normal planning prompts:

- Do not ask the plan-breakdown quiz in Step 3.
- Do not ask the execution-choice question in Execution Handoff.
- Write the full implementation plan directly from the approved design.
- Include a task dependency graph suitable for `/skill:subagent-driven-development`.
- Mark every task as AFK or HITL.
- If any HITL task needs a human decision not already resolved by the approved design, stop and report the blocker to the `go` coordinator instead of guessing.
- Include `**GO_EXECUTION_READY:** true` in the plan header only when the saved plan has no unresolved HITL decisions and can be executed autonomously.

After saving the plan in go mode, return the saved plan path to the `go` coordinator so it can continue automatically.

## Process

### 1. Gather Context

Before writing the plan, explore the relevant code directories and existing implementation patterns:

- Identify the domain modules involved and read their current structure (use `read`, `bash` with `fd`/`rg`)
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

Use the `question` tool to ask the user:

- Does the granularity feel right? (too coarse / too fine)
- Are the dependency relationships correct?
- Should any slices be merged or split further?
- Are the correct slices marked as HITL and AFK?

Outside `/go`, **STOP HERE.** Do not write the full plan or proceed to Step 4 until the user explicitly approves the breakdown and answers the quiz. 在这里停止。在用户明确批准拆分并回答问题之前，不要编写完整的计划或进入第 4 步。

In `/go`, skip this quiz only after there is an approved brainstorming design and no unresolved HITL decisions. If unresolved HITL decisions exist, stop and report the blocker to the `go` coordinator instead of proceeding.

### 4. Write the Full Plan

Expand each approved slice into the full Task Structure (see below) and save the plan document using the `write` tool.

## Task & Step Granularity

- **A Task** represents one Vertical Slice — a complete end-to-end feature path (typically 30-60 minutes).
- **A Step** represents a tiny TDD cycle within that slice (2-5 minutes each).

Because a Vertical Slice touches multiple layers, a single Task MUST contain multiple TDD steps organized by layer (e.g., Step 1-4 for Repository layer, Step 5-8 for UseCase layer, Step 9-12 for HTTP Handler layer).

## Plan Document Header

**Every plan MUST start with this header:**

```markdown
# [Feature Name] Implementation Plan

> **For Pi:** Execute this plan using /skill:executing-plans (separate session with checkpoints) or /skill:subagent-driven-development (current session with subagents).

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

**GO_EXECUTION_READY:** true/false (only required when invoked by `/go`)

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

#### Layer 1: Repository / DB

**Step 1: Write the failing Repository test**

```go
func TestRepo_Create(t *testing.T) {
    // Arrange: set up test DB / container
    repo := NewRepository(testDB)

    // Act
    err := repo.Create(context.Background(), &Entity{Name: "test"})

    // Assert
    assert.NoError(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./internal/{domain}/adapter/repository/... -run TestRepo_Create`
Expected: FAIL with "undefined: NewRepository"

**Step 3: Write minimal Repository implementation**

```go
// Repository 实体仓储
type Repository struct {
    db *gorm.DB
}

// NewRepository 创建仓储实例
func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// Create 持久化实体
func (r *Repository) Create(ctx context.Context, entity *Entity) error {
    return r.db.WithContext(ctx).Create(entity).Error
}
```

**Step 4: Run test to verify it passes**

Run: `go test -v ./internal/{domain}/adapter/repository/... -run TestRepo_Create`
Expected: PASS

---

#### Layer 2: UseCase

**Step 5: Write the failing UseCase test (mocking Repository)**

```go
func TestManager_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mock.NewMockRepository(ctrl)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

    manager := NewManager(mockRepo)

    err := manager.Create(context.Background(), &Entity{Name: "test"})
    assert.NoError(t, err)
}
```

**Step 6: Run test to verify it fails**

Run: `go test -v ./internal/{domain}/usecase/... -run TestManager_Create`
Expected: FAIL with "undefined: NewManager"

**Step 7: Write minimal UseCase implementation**

```go
// Manager 实体管理器
type Manager struct {
    repo Repository
}

// NewManager 创建管理器实例
func NewManager(repo Repository) *Manager {
    return &Manager{repo: repo}
}

// Create 创建实体
func (m *Manager) Create(ctx context.Context, entity *Entity) error {
    return m.repo.Create(ctx, entity)
}
```

**Step 8: Run test to verify it passes**

Run: `go test -v ./internal/{domain}/usecase/... -run TestManager_Create`
Expected: PASS

---

#### Layer 3: HTTP Handler & Route

**Step 9: Write the failing Handler test**

```go
func TestHandler_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockManager := mock.NewMockManager(ctrl)
    mockManager.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

    handler := NewHandler(mockManager)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request = httptest.NewRequest("POST", "/api/v1/entities", strings.NewReader(`{"name":"test"}`))

    handler.Create(c)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

**Step 10: Implement Handler, register Route, wire DI**

```go
// Handler HTTP 处理器
type Handler struct {
    manager Manager
}

func NewHandler(manager Manager) *Handler {
    return &Handler{manager: manager}
}

func (h *Handler) Create(c *gin.Context) {
    // parse, call manager, respond
}
```

**Step 11: Run all tests for this slice**

Run: `go test -v ./internal/{domain}/...`
Expected: ALL PASS

**Step 12: Lint & commit**

```bash
make fmt && make lint
git add internal/{domain}/
```
```

## Remember
- Exact file paths always
- Complete code in plan (not "add validation")
- Exact commands with expected output
- Reference relevant skills with /skill: syntax
- DRY, YAGNI, TDD
- If the plan adds/modifies API endpoints, identify the E2E requirement and include E2E authoring in the relevant task or as a separate vertical slice (use `/skill:e2e-testing`)
- Vertical slices, not horizontal layers — each task must be independently verifiable
- Declare dependencies explicitly — enable parallel execution where possible

## Execution Handoff

Outside `/go`, after saving the plan, use `question` to offer execution choice.

In `/go`, do not ask for execution choice. Return the saved plan path to the `go` coordinator so it can continue with `/skill:subagent-driven-development`.

Execution choice question:

```
question({
  questions: [{
    id: "execution_approach",
    prompt: "Plan saved to docs/plans/<filename>.md. How would you like to execute?",
    options: [
      { label: "Subagent-Driven (this session) - Fresh subagent per task, fast iteration", value: "subagent" },
      { label: "Executing Plans (separate session) - Batch execution with checkpoints", value: "executing_plans" }
    ]
  }]
})
```

**If Subagent-Driven chosen:**
- **REQUIRED SUB-SKILL:** Use /skill:subagent-driven-development
- Stay in this session
- Fresh subagent per task + code review
- Use dependency graph to parallelize independent AFK tasks

**If Executing Plans chosen:**
- **REQUIRED SUB-SKILL:** New session uses /skill:executing-plans
