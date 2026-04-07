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
Requirements → Design & Planning → Implementation → Code Review → Simplification
(brainstorming)  (writing-plans)  (executing-plans) (code-reviewer) (code-simplifier)
```

- Only the first step (brainstorming) requires human confirmation; the rest runs automatically
- Each step includes validation (`go build`, `go vet`, `go test`)
- Issues found during review are auto-fixed and re-verified

### Step-by-Step Workflow

For finer-grained control, use individual skills:

| Stage | Skill | Description |
|-------|-------|-------------|
| Requirements | `brainstorming` | Turn ideas into complete designs and specs |
| Planning | `writing-plans` | Write detailed implementation plans (with TDD) |
| Execution | `executing-plans` | Execute plans in batches with review checkpoints |
| Execution (parallel) | `subagent-driven-development` | Subagent-driven development with two-stage review |
| Review | `code-review-expert` | SOLID, security, architecture review |

## Skills Overview

### Flow (decides how to handle a task)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `go` | End-to-end automated development | `/go` + requirement description |
| `brainstorming` | Requirements exploration & design | **Must use** before creating features, components, or modifying behavior |
| `using-superpowers` | Skill usage guide | Ensures correct skill invocation |

### Planning (creates implementation plans)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `writing-plans` | Write implementation plans | Multi-step tasks with specs or requirements |

### Execution (executes plans)

| Skill | Purpose | Trigger |
|-------|---------|---------|
| `executing-plans` | Execute plans in batches | Separate session with checkpoints |
| `subagent-driven-development` | Subagent-driven development | Current session, one subagent per task |
| `dispatching-parallel-agents` | Dispatch parallel agents | 2+ independent tasks |
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
3. **Execution** (`executing-plans`, `subagent-driven-development`, `test-driven-development`) — Execute the plan
4. **Dispatch** (`dispatching-parallel-agents`) — Solve multiple independent problems in parallel
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
