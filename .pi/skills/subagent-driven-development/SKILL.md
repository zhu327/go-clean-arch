---
name: subagent-driven-development
description: Use when executing a written implementation plan in the current session. Dispatches implementation subagents per task with dependency-aware wave parallelism, file-conflict checks, spec review, code-quality review, and validation.
disable-model-invocation: true
---

# Subagent-Driven Development

Execute plan by dispatching subagents per task with **dependency-aware parallel scheduling**. Tasks without mutual dependencies run concurrently in waves. Each task gets two-stage review (spec compliance then code quality).

**Core principle:** Dependency graph drives parallelism. Fresh subagent per task via the `subagent` tool. Two-stage review gates. Maximum concurrency within safety bounds.

**Tool used:** `subagent` — use top-level `agent` + `task` for one-off dispatch, `tasks` array for parallel dispatch, and `chain` for sequential chains.

## When to Use

```
Have implementation plan?
  ├─ yes → Plan has dependency graph?
  │         ├─ yes → Stay in this session?
  │         │         ├─ yes → subagent-driven-development (this skill)
  │         │         └─ no  → executing-plans (separate session)
  │         └─ no  → treat all as sequential
  └─ no → Manual execution or brainstorm first
```

**vs. Executing Plans:**
- Same session (no context switch)
- Fresh subagent per task (no context pollution)
- **Parallel waves** — independent tasks run concurrently via `subagent` `tasks` payloads
- Two-stage review after each task: spec compliance first, then code quality
- Faster iteration (no human-in-loop between AFK tasks)

**vs. Dispatching Parallel Agents:**
- This skill: **wave-parallel** execution of **planned** tasks with review gates
- Parallel Agents: **concurrent** execution of independent **ad-hoc** problems (bugs, investigations)
- Use /skill:dispatching-parallel-agents when you have N independent bugs/issues, NOT a plan to execute

## When Invoked by /go

When this skill is invoked as part of `/go`, the main agent is the workflow coordinator and **must not implement planned changes directly**.

Go mode requirements:
- At least one `subagent` call is mandatory before execution can be considered complete.
- If the implementation plan has exactly one task, dispatch one implementation subagent for that task instead of implementing it in the coordinator.
- If independent tasks exist, use the normal dependency-aware parallel waves and file-conflict checks.
- Sequential fallback is still subagent-driven: dispatch one task at a time via the `subagent` tool with top-level `agent` + `task`; never replace subagent execution with coordinator implementation.
- `/go` mode must not ask new HITL questions during execution. If a task still needs an unresolved human decision, stop and report GO-4 as blocked.
- If no suitable implementation subagent can be dispatched, **STOP** and report the blocker to the `/go` coordinator.
- After implementation and validation, return a structured execution report to the `go` coordinator; the `go` coordinator owns the separate whole-change review gate.

Go mode execution report format:

```markdown
## Subagent Execution Report

**Tasks run:**
- <task name/id> — <completed/blocked/deviated: reason>

**Agents used:**
- <agent name or default> — <task/review/fix role>

**Files changed:**
- `<path>` — <summary or task responsible>

**Validation:**
- `<command>` — <passed/failed/skipped: reason>

**Blockers / deviations:**
- None / <details>
```

## The Process

```
1. Parse plan: extract tasks, dependency graph, file lists
2. Compute execution waves via topological sort
3. Verify file conflicts within each wave
4. Create todo items with all tasks

Per Wave:
  a. Handle HITL tasks first (outside `/go`: ask user via question; inside `/go`: block if unresolved)
  b. Dispatch AFK implementers in parallel via the `subagent` tool's `tasks` payload (max 4 per batch)
  c. Collect results as subagents complete
  d. Dispatch ALL spec reviewers in parallel via the `subagent` tool's `tasks` payload
  e. Fix spec issues (implementer fixes → re-review)
  f. Dispatch ALL code quality reviewers in parallel via the `subagent` tool's `tasks` payload
  g. Fix quality issues (implementer fixes → re-review)
  h. Run validation (go build, go test, go vet) via bash
  i. If validation fails → dispatch Integration Fixer subagent
  j. Mark all wave tasks complete in todo

After all waves:
  Non-go mode: dispatch final code-reviewer agent for entire implementation
  /go mode: return execution report; go coordinator owns whole-change review
```

## Phase 1: Parse Plan & Build Schedule

1. Read the plan file once (use `read` tool)
2. Extract the **Task Dependency Graph** (the table with "Blocked by" and "Parallelizable with" columns)
3. Extract the **Files** section of each task (Create/Modify lists)
4. Compute execution waves via topological sort:
   - **Wave 0:** tasks with no dependencies (can start immediately)
   - **Wave N:** tasks whose ALL dependencies are in waves < N
5. Verify file safety within each wave (see File Conflict Safety below)
6. Separate HITL from AFK tasks within each wave
7. Create `todo` items with all tasks, annotated with wave number

### Wave Computation Example

Given dependency graph:
```
Task 1 (AFK) ──┐
                ├── Task 4 (AFK)
Task 2 (AFK) ──┘
Task 3 (HITL) ───── Task 5 (AFK)
```

| Task | Blocked by | Wave |
|------|------------|------|
| 1    | None       | 0    |
| 2    | None       | 0    |
| 3    | None       | 0    |
| 4    | 1, 2       | 1    |
| 5    | 3          | 1    |

→ Wave 0: [Task 1, Task 2, Task 3] — all parallel
→ Wave 1: [Task 4, Task 5] — all parallel (after Wave 0)

## Phase 2: Execute Waves

For each wave, in order:

### Step 1: Handle HITL Tasks

HITL tasks need human decisions before agent can proceed:
1. Present the HITL task to the user with the decision needed (use `question` tool)
2. Wait for user response
3. Dispatch implementer subagent with the decision as additional context
4. Run spec review + code quality review (sequential per HITL task)
5. Mark complete in `todo`

### Step 2: Dispatch AFK Implementers in Parallel

**Concurrency Limit:** Dispatch at most **4 AFK implementers** simultaneously. If a wave (or sub-wave) has more than 4 AFK tasks, split into batches within the wave (e.g., 7 tasks = batch of 4 + batch of 3). This prevents API rate limits and resource exhaustion. Batches within the same wave run sequentially, but tasks within each batch run in parallel.

**Dispatch tasks using the `subagent` tool with `tasks` array (max 4):**

```json
{
  "tasks": [
    { "task": "Implement Task 1: Certificate CRUD\n\n[full task text]" },
    { "task": "Implement Task 2: SubDomain DNS\n\n[full task text]" },
    { "task": "Implement Task 3: WorkOrder Integration\n\n[full task text]" }
  ]
}
```
All three run concurrently. Add `"agent": "worker"` only when a suitable implementation agent exists.

**Note:** If no `worker` implementation agent is configured, omit the `agent` field to use the default subagent. Do not use `code-reviewer` for implementation; it is review-only.

Each subagent gets:
- Full task text from plan (never make subagent read plan file)
- Scene-setting context (where this fits, what was built in prior waves)
- Prior wave outputs if relevant (e.g., "Task 1 created `CertificateRepository` interface at `internal/certificate/usecase/interfaces.go` — you depend on it")

### Step 3: Collect Results

As subagents complete, collect their reports. Wait for ALL implementers in the wave to finish before proceeding to reviews. A parallel `subagent` call with a `tasks` payload returns all results when done.

### Step 4: Dispatch Spec Reviewers in Parallel

**Dispatch ALL spec reviewers simultaneously** for all completed tasks in the wave:

```json
{
  "tasks": [
    { "agent": "code-reviewer", "task": "Review spec compliance for Task 1\n\n[requirements + implementer report]" },
    { "agent": "code-reviewer", "task": "Review spec compliance for Task 2\n\n[requirements + implementer report]" },
    { "agent": "code-reviewer", "task": "Review spec compliance for Task 3\n\n[requirements + implementer report]" }
  ]
}
```

### Step 5: Fix Spec Issues

For each task where spec review found issues:
1. Resume the original implementer subagent to fix issues (dispatch a new `subagent` with fix instructions and context of what was built)
2. Re-dispatch spec reviewer
3. Repeat until ✅

Fix loops are **per-task and serialized** (implementer must fix before re-review).
Tasks that passed spec review can proceed to code quality review immediately.

### Step 6: Dispatch Code Quality Reviewers in Parallel

Same pattern: dispatch all code quality reviewers simultaneously for spec-approved tasks via the `subagent` tool's `tasks` payload.

### Step 7: Fix Quality Issues

Same fix loop pattern as spec issues.

### Step 8: Validate & Complete Wave

1. Run validation across the entire project: `bash("go build ./...")`, `bash("go test ./...")`, `bash("go vet ./...")`
2. **If validation fails (Integration Issue):**
   - Individual tasks passed their own tests, but merged code causes build/test failures (e.g., interface mismatch, import cycle, conflicting type definitions).
   - Analyze the error output to identify which tasks' code is conflicting.
   - Dispatch a single **"Integration Fixer"** subagent via the `subagent` tool with top-level `task` (and optional `agent`) with:
     - The exact compilation/test error output
     - The list of files changed by each task in this wave
     - Instructions to resolve the integration conflict with minimal changes
   - Re-run validation after the fix.
   - Repeat until validation passes. If 3 fix attempts fail, **STOP and report** to the user — the plan likely has a design issue that needs human judgment.
3. Mark all wave tasks complete in `todo`
4. Proceed to next wave

## Phase 3: Final Review

After all waves complete, dispatch `code-reviewer` agent for the entire implementation (all files changed across all waves).

When invoked by `/go`, skip this phase and return the execution report instead; the `go` coordinator owns the separate whole-change review gate.

```json
{
  "agent": "code-reviewer",
  "task": "Review entire implementation across all waves..."
}
```

## File Conflict Safety

Before dispatching parallel implementers within a wave, verify no file conflicts:

1. Read the **Files** section of each task in the wave
2. Build a file → tasks map
3. Check for conflicts:
   - Two tasks **CREATE** the same file → **CONFLICT** — cannot run in parallel
   - Two tasks **MODIFY** the same file → **CONFLICT** — cannot run in parallel
   - One CREATES, another READS → safe (read is implicit, not a conflict)

**If conflicts found within a wave:**
- Split into sub-waves: group non-conflicting tasks together
- Sub-wave A runs first (parallel), then sub-wave B (parallel), etc.
- Log the split decision for transparency

**Example conflict resolution:**
```
Wave 0 has: Task 1, Task 2, Task 3
Task 1 modifies: internal/di/wire.go
Task 2 modifies: internal/di/wire.go  ← conflicts with Task 1
Task 3 modifies: internal/shared/router.go (no conflict)

→ Sub-wave 0a: Task 1, Task 3 (parallel)
→ Sub-wave 0b: Task 2 (after 0a completes)
```

## Prompt Templates

Resolve these paths relative to `.pi/skills/subagent-driven-development/`.

- `./implementer-prompt.md` - Dispatch implementer subagent
- `./spec-reviewer-prompt.md` - Dispatch spec compliance reviewer subagent
- `./code-quality-reviewer-prompt.md` - Dispatch code quality reviewer subagent

### Additional Context for Parallel Implementers

When dispatching parallel implementers, add to each prompt:

```
## Parallel Execution Notice

You are running in parallel with other implementer agents in this wave.
Other agents are working on: [list other task names]

**DO NOT** modify any files outside your task's file list.
Your files: [list from plan]

If you discover you need to modify a file not in your list, STOP and report it
instead of making the change.
```

### Prior Wave Context (Artifact Passing)

When dispatching implementers in Wave N > 0, the Controller **MUST read the actual source files** generated by prior waves and extract the exact signatures/structs that downstream tasks depend on. A 1-line summary is NOT enough — subagents need compilable interface definitions.

````markdown
## Prior Wave Dependencies

You depend on the following components created in previous waves.
Here are their **actual signatures** (read from the generated source files):

### From Task 1 (Certificate CRUD):
`internal/certificate/usecase/interfaces.go`:
```go
type CertificateRepository interface {
    Create(ctx context.Context, cert *domain.Certificate) error
    GetByID(ctx context.Context, id string) (*domain.Certificate, error)
}
```

You MUST use these exact type names and method signatures.
If you find a mismatch between the plan and the actual generated code, STOP and report it.
````

**How the Controller does this:**
1. After a wave completes, identify which files downstream tasks depend on (from the plan's "Blocked by" + "Files" sections)
2. Read those files using the `read` tool
3. Extract exported interfaces, structs, and function signatures
4. Paste the actual code snippets into the next wave's implementer prompts

## Falling Back to Sequential

If the plan has **no dependency graph** or all tasks share files:
- Fall back to sequential execution (one task at a time via `subagent` with top-level `agent` + `task`)
- Same review gates still apply
- This is equivalent to the original sequential behavior
- In `/go` mode, sequential fallback still requires subagent dispatch for every implementation task; the coordinator must never implement planned changes directly

## Red Flags

**Never:**
- Start implementation on main/master branch without explicit user consent
- Skip reviews (spec compliance OR code quality)
- Proceed with unfixed issues
- **Dispatch parallel subagents that modify the same files** — verify file lists first
- Make subagent read plan file (provide full text instead)
- Skip scene-setting context (subagent needs to understand where task fits)
- Ignore subagent questions (answer before letting them proceed)
- Accept "close enough" on spec compliance
- Skip review loops (reviewer found issues = implementer fixes = review again)
- Let implementer self-review replace actual review (both are needed)
- **Start code quality review before spec compliance is ✅** (wrong order)
- Move to next task/wave while either review has open issues
- **Start next wave before current wave fully completes** (dependency violation)

## Integration

**Required workflow skills:**
- **/skill:writing-plans** - Creates the plan with dependency graph that this skill executes
- **/skill:code-review-expert** - Code review template for reviewer subagents

**Subagents should use:**
- **/skill:test-driven-development** - Subagents follow TDD for each task
- **/skill:e2e-testing** - Subagents add E2E tests when implementing API endpoints

**Alternative workflows:**
- **/skill:executing-plans** - Use for parallel session instead of same-session execution
- **/skill:dispatching-parallel-agents** - Use for concurrent independent ad-hoc problems, NOT for plan execution
