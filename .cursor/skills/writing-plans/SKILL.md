---
name: writing-plans
description: Use when you have a spec or requirements for a multi-step task, before touching code
---

# Writing Plans

## Overview

Write comprehensive implementation plans assuming the engineer has zero context for our codebase and questionable taste. Document everything they need to know: which files to touch for each task, code, testing, docs they might need to check, how to test it. Give them the whole plan as bite-sized tasks. DRY. YAGNI. TDD.

Assume they are a skilled developer, but know almost nothing about our toolset or problem domain. Assume they don't know good test design very well.

**Announce at start:** "I'm using the writing-plans skill to create the implementation plan."

**Save plans to:** `docs/plans/YYYY-MM-DD-<feature-name>.md`

## Bite-Sized Task Granularity

**Each step is one action (2-5 minutes):**
- "Write the failing test" - step
- "Run it to make sure it fails" - step
- "Implement the minimal code to make the test pass" - step
- "Run the tests and make sure they pass" - step

## Plan Document Header

**Every plan MUST start with this header:**

```markdown
# [Feature Name] Implementation Plan

> **For Cursor:** Execute this plan using /executing-plans (separate session with checkpoints) or /subagent-driven-development (current session with subagents).

**Goal:** [One sentence describing what this builds]

**Architecture:** [2-3 sentences about approach]

**Tech Stack:** [Key technologies/libraries]

---
```

## Task Structure

```markdown
### Task N: [Component Name]

**Layer:** Domain / UseCase / Adapter (指明所属架构层)

**Files:**
- Create: `internal/{domain}/usecase/manager.go`
- Create: `internal/{domain}/usecase/manager_test.go`
- Modify: `internal/{domain}/usecase/interfaces.go:23-45`
- Modify: `internal/di/wire.go` (依赖注入)

**Step 1: Write the failing test**

```go
func TestManager_CreateEntity(t *testing.T) {
    // Arrange
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mock.NewMockRepository(ctrl)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    
    manager := NewManager(mockRepo)
    
    // Act
    err := manager.Create(context.Background(), &Entity{Name: "test"})
    
    // Assert
    assert.NoError(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v ./internal/{domain}/usecase/... -run TestManager_CreateEntity`
Expected: FAIL with "undefined: NewManager"

**Step 3: Write minimal implementation**

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

**Step 4: Run test to verify it passes**

Run: `go test -v ./internal/{domain}/usecase/... -run TestManager_CreateEntity`
Expected: PASS

**Step 5: Lint**

```bash
make fmt && make lint
git add internal/{domain}/usecase/
```
```

## Remember
- Exact file paths always
- Complete code in plan (not "add validation")
- Exact commands with expected output
- Reference relevant skills with @ syntax
- DRY, YAGNI, TDD

## Execution Handoff

After saving the plan, use `AskQuestion` to offer execution choice:

```
AskQuestion({
  title: "执行方式选择",
  questions: [{
    id: "execution_approach",
    prompt: "Plan saved to docs/plans/<filename>.md. How would you like to execute?",
    options: [
      { id: "subagent", label: "Subagent-Driven (this session) - Fresh subagent per task, fast iteration" },
      { id: "parallel", label: "Parallel Session (separate) - Batch execution with checkpoints" }
    ]
  }]
})
```

**If Subagent-Driven chosen:**
- **REQUIRED SUB-SKILL:** Use /subagent-driven-development
- Stay in this session
- Fresh subagent per task + code review

**If Parallel Session chosen:**
- **REQUIRED SUB-SKILL:** New session uses /executing-plans
