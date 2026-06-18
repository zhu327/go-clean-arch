# Go Clean Architecture Template — AI Agent Configuration Index

> See `.cursor/rules/` for detailed coding standards.

## Invocation

In Cursor, use `/` to reference skill files or describe requirements in conversation. AI will automatically identify and apply the relevant skill.

## Recommended Development Workflow

### `/go` — End-to-End Development (Recommended)

The most efficient way to develop features. Simply describe your requirement and AI handles the full pipeline:

```
/go Add an Article domain module with CRUD operations
```

Automated pipeline:
```
Requirements → Planning → Wave Execution → Code Review → Simplification
(brainstorming)  (writing-plans)  (subagent-driven-development) (code-reviewer) (code-simplifier)
```

- Only the first step (brainstorming) requires human confirmation; the rest runs automatically
- Planning includes a coverage checklist before execution
- Execution uses per-task spec review plus wave validation (`go build`, `go vet`, `go test`)
- Final review focuses on global architecture/code quality; issues found during review are auto-fixed and re-verified

### Step-by-Step Workflow

For finer-grained control, use individual skills:

| Stage | Skill | Description |
|-------|-------|-------------|
| Requirements | `brainstorming` | Turn ideas into complete designs and specs |
| Planning | `writing-plans` | Write detailed implementation plans (with TDD) |
| Execution | `subagent-driven-development` | Execute plans with dependency-aware wave parallelism, per-task spec review, and final global code review |
| Review | `code-review-expert` | SOLID, security, architecture review |

## Skills Overview

### Flow (decides how to handle a task)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `go` | End-to-end automated development | `/go` + requirement description |
| `brainstorming` | Requirements exploration & design | **Must use** before creating features, components, or modifying behavior |

### Planning (creates implementation plans)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `writing-plans` | Write implementation plans | Multi-step tasks with specs or requirements |

### Execution (executes plans)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `subagent-driven-development` | Subagent-driven development | Current session, one subagent per task, wave-parallel when safe |
| `test-driven-development` | TDD development | New features, bug fixes, refactoring |

### Review (quality assurance)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `code-review-expert` | Code review | Review git changes, SOLID/security checks |

## Subagents

| Agent | Purpose | Auto-Dispatch Scenario |
|-------|---------|------------------------|
| `code-reviewer` | Code review agent | `/go` step 4 (auto); after completing major project steps (manual) |
| `code-simplifier` | Code simplification agent | `/go` step 5 (auto); when code needs simplification (manual) |

## Skill Priority

When multiple skills may apply, use them in this order:

1. **Flow** (`brainstorming`) — Decide how to approach the task
2. **Planning** (`writing-plans`) — Create the implementation plan
3. **Execution** (`subagent-driven-development`, `test-driven-development`) — Execute the plan
4. **Review** (`code-review-expert`) — Quality assurance
5. **Review** (`code-review-expert`) — Quality assurance

## Learned Patterns

Experience patterns extracted from project practice, located in `.cursor/skills/learned/`:

| Pattern | Description |
|---------|-------------|
| `gorm-preload-pattern` | GORM Preload to avoid N+1 queries |

Use `/learn` to extract new experience patterns from the current session.

## Common Commands

```bash
make build   # Build
make serve   # Start the server
make di      # Generate dependency injection code (Wire)
make doc     # Generate Swagger documentation
make lint    # Run linter (golangci-lint)
make test    # Run tests (with coverage)
make fmt     # Format code (gofumpt + golines)
make mock    # Generate mocks (mockgen)
```

## Architecture Overview

Uses **DDD + Clean Architecture**, organized by domain modules:

```
internal/
├── {domain}/              # Domain module (e.g., user/)
│   ├── domain/            # Entities, business rules
│   ├── usecase/           # Manager, interfaces, DTOs, mocks
│   └── adapter/           # Handler, Router, Repository
├── shared/                # Shared infrastructure (Server, Router, Middleware)
└── di/                    # Wire dependency injection
```

Feature development order: Domain → UseCase → Repository → Handler → Router → Swagger → Tests → Wire DI

## Agent skills

### Issue tracker

Issues and PRDs live as markdown files under `.scratch/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Default canonical labels (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context repo: `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.
