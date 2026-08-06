---
name: writing-plans
description: Use after requirements or design are understood and before coding a multi-step change. Creates concrete implementation plans with vertical slices, dependency graphs, file lists, test cases, and project-derived validation commands.
---

# Writing Plans

Write an executable plan when dependencies, contracts, risk, or handoff justify a durable artifact. Prefer small vertical slices, DRY, YAGNI, and behavior-focused verification.

For a small sequential change, use a lightweight in-session plan instead of this full template:

```markdown
## Goal
## Assumptions
## Steps
1. [change and focused verification]
2. [change and focused verification]
## Final validation
```

Use the full format for cross-domain/layer work, parallel execution, high-risk changes, work spanning sessions, or approval. It defines contracts, files, acceptance criteria, and test scenarios without production method bodies or complete test code.

Save durable plans to the project's existing planning location; use `docs/plans/YYYY-MM-DD-<feature-name>.md` only when no convention exists.

## 1. Gather Context

Before planning:

- inspect the relevant implementation, tests, project instructions, architecture documents, ADRs, CI configuration, and existing commands;
- identify the project’s terminology, module boundaries, naming conventions, test strategy, and validation commands;
- record unknowns and resolve them before planning rather than inventing architecture or tooling.

Do not assume a language, framework, directory layout, dependency-injection mechanism, API style, test runner, or E2E harness.

## 2. Draft Vertical Slices

A task is a narrow, complete, independently verifiable behavior—not a horizontal “implement all models” or “write all handlers” phase.

For each proposed slice, provide:

- **Title**
- **Type:** AFK or HITL
- **Blocked by:** explicit dependencies
- **Areas touched:** use the project’s actual terminology

When this skill runs standalone, ask the user to approve the breakdown before writing the detailed plan. In an approved autonomous pipeline, continue unless an unresolved HITL decision remains.

## 3. Write the Plan

Every plan begins with:

```markdown
# [Feature Name] Implementation Plan

**Goal:** [One sentence]

**Architecture:** [How this fits existing project architecture]

**Validation:** [Commands discovered from this project]

## Task Dependency Graph

| Task | Type | Blocked by | Parallelizable with |
|------|------|------------|---------------------|
| 1 | AFK | None | 2 |
```

For every task in a full plan, use this structure. Omit sections that genuinely do not apply rather than manufacturing content:

```markdown
### Task N: [Vertical Slice Name]

**Type:** AFK / HITL
**Blocked by:** Task X / None
**Areas touched:** [project-specific areas]

**Goal:** [One or two sentences]

**Acceptance Criteria:**
- [ ] Observable behavior and outcome
- [ ] Relevant error, boundary, or edge case

**Files:**
- Create: `actual/project/path`
- Modify: `actual/project/path`

#### Interface Contracts

Describe only new or changed public/internal contracts that downstream work needs:

```text
[Use the project’s language and syntax. Include signatures, schemas, events, or configuration shape—not implementations.]
```

#### Test Cases to Cover

- Primary success behavior
- Relevant validation, error, and boundary behavior
- Integration/E2E behavior when the project has a suitable harness or requirements demand it

#### Validation

```text
[Focused test command discovered from the project]
[Relevant build, lint, typecheck, or integration command]
```
```

## 4. E2E Planning

If public behavior changes, determine from the repository whether E2E coverage exists or is required. If so, add a dependent E2E task that names:

- the public interfaces and scenarios to cover;
- the existing harness, fixtures, and cleanup strategy;
- the actual project command that must pass.

Do not add an E2E task solely because a generic workflow says so when the project has no such harness and the requirements do not call for it.

## 5. Plan Coverage Check

Before handoff, confirm:

- [ ] Every approved requirement maps to a task.
- [ ] Every task has observable acceptance criteria and behavior-focused tests.
- [ ] Every task lists exact file paths after repository inspection.
- [ ] Dependencies have no cycles.
- [ ] Parallel tasks do not modify the same file.
- [ ] Public-interface changes have appropriate integration/E2E coverage when applicable.
- [ ] Validation commands come from this project.
- [ ] Assumptions and deviations are explicit.

## Handoff

Execute sequentially when tasks share evolving context, interfaces, generated outputs, or files. Hand the plan to `subagent-driven-development` only when independent slices have stable contracts, non-overlapping write sets, and plausible parallel benefit. In standalone use, ask before execution; in an approved autonomous pipeline, continue. Provide each implementer the full task text rather than asking it to rediscover the plan.
