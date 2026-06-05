---
name: dispatching-parallel-agents
description: Use for 2+ independent ad-hoc bugs, investigations, or fixes that can run concurrently without shared files or sequential dependencies. Not for executing written plans; use subagent-driven-development for planned work.
disable-model-invocation: true
---

# Dispatching Parallel Agents

## Overview

When you have multiple **independent problems** to solve (different bugs, different subsystems, different investigations), dispatch one agent per problem and let them work concurrently.

**Core principle:** One agent per independent problem domain. Concurrent execution. No shared state.

**Tool used:** `subagent` — use `tasks` array to dispatch multiple agents simultaneously.

**This skill is for: solving N independent problems in parallel.**
**NOT for: executing a sequential implementation plan** — use /skill:subagent-driven-development for that.

## When to Use

```
Multiple tasks/problems?
  ├─ no → Single agent handles all
  └─ yes → Following an implementation plan?
            ├─ yes → Use /skill:subagent-driven-development
            └─ no → Are tasks independent (no shared files)?
                     ├─ yes → Parallel dispatch (this skill)
                     └─ no → Sequential agents
```

**Use when:**
- 2+ independent problems (bugs, failures, investigations) across different domains
- Each problem can be understood and solved without context from others
- Agents won't edit the same files
- No implementation plan — just problems to fix

**Don't use when:**
- Following an implementation plan (use /skill:subagent-driven-development)
- Failures are related (fix one might fix others)
- Need to understand full system state first
- Agents would interfere (editing same files)

## The Pattern

### 1. Identify Independent Domains

Group failures by what's broken (按领域模块分组):
- Certificate domain: 证书同步、拓扑构建
- SubDomain domain: 子域名 CRUD、DNS 集成
- WorkOrder domain: 工单创建、状态流转

Each domain is independent - fixing Certificate doesn't affect WorkOrder tests.

### 2. Create Focused Agent Tasks

Each agent gets:
- **Specific scope:** One test file or subsystem
- **Clear goal:** Make these tests pass
- **Constraints:** Don't change other code
- **Expected output:** Summary of what you found and fixed

### 3. Dispatch in Parallel

Use the `subagent` tool with `tasks` array:

```json
{
  "tasks": [
    { "task": "Fix certificate manager test failures\n\n[details]" },
    { "task": "Fix subdomain repository test failures\n\n[details]" },
    { "task": "Fix workorder usecase test failures\n\n[details]" }
  ]
}
```
All three run concurrently. Add an `agent` field only when a suitable implementation agent exists.

### 4. Review and Integrate

When agents return:
- Read each summary from the `subagent` output
- Verify fixes don't conflict
- Run `bash("make test && make lint")` to verify
- Integrate all changes

## Agent Prompt Structure

Good agent prompts are:
1. **Focused** - One clear problem domain
2. **Self-contained** - All context needed to understand the problem
3. **Specific about output** - What should the agent return?

```markdown
Fix the 3 failing tests in internal/certificate/usecase/manager_test.go:

1. TestManager_SyncCertificate - expects cert to be updated but got nil
2. TestManager_GetExpiredCerts - returns wrong count
3. TestManager_BuildTopology - missing edge connections

These are likely mock setup or assertion issues. Your task:

1. Read the test file and understand what each test verifies
2. Identify root cause - mock expectations or actual bugs?
3. Fix by:
   - Correcting mock EXPECT setup with proper DoAndReturn
   - Fixing bugs in implementation if found
   - Ensuring complete mock response structs

Do NOT just skip assertions - find the real issue.

Return: Summary of what you found and what you fixed.
```

## Common Mistakes

**❌ Too broad:** "Fix all the tests" - agent gets lost
**✅ Specific:** "Fix internal/certificate/usecase/manager_test.go" - focused scope

**❌ No context:** "Fix the race condition" - agent doesn't know where
**✅ Context:** Paste the error messages and test names

**❌ No constraints:** Agent might refactor everything
**✅ Constraints:** "Do NOT change production code" or "Fix tests only"

**❌ Vague output:** "Fix it" - you don't know what changed
**✅ Specific:** "Return summary of root cause and changes"

## vs. Subagent-Driven Development

| | Parallel Agents (this skill) | Subagent-Driven Development |
|---|---|---|
| **Input** | Ad-hoc problems/failures | Implementation plan from /skill:writing-plans |
| **Execution** | All agents run concurrently | Dependency-aware waves; parallel within safe waves, sequential across dependency layers |
| **Review** | Post-integration only | Two-stage review after each task |
| **When** | Multiple independent bugs/issues | Building features from a plan |
| **Agent count** | N agents simultaneously | 1 implementer + 2 reviewers per task |
| **Coordination** | None (independent) | Controller orchestrates sequence |

**Rule of thumb:** "Fix N broken things" → this skill. "Build something from a plan" → /skill:subagent-driven-development.

## When NOT to Use

**Related failures:** Fixing one might fix others - investigate together first
**Need full context:** Understanding requires seeing entire system
**Exploratory debugging:** You don't know what's broken yet
**Shared state:** Agents would interfere (editing same files, using same resources)
**Implementation plan:** Use /skill:subagent-driven-development for planned, sequential work

## Real Example from Session

**Scenario:** 6 test failures across 3 domains after major refactoring

**Dispatch:**
```json
{
  "tasks": [
    { "task": "Fix internal/certificate/usecase/manager_test.go" },
    { "task": "Fix internal/subdomain/adapter/repository/subdomain_test.go" },
    { "task": "Fix internal/workorder/usecase/manager_test.go" }
  ]
}
```

**Integration:** All fixes independent, no conflicts, `make test` green

**Time saved:** 3 problems solved in parallel vs sequentially
